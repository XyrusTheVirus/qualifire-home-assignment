package tests

import (
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetError(t *testing.T) {
	err := errors.GetError("TEST_ERROR", "Test error message", http.StatusBadRequest)

	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "Test error message", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

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

func TestValidation_GetError(t *testing.T) {
	v := errors.Validation{}
	err := v.GetError("Validation failed", http.StatusBadRequest)

	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Equal(t, "Validation failed", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestApiProvider_GetError(t *testing.T) {
	a := errors.ApiProvider{}
	err := a.GetError("Provider request failed", http.StatusInternalServerError)

	assert.Equal(t, "LLM_PROVIDER_ERROR", err.Code)
	assert.Equal(t, "Provider request failed", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
}
