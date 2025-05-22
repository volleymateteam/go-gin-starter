package video

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-gin-starter/pkg/logger"
	"go-gin-starter/repositories"
)

// QueueManager handles SQS operations for video processing
type QueueManager struct {
	sqs        *sqs.SQS
	queueURL   string
	processor  *VideoProcessor
	maxRetries int
}

// NewQueueManager creates a new queue manager instance
func NewQueueManager(sqsClient *sqs.SQS, queueURL string, processor *VideoProcessor) *QueueManager {
	return &QueueManager{
		sqs:        sqsClient,
		queueURL:   queueURL,
		processor:  processor,
		maxRetries: 3,
	}
}

// EnqueueVideo adds a video processing job to the queue
func (q *QueueManager) EnqueueVideo(job *VideoProcessingJob) error {
	job.CreatedAt = time.Now()
	job.Status = StatusPending

	messageBody, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = q.sqs.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	return err
}

// StartProcessing starts processing videos from the queue
func (q *QueueManager) StartProcessing() {
	for {
		result, err := q.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(q.queueURL),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(20),
		})

		if err != nil {
			logger.Error("Failed to receive message from SQS", zap.Error(err))
			continue
		}

		for _, message := range result.Messages {
			var job VideoProcessingJob
			if err := json.Unmarshal([]byte(*message.Body), &job); err != nil {
				logger.Error("Failed to unmarshal job", zap.Error(err))
				continue
			}

			job.Status = StatusProcessing

			thumbnailURL, err := q.processor.ProcessVideo(&job)
			if err != nil {
				logger.Error("Failed to process video",
					zap.String("match_id", job.MatchID),
					zap.Error(err))
				job.Status = StatusFailed
				job.Error = err.Error()
			} else {
				job.Status = StatusCompleted
			}

			matchID, parseErr := uuid.Parse(job.MatchID)
			if parseErr == nil {
				matchRepo := repositories.NewMatchRepository()
				match, getErr := matchRepo.GetByID(matchID)
				if getErr == nil {
					match.ThumbnailURL = thumbnailURL
					if uploadErr := matchRepo.Update(match); uploadErr != nil {
						logger.Error("failed to update match thumbnail",
							zap.Error(uploadErr))
					}
				} else {
					logger.Error("failed to fetch match for thumbnail update",
						zap.Error(getErr))
				}
			} else {
				logger.Error("invalid match UUID", zap.Error(parseErr))
			}

			// Delete message from queue if processed successfully
			if job.Status == StatusCompleted {
				_, err = q.sqs.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(q.queueURL),
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					logger.Error("Failed to delete message", zap.Error(err))
				}
			}
		}
	}
}
