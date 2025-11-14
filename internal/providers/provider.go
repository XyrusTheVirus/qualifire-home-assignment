package providers

import (
	"net/http"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/models"
	"qualifire-home-assignment/internal/transports"
	"time"
)

// Provider defines the interface for all providers
type Provider interface {
	// SendRequest sends a request to the provider and returns the response
	SendRequest() *Response
}

// Response represents a generic response from a provider
type Response struct {
	Choices      []Message `json:"choices"`
	TokensUsed   int       `json:"tokens_used,omitempty"`
}

// ProviderBase provides common fields for all providers
type ProviderBase struct {
	Request models.ProxyRequest
}

// Message represents a single message in the chat completion request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	OPENAI    = "openai"
	ANTHROPIC = "anthropic"
)

// Factory creates a Provider based on the ProxyRequest provider field
func Factory(req models.ProxyRequest) Provider {
	switch req.Provider {
	case OPENAI:
		return OpenAI{ProviderBase{req}}
	case ANTHROPIC:
		return Anthropic{ProviderBase{req}}
	default:
		panic("unknown provider")
	}
}

// GetHttpClient returns a configured HTTP client with logging and timeout
func (p ProviderBase) GetHttpClient() *http.Client {
	return &http.Client{
		Transport: &transports.LoggingTransport{Base: &http.Transport{
			TLSHandshakeTimeout: time.Duration(configs.EnvInt("TLS_HANDSHAKE_TIMEOUT", "30")) * time.Second},
			Req: p.Request,
		},
		Timeout: time.Duration(configs.EnvInt("PROVIDER_REQUEST_TIMEOUT", "30")) * time.Second,
	}
}
