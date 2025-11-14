package controllers

import (
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/http/validators"
	"qualifire-home-assignment/internal/models"
	"qualifire-home-assignment/internal/providers"
	"qualifire-home-assignment/internal/services"

	"github.com/gin-gonic/gin"
)

// ChatCompletions handles chat completion requests
type ChatCompletions struct{}

// RouteRequests handles the chat completion requests
func (cp ChatCompletions) RouteRequests(c *gin.Context) {
	req := validators.ChatCompletion{}.Validate(c)
	proxyReq := req.(models.ProxyRequest)

	// Store provider in context for metrics middleware
	c.Set("provider", proxyReq.Provider)

	// Check quota using virtual key
	virtualKey := proxyReq.VirtualKey
	quotaService := services.GetQuotaService()
	allowed, reason := quotaService.CheckQuota(virtualKey)
	if !allowed {
		err := errors.GetError("QUOTA_EXCEEDED", reason, http.StatusTooManyRequests)
		c.JSON(err.StatusCode, err.ToGin())
		return
	}

	result := providers.Factory(proxyReq).SendRequest()

	// Store token count in context for metrics middleware
	c.Set("token_count", result.TokensUsed)

	Success(c, http.StatusOK, result)
}
