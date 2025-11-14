package middleware

import (
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/services"
	"qualifire-home-assignment/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware tracks request metrics and quota usage
// This middleware ensures metrics are recorded even if panics occur
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		metricsService := services.GetMetricsService()
		quotaService := services.GetQuotaService()

		var virtualKey string
		var provider string
		var tokenCount int
		statusCode := 200

		// Defer to ensure metrics are always recorded
		defer func() {
			duration := time.Since(startTime)
			// After request processing, extract provider and token info from context
			if providerVal, exists := c.Get("provider"); exists {
				provider = providerVal.(string)
			}
			if tokenVal, exists := c.Get("token_count"); exists {
				tokenCount = tokenVal.(int)
			}

			// Capture panic information if present
			if err := recover(); err != nil {
				// Determine status code from error
				switch e := err.(type) {
				case errors.Error:
					statusCode = e.StatusCode
				case errors.ApiProvider:
					statusCode = e.StatusCode
				default:
					statusCode = 500
				}

				// Record metrics even on panic
				if virtualKey != "" {
					metricsService.RecordRequest(services.RequestMetric{
						Provider:   provider,
						VirtualKey: virtualKey,
						Duration:   duration,
						Status:     statusCode,
						Timestamp:  startTime,
					})

					// Increment quota even on failure to prevent abuse
					if tokenCount > 0 {
						quotaService.IncrementRequest(virtualKey, tokenCount)
					}
				}

				// Re-panic to let recovery middleware handle it
				panic(err)
			}

			// Normal execution - record metrics
			if virtualKey != "" {
				// Get actual status code from response
				statusCode = c.Writer.Status()

				metricsService.RecordRequest(services.RequestMetric{
					Provider:   provider,
					VirtualKey: virtualKey,
					Duration:   duration,
					Status:     statusCode,
					Timestamp:  startTime,
				})

				// Increment quota with token count
				if tokenCount > 0 {
					quotaService.IncrementRequest(virtualKey, tokenCount)
				}
			}
		}()

		// Extract virtual key from the Authorization header early
		authHeader := c.GetHeader("Authorization")
		if extractedKey, ok := utils.ExtractVirtualKey(authHeader); ok {
			virtualKey = extractedKey
			c.Set("virtual_key", virtualKey)
		}

		// Continue with request processing
		c.Next()
	}
}
