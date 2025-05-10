package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger
func Init() {
	// Set up logger configuration
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	// Check if we're in production
	if os.Getenv("ENV") == "production" {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true
	}

	var err error
	log, err = config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// Middleware returns a gin middleware for logging HTTP requests
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Process request
		c.Next()

		// Log request details after completion
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		// Get the error if any
		errors := c.Errors.Errors()

		// Create structured log fields
		fields := []zapcore.Field{
			zap.String("request_id", requestID),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.Int("status", statusCode),
			zap.String("path", path),
			zap.Duration("latency", latency),
		}

		if len(errors) > 0 {
			fields = append(fields, zap.Strings("errors", errors))
		}

		// Log based on status code
		if statusCode >= 500 {
			Error("Server error", fields...)
		} else if statusCode >= 400 {
			Warn("Client error", fields...)
		} else {
			Info("Request completed", fields...)
		}
	}
}

// With creates a child logger with additional fields
func With(fields ...zapcore.Field) *zap.Logger {
	return log.With(fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zapcore.Field) {
	log.Fatal(msg, fields...)
}

// GetRequestLogger returns a logger with request context
func GetRequestLogger(c *gin.Context) *zap.Logger {
	requestID, exists := c.Get("request_id")
	if !exists {
		requestID = uuid.New().String()
	}

	userID, _ := c.Get("user_id")

	return log.With(
		zap.String("request_id", requestID.(string)),
		zap.Any("user_id", userID),
	)
}

// Sync flushes any buffered log entries
func Sync() {
	_ = log.Sync()
}
