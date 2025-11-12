package providers

import "qualifire-home-assignment/internal/models"

type Provider interface {
	SendRequest() *Response
}

type Response struct {
	Choices []Message `json:"choices"`
}

type ProviderBase struct {
	Request models.ProxyRequest
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	OPENAI    = "openai"
	ANTHROPIC = "anthropic"
)

func Factory(req models.ProxyRequest) Provider {
	switch req.Provider {
	case OPENAI:
		return OpenAI{ProviderBase{req}}
	case ANTHROPIC:
		return Anthropic{}
	default:
		panic("unknown provider")
	}
}
