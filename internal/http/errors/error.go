package errors // Contract is the interface for the errors
import (
	"qualifire-home-assignment/internal/configs"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type Contract interface {
	GetError(message string, status int) Error
}

// ValidationErrorResponse is a struct that defines the validation error response
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error is the error struct
type Error struct {
	Code       string
	Message    string
	Details    string
	StatusCode int
}

// GetDetails returns the details
func GetDetails() string {
	if configs.IsDevelopment() {
		return ""
	}

	return string(debug.Stack())
}

func GetError(code string, message string, status int) Error {
	return Error{
		Code:       code,
		Message:    message,
		Details:    GetDetails(),
		StatusCode: status,
	}
}

// ToGin returns the error in gin format
func (e Error) ToGin() gin.H {
	return gin.H{
		"code":    e.Code,
		"message": e.Message,
		"details": e.Details,
	}
}
