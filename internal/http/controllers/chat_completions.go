package controllers

import (
	"net/http"
	"qualifire-home-assignment/internal/http/validators"
	"qualifire-home-assignment/internal/models"
	"qualifire-home-assignment/internal/providers"

	"github.com/gin-gonic/gin"
)

type ChatCompletions struct {
}

func (cp ChatCompletions) RouteRequests(c *gin.Context) {
	req := validators.ChatCompletion{}.Validate(c)
	result := providers.Factory(req.(models.ProxyRequest)).SendRequest()
	Success(c, http.StatusOK, result)
}
