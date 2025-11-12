package validators

import (
	"net/http"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/models"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ChatCompletion represents the chat completion request payload
type ChatCompletion struct {
	// From body
	Messages []Message `json:"messages" binding:"required,min=1,max=100,dive"`

	// From headers
	AuthToken string `header:"Authorization" validate:"required"`
}

// Message represents a single message in the chat completion request
type Message struct {
	Role    string `json:"role" binding:"required,oneof=admin user assistant"`
	Content string `json:"content" binding:"required,max=255"`
}

// VirtualKey represents the virtual key structure
type VirtualKey struct {
	Provider string `json:"provider"`
	ApiKey   string `json:"api_key"`
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
	r := ConvertDto(c, cc)
	return getProxyRequest(r)
}

// getProxyRequest maps ChatCompletion to ProxyRequest based on virtual keys configuration
func getProxyRequest(cc *ChatCompletion) models.Model {
	messages := make([]models.Message, 0)
	for _, msg := range cc.Messages {
		messages = append(messages, models.Message{Role: msg.Role, Content: msg.Content})
	}
	vk := getKeyInfo(cc)
	return models.ProxyRequest{
		Provider: vk["provider"].(string),
		ApiKey:   vk["api_key"].(string),
		Messages: messages,
	}
}

func getKeyInfo(cc *ChatCompletion) map[string]interface{} {
	virtualKeys := configs.Config("virtual_keys", "")
	if virtualKeys != "" {
		re := regexp.MustCompile(`^Bearer ([\w-]+)$`)
		matches := re.FindStringSubmatch(cc.AuthToken)
		// If the regex yields different result then 2, the Authorization header format is wrong
		if len(matches) != 2 {
			panic(errors.Validation{}.GetError("wrong authorization header format", http.StatusBadRequest))
		}

		// Checks whether the virtual key exists in the configurations
		if vk, ok := virtualKeys.(map[string]interface{})[matches[1]]; ok {
			return vk.(map[string]interface{})
		} else {
			panic(errors.Validation{}.GetError("wrong virtual key", http.StatusBadRequest))
		}
	} else {
		panic(errors.GetError("MISSING_CONFIGURATIONS", "virtual Keys weren't load to configurations", http.StatusInternalServerError))
	}
}
