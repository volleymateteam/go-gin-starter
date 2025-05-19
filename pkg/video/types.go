package video

import "time"

// VideoProcessingJob represents a video processing task
type VideoProcessingJob struct {
	MatchID   string    `json:"match_id"`
	InputKey  string    `json:"input_key"`  // S3 key for raw video
	OutputKey string    `json:"output_key"` // S3 key for processed video
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
}

// VideoFormat represents different video formats
type VideoFormat struct {
	Resolution string
	Bitrate    string
	MaxSize    int64 // in bytes
}

const (
	// Status constants
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"

	// Folder structure
	RawVideoFolder   = "raw"
	CompressedFolder = "compressed"
	ThumbnailsFolder = "thumbnails"

	// Default formats
	Format1080p = "1080p"
	Format720p  = "720p"
	Format480p  = "480p"
)

// DefaultVideoFormats defines standard processing formats
var DefaultVideoFormats = map[string]VideoFormat{
	Format1080p: {
		Resolution: "1920x1080",
		Bitrate:    "4M",
		MaxSize:    1.5 * 1024 * 1024 * 1024, // 1.5GB
	},
	Format720p: {
		Resolution: "1280x720",
		Bitrate:    "2.5M",
		MaxSize:    800 * 1024 * 1024, // 800MB
	},
	Format480p: {
		Resolution: "854x480",
		Bitrate:    "1M",
		MaxSize:    400 * 1024 * 1024, // 400MB
	},
}
