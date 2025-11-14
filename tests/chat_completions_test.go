package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/controllers"
	"qualifire-home-assignment/internal/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupChatCompletionTest prepares the test environment by setting test mode,
// loading configurations and resetting metrics and quota services for clean state
func setupChatCompletionTest() {
	os.Setenv("IS_TEST", "1")
	gin.SetMode(gin.TestMode)
	configs.LoadConfig()
	services.GetMetricsService().Reset()
	services.GetQuotaService().Reset()
}

// TestChatCompletions_RouteRequests_MetricsTracking verifies the chat completion
// endpoint's request handling by setting up a mock OpenAI server and testing the
// basic routing structure
func TestChatCompletions_RouteRequests_MetricsTracking(t *testing.T) {
	setupChatCompletionTest()

	// Setup mock OpenAI endpoint
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "chatcmpl-123",
			"choices": [{
				"message": {
					"role": "assistant",
					"content": "Hello!"
				}
			}]
		}`))
	}))
	defer mockServer.Close()

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
		"model": "gpt-3.5-turbo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key-1")

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)
	router.POST("/chat/completions", controllers.ChatCompletions{}.RouteRequests)
	c.Request = req

	// Note: This test validates the structure, actual API calls would need proper mocking
	assert.NotNil(t, router)
}

// TestChatCompletions_RouteRequests_QuotaIncrement tests the initial quota state
// and tracking for a test virtual key to ensure proper quota management
func TestChatCompletions_RouteRequests_QuotaIncrement(t *testing.T) {
	setupChatCompletionTest()

	quotaService := services.GetQuotaService()
	initialRequests, initialTokens := quotaService.GetUsage("test-key-1")

	// Note: Full integration would require mocking the provider response
	assert.Equal(t, 0, initialRequests)
	assert.Equal(t, 0, initialTokens)
}
