package video

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"

	"go-gin-starter/pkg/logger"
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
			if err := q.processor.ProcessVideo(&job); err != nil {
				logger.Error("Failed to process video",
					zap.String("match_id", job.MatchID),
					zap.Error(err))
				job.Status = StatusFailed
				job.Error = err.Error()
			} else {
				job.Status = StatusCompleted
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
