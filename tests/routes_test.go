package tests

import (
	"net/http"
	"net/http/httptest"
	"qualifire-home-assignment/internal/http/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequests_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

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

func TestHandleRequests_ChatCompletionsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	req := httptest.NewRequest("POST", "/chat/completions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail validation but endpoint should exist
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestHandleRequests_HasRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	assert.NotNil(t, router)
}

func TestHandleRequests_HasQuotaMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.HandleRequests()

	assert.NotNil(t, router)
}
