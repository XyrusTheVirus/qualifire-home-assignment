package controllers

import (
	"net/http"
	"qualifire-home-assignment/internal/services"

	"github.com/gin-gonic/gin"
)

// Metrics handles metrics requests
type Metrics struct{}

// GetMetrics returns usage statistics
func (m Metrics) GetMetrics(c *gin.Context) {
	metricsService := services.GetMetricsService()
	stats := metricsService.GetStats()
	Success(c, http.StatusOK, stats)
}
