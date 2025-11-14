package middleware

import (
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/services"
	"qualifire-home-assignment/internal/utils"

	"github.com/gin-gonic/gin"
)

// QuotaMiddleware checks if the virtual key has quota available
// Note: Virtual key must be set in context by validator before this runs
func QuotaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check quota after validation (validator sets virtual key in context)
		virtualKeyInterface, exists := utils.ExtractVirtualKey(c.GetHeader("Authorization"))
		if !exists {
			return
		}

		virtualKey := virtualKeyInterface
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
