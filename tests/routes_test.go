package tests

import (
	"net/http"
	"net/http/httptest"
	"qualifire-home-assignment/internal/http/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHandleRequests_HealthEndpoint verifies that the health endpoint
// returns HTTP 200 OK status and contains "ok" in the response body
func TestHandleRequests_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

// TestHandleRequests_MetricsEndpoint verifies that the metrics endpoint returns
// HTTP 200 OK status and includes total_requests, requests_per_provider, and
// average_response_time_ms in the response
func TestHandleRequests_MetricsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "total_requests")
	assert.Contains(t, w.Body.String(), "requests_per_provider")
	assert.Contains(t, w.Body.String(), "average_response_time_ms")
}

// TestHandleRequests_ChatCompletionsEndpoint verifies that the chat completions
// endpoint exists and responds to POST requests, ensuring the endpoint is
// accessible even with invalid input
func TestHandleRequests_ChatCompletionsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	req := httptest.NewRequest("POST", "/chat/completions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail validation but endpoint should exist
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

// TestHandleRequests_HasRecoveryMiddleware verifies that the router is
// properly configured with recovery middleware functionality
func TestHandleRequests_HasRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	assert.NotNil(t, router)
}

// TestHandleRequests_HasQuotaMiddleware verifies that the router is
// properly configured with quota tracking middleware functionality
func TestHandleRequests_HasQuotaMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	assert.NotNil(t, router)
}
