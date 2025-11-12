package controllers

import (
	"fmt"
	"log"
	"net/http"
	"qualifire-home-assignment/internal/http/errors"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Success is the success response
func Success(c *gin.Context, status int, results interface{}) {
	c.JSON(status, results)
}

// Recovery is the recovery middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovery error: %v", err)
				log.Printf("Panic recovery debug stack:" + string(debug.Stack()))
				switch e := err.(type) {
				case errors.Error:
					c.JSON(e.StatusCode, e.ToGin())
				default:
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": fmt.Sprintf("internal error: %v", err),
					})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
