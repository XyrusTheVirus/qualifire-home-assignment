package controllers

import (
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/http/validators"
	"qualifire-home-assignment/internal/models"
	"qualifire-home-assignment/internal/providers"
	"qualifire-home-assignment/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

// ChatCompletions handles chat completion requests
type ChatCompletions struct{}

// RouteRequests handles the chat completion requests
func (cp ChatCompletions) RouteRequests(c *gin.Context) {
	startTime := time.Now()

	req := validators.ChatCompletion{}.Validate(c)
	proxyReq := req.(models.ProxyRequest)

	// Check quota using a virtual key from context
	virtualKey := proxyReq.VirtualKey
	quotaService := services.GetQuotaService()
	allowed, reason := quotaService.CheckQuota(virtualKey)
	if !allowed {
		err := errors.GetError("QUOTA_EXCEEDED", reason, http.StatusTooManyRequests)
		c.JSON(err.StatusCode, err.ToGin())
		return
	}

	result := providers.Factory(proxyReq).SendRequest()

	// Track metrics
	duration := time.Since(startTime)
	metricsService := services.GetMetricsService()
	metricsService.RecordRequest(services.RequestMetric{
		Provider:   proxyReq.Provider,
		VirtualKey: virtualKey,
		Duration:   duration,
		Status:     http.StatusOK,
		Timestamp:  startTime,
	})

	// Increment quota (estimate 100 tokens per request for simplicity)
	quotaService.IncrementRequest(virtualKey, 100)

	Success(c, http.StatusOK, result)
}
