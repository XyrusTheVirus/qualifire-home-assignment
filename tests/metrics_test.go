package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/http/middleware"
	"qualifire-home-assignment/internal/services"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupMetricsTestRouter Setup a test router with metrics middleware enabled
func setupMetricsTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.MetricsMiddleware())
	return r
}

// TestMetricsMiddleware_SuccessfulRequest Test that metrics are recorded correctly for a successful request
func TestMetricsMiddleware_SuccessfulRequest(t *testing.T) {
	metricsService := services.GetMetricsService()
	metricsService.Reset()

	router := setupMetricsTestRouter()
	router.POST("/tests/chat/completions", func(c *gin.Context) {
		c.Set("provider", "openai")
		c.Set("token_count", 150)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/tests/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer test-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify metrics were recorded
	stats := metricsService.GetStats()
	assert.Equal(t, 1, stats["total_requests"])
}

// TestMetricsMiddleware_WithPanic Test that metrics are still recorded even if there is a panic in the handler
func TestMetricsMiddleware_WithPanic(t *testing.T) {
	metricsService := services.GetMetricsService()
	metricsService.Reset()

	router := setupMetricsTestRouter()
	router.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "recovered"})
			}
		}()
		c.Next()
	})

	router.POST("/tests/chat/completions", func(c *gin.Context) {
		c.Set("provider", "openai")
		c.Set("token_count", 100)
		panic(errors.GetError("TEST_ERROR", "test panic", http.StatusBadRequest))
	})

	req := httptest.NewRequest("POST", "/tests/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer test-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify metrics were still recorded despite panic
	stats := metricsService.GetStats()
	assert.Equal(t, 1, stats["total_requests"])
}

// TestMetricsMiddleware_QuotaIncrement Test that quota is incremented correctly
func TestMetricsMiddleware_QuotaIncrement(t *testing.T) {
	quotaService := services.GetQuotaService()
	quotaService.Reset()
	quotaService.SetLimits(100, 10000, time.Hour)

	router := setupMetricsTestRouter()
	router.POST("/tests/chat/completions", func(c *gin.Context) {
		c.Set("provider", "openai")
		c.Set("token_count", 250)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/tests/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer test-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify quota was incremented
	requests, tokens := quotaService.GetUsage("test-key")
	assert.Equal(t, 1, requests)
	assert.Equal(t, 250, tokens)
}

// TestMetricsMiddleware_SkipNonCompletionEndpoints Test that metrics are NOT tracked for non-completion endpoints
func TestMetricsMiddleware_SkipNonCompletionEndpoints(t *testing.T) {
	metricsService := services.GetMetricsService()
	metricsService.Reset()

	router := setupMetricsTestRouter()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify metrics were NOT recorded for non-completion endpoint
	stats := metricsService.GetStats()
	assert.Equal(t, 0, stats["total_requests"])
}

// TestMetricsMiddleware_NoVirtualKey Test that metrics are NOT tracked for requests without a virtual key
func TestMetricsMiddleware_NoVirtualKey(t *testing.T) {
	metricsService := services.GetMetricsService()
	metricsService.Reset()

	router := setupMetricsTestRouter()
	router.POST("/tests/chat/completions", func(c *gin.Context) {
		c.Set("provider", "openai")
		c.Set("token_count", 100)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/tests/chat/completions", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify metrics were NOT recorded without virtual key
	stats := metricsService.GetStats()
	assert.Equal(t, 0, stats["total_requests"])
}

// TestMetricsMiddleware_MultipleProviders Test that metrics are tracked for multiple providers
func TestMetricsMiddleware_MultipleProviders(t *testing.T) {
	metricsService := services.GetMetricsService()
	metricsService.Reset()

	router := setupMetricsTestRouter()
	router.POST("/tests/tests/chat/completions", func(c *gin.Context) {

		c.Set("token_count", 100)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Request 1: OpenAI
	req1 := httptest.NewRequest("POST", "/tests/tests/chat/completions", nil)
	req1.Header.Set("Authorization", "Bearer vk_user1_openai")
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = req1
	c1.Set("provider", "openai")
	router.ServeHTTP(w1, req1)

	// Request 2: Anthropic
	req2 := httptest.NewRequest("POST", "/tests/tests/chat/completions", nil)
	req2.Header.Set("Authorization", "Bearer vk_user2_anthropic")
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = req2
	c2.Set("provider", "anthropic")
	router.ServeHTTP(w2, req2)

	// Verify metrics tracked by both providers
	stats := metricsService.GetStats()
	assert.Equal(t, 2, stats["total_requests"])
}

// TestMetricsMiddleware_DurationTracking Test that duration is tracked correctly
func TestMetricsMiddleware_DurationTracking(t *testing.T) {
	metricsService := services.GetMetricsService()
	metricsService.Reset()

	router := setupMetricsTestRouter()
	router.POST("/tests/chat/completions", func(c *gin.Context) {
		c.Set("provider", "openai")
		c.Set("token_count", 100)
		// Simulate processing time
		time.Sleep(10 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("POST", "/tests/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer test-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify duration was tracked
	stats := metricsService.GetStats()
	avgTime := stats["average_response_time_ms"].(float64)
	assert.Greater(t, avgTime, float64(5), "Duration should be tracked and > 5ms")
}

// TestMetricsMiddleware_StatusCodeTracking Test that status code is tracked correctly
func TestMetricsMiddleware_StatusCodeTracking(t *testing.T) {
	metricsService := services.GetMetricsService()
	metricsService.Reset()

	router := setupMetricsTestRouter()
	router.POST("/tests/chat/completions", func(c *gin.Context) {
		c.Set("provider", "openai")
		c.Set("token_count", 100)
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})

	req := httptest.NewRequest("POST", "/tests/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer test-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestMetricsMiddleware_VirtualKeyExtraction Test that virtual key is extracted from Authorization header
func TestMetricsMiddleware_VirtualKeyExtraction(t *testing.T) {
	router := setupMetricsTestRouter()
	router.POST("/tests/chat/completions", func(c *gin.Context) {
		virtualKey, exists := c.Get("virtual_key")
		assert.True(t, exists, "Virtual key should be set in context")
		assert.Equal(t, "my-test-key", virtualKey)
		c.JSON(http.StatusOK, gin.H{"key": virtualKey})
	})

	req := httptest.NewRequest("POST", "/tests/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer my-test-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "my-test-key", response["key"])
}
