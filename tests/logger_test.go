package tests

import (
	"qualifire-home-assignment/internal/loggers"
	"qualifire-home-assignment/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger_Info(t *testing.T) {
	entry := &models.LogEntry{
		VirtualKey: "test-key",
		Provider:   "openai",
		Status:     200,
		DurationMS: int64(200 * time.Millisecond),
	}

	logger := loggers.Logger{Entry: entry}

	// Should not panic
	assert.NotPanics(t, func() {
		logger.Info()
	})
}

func TestLogger_Error(t *testing.T) {
	entry := &models.LogEntry{
		VirtualKey: "test-key",
		Provider:   "openai",
		Status:     500,
		DurationMS: int64(200 * time.Millisecond),
	}

	logger := loggers.Logger{Entry: entry}

	// Should not panic
	assert.NotPanics(t, func() {
		logger.Error()
	})
}

func TestLogger_InfoWithNilEntry(t *testing.T) {
	logger := loggers.Logger{Entry: nil}

	// Should handle nil gracefully
	assert.NotPanics(t, func() {
		logger.Info()
	})
}

func TestLogger_ErrorWithNilEntry(t *testing.T) {
	logger := loggers.Logger{Entry: nil}

	// Should handle nil gracefully
	assert.NotPanics(t, func() {
		logger.Error()
	})
}
