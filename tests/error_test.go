package tests

import (
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetError verifies that GetError function correctly creates an error with the specified code, message and status code
func TestGetError(t *testing.T) {
	err := errors.GetError("TEST_ERROR", "Test error message", http.StatusBadRequest)

	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "Test error message", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

// TestError_ToGin verifies that Error.ToGin correctly converts an error to the Gin response format with all required fields
func TestError_ToGin(t *testing.T) {
	err := errors.Error{
		Code:       "TEST_ERROR",
		Message:    "Test error message",
		Details:    "Stack trace details",
		StatusCode: http.StatusBadRequest,
	}

	ginResponse := err.ToGin()

	assert.Contains(t, ginResponse, "code")
	assert.Contains(t, ginResponse, "message")
	assert.Contains(t, ginResponse, "details")
	assert.Equal(t, "TEST_ERROR", ginResponse["code"])
	assert.Equal(t, "Test error message", ginResponse["message"])
	assert.Equal(t, "Stack trace details", ginResponse["details"])
}

// TestValidation_GetError verifies that Validation.GetError correctly creates a validation error with appropriate code and status
func TestValidation_GetError(t *testing.T) {
	v := errors.Validation{}
	err := v.GetError("Validation failed", http.StatusBadRequest)

	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Equal(t, "Validation failed", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

// TestApiProvider_GetError verifies that ApiProvider.GetError correctly creates a provider error with appropriate code and status
func TestApiProvider_GetError(t *testing.T) {
	a := errors.ApiProvider{}
	err := a.GetError("Provider request failed", http.StatusInternalServerError)

	assert.Equal(t, "LLM_PROVIDER_ERROR", err.Code)
	assert.Equal(t, "Provider request failed", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
}
