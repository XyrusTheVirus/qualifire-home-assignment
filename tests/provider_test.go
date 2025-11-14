package tests

import (
	"qualifire-home-assignment/internal/models"
	"qualifire-home-assignment/internal/providers"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFactory_OpenAI tests the Factory method with OpenAI provider configuration
// and verifies it returns the correct provider instance
func TestFactory_OpenAI(t *testing.T) {
	req := models.ProxyRequest{
		Provider: "openai",
		ApiKey:   "test-key",
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
		Model:      "gpt-3.5-turbo",
		VirtualKey: "test-key-1",
	}

	provider := providers.Factory(req)

	assert.NotNil(t, provider)
	assert.IsType(t, providers.OpenAI{}, provider)
}

// TestFactory_Anthropic tests the Factory method with Anthropic provider configuration
// and verifies it returns the correct provider instance
func TestFactory_Anthropic(t *testing.T) {
	req := models.ProxyRequest{
		Provider: "anthropic",
		ApiKey:   "test-key",
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
		Model:      "claude-3-opus-20240229",
		VirtualKey: "test-key-2",
	}

	provider := providers.Factory(req)

	assert.NotNil(t, provider)
	assert.IsType(t, providers.Anthropic{}, provider)
}

// TestFactory_UnknownProvider verifies that the Factory method panics
// when an unknown provider is specified
func TestFactory_UnknownProvider(t *testing.T) {
	req := models.ProxyRequest{
		Provider: "unknown",
		ApiKey:   "test-key",
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
		Model:      "test-model",
		VirtualKey: "test-key-1",
	}

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "unknown provider", r)
		}
	}()

	providers.Factory(req)
}

// TestProviderBase_GetHttpClient tests the GetHttpClient method of the ProviderBase
// and ensures it returns a properly configured HTTP client
func TestProviderBase_GetHttpClient(t *testing.T) {
	req := models.ProxyRequest{
		Provider: "openai",
		ApiKey:   "test-key",
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
		Model:      "gpt-3.5-turbo",
		VirtualKey: "test-key-1",
	}

	base := providers.ProviderBase{Request: req}
	client := base.GetHttpClient()

	assert.NotNil(t, client)
	assert.NotNil(t, client.Transport)
	assert.NotZero(t, client.Timeout)
}

// TestMessage_Structure verifies the Message struct fields
// are properly set and accessible
func TestMessage_Structure(t *testing.T) {
	msg := providers.Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	assert.Equal(t, "user", msg.Role)
	assert.Equal(t, "Hello, world!", msg.Content)
}

// TestResponse_Structure verifies the Response struct fields
// are properly set and contain the expected message data
func TestResponse_Structure(t *testing.T) {
	resp := providers.Response{
		Choices: []providers.Message{
			{Role: "assistant", Content: "Hello!"},
		},
	}

	assert.Len(t, resp.Choices, 1)
	assert.Equal(t, "assistant", resp.Choices[0].Role)
	assert.Equal(t, "Hello!", resp.Choices[0].Content)
}
