package validators

import (
	"net/http"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/models"
	"qualifire-home-assignment/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ChatCompletion represents the chat completion request payload
type ChatCompletion struct {
	// From body
	Messages []models.Message `json:"messages" binding:"required,dive" validate:"dive,min=1,max=100"`
	Model    string           `json:"model" binding:"required"`

	// From headers
	AuthToken string `header:"Authorization" validate:"required"`
}

// Validate validates the ChatCompletion request
func (cc ChatCompletion) Validate(c *gin.Context) models.Model {
	if err := c.ShouldBindJSON(&cc); err != nil {
		// Iterate through all validation errors and return them as a map
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			errs := make(map[string]string)
			for _, e := range validationErrs {
				errs[e.Tag()] = e.Error()
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
			return nil
		}

	}

	if err := c.ShouldBindHeader(&cc); err != nil {
		panic(errors.Validation{}.GetError(err.Error(), http.StatusBadRequest))
	}

	return GetProxyRequest(&cc)
}

// GetProxyRequest maps ChatCompletion to ProxyRequest based on virtual keys configuration
func GetProxyRequest(cc *ChatCompletion) models.Model {
	keyInfo, virtualKey := GetKeyInfo(cc)
	return models.ProxyRequest{
		Provider:   keyInfo["provider"].(string),
		ApiKey:     keyInfo["api_key"].(string),
		Messages:   cc.Messages,
		Model:      cc.Model,
		VirtualKey: virtualKey,
	}
}

func GetKeyInfo(cc *ChatCompletion) (map[string]interface{}, string) {
	virtualKeys := configs.Config("virtual_keys", "")
	if virtualKeys != "" {
		virtualKey, ok := utils.ExtractVirtualKey(cc.AuthToken)
		// If the extraction fails, the Authorization header format is wrong
		if !ok {
			panic(errors.Validation{}.GetError("wrong authorization header format", http.StatusBadRequest))
		}

		// Checks whether the virtual key exists in the configurations
		if vk, ok := virtualKeys.(map[string]interface{})[virtualKey]; ok {
			return vk.(map[string]interface{}), virtualKey
		} else {
			panic(errors.Validation{}.GetError("wrong virtual key", http.StatusBadRequest))
		}
	} else {
		panic(errors.GetError("MISSING_CONFIGURATIONS", "virtual Keys weren't load to configurations", http.StatusInternalServerError))
	}
}
