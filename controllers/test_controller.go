package controllers

import (
	"go-gin-starter/pkg/video"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	Queue *video.QueueManager
}

func NewTestController(queue *video.QueueManager) *TestController {
	return &TestController{
		Queue: queue,
	}
}

func (tc *TestController) EnqueueTestJob(c *gin.Context) {
	job := &video.VideoProcessingJob{
		MatchID:   "test-match-id",
		InputKey:  "raw/2025_germany/test/testfile.mp4",
		OutputKey: "compressed/2025_germany/test/testfile.mp4",
		CreatedAt: time.Now(),
		Status:    video.StatusPending,
	}

	if err := tc.Queue.EnqueueVideo(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to enqueue job",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test job enqueued"})
}
