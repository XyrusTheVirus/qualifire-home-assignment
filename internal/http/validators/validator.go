package validators

import (
	"qualifire-home-assignment/internal/models"

	"github.com/gin-gonic/gin"
)

type Validator interface {
	Validate(c *gin.Context) models.Model
}
