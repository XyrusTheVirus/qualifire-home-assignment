package middleware

import (
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/services"

	"github.com/gin-gonic/gin"
)

// QuotaMiddleware checks if the virtual key has quota available
// Note: Virtual key must be set in context by validator before this runs
func QuotaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check quota after validation (virtual key is set in context by validator)
		virtualKeyInterface, exists := c.Get("virtual_key")
		if !exists {
			return
		}

		virtualKey := virtualKeyInterface.(string)
		quotaService := services.GetQuotaService()

		// Check quota
		allowed, reason := quotaService.CheckQuota(virtualKey)
		if !allowed {
			err := errors.GetError("QUOTA_EXCEEDED", reason, http.StatusTooManyRequests)
			c.JSON(err.StatusCode, err.ToGin())
			c.Abort()
			return
		}
	}
}
