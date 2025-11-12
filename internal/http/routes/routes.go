package routes

import (
	"qualifire-home-assignment/internal/http/controllers"

	"github.com/gin-gonic/gin"
)

func HandleRequests() *gin.Engine {
	r := gin.Default()
	r.Use(controllers.Recovery())
	r.POST("/chat/completions", controllers.ChatCompletions{}.RouteRequests)

	return r
}
