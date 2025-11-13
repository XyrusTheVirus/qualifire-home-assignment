package routes

import (
	"net/http"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/controllers"
	"qualifire-home-assignment/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func HandleRequests() *gin.Engine {
	r := gin.Default()
	if !configs.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(controllers.Recovery())

	g := r.Group("/chat")
	g.Use(middleware.QuotaMiddleware())
	{
		r.POST("/completions", controllers.ChatCompletions{}.RouteRequests)
	}

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/metrics", controllers.Metrics{}.GetMetrics)

	return r
}
