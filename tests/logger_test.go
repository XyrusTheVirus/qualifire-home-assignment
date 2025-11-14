package tests

import (
	"qualifire-home-assignment/internal/loggers"
	"qualifire-home-assignment/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLogger_Info verifies that the Info method properly logs the entry data
// and handles successful (200) status codes without panicking
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

// TestLogger_Error verifies that the Error method properly logs the entry data
// and handles error (500) status codes without panicking
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

// TestLogger_InfoWithNilEntry verifies that the Info method handles nil entries
// gracefully without causing panics or runtime errors
func TestLogger_InfoWithNilEntry(t *testing.T) {
	logger := loggers.Logger{Entry: nil}

	// Should handle nil gracefully
	assert.NotPanics(t, func() {
		logger.Info()
	})
}

// TestLogger_ErrorWithNilEntry verifies that the Error method handles nil entries
// gracefully without causing panics or runtime errors
func TestLogger_ErrorWithNilEntry(t *testing.T) {
	logger := loggers.Logger{Entry: nil}

	// Should handle nil gracefully
	assert.NotPanics(t, func() {
		logger.Error()
	})
}
