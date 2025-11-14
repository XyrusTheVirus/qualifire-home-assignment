package models

// ProxyRequest represents a request to be proxied to an external provider
type ProxyRequest struct {
	Provider   string    `json:"provider"`
	ApiKey     string    `json:"api_key"`
	Messages   []Message `json:"messages"`
	Model      string    `json:"model"`
	VirtualKey string    `json:"virtual_key"`
}

// Message represents a single message in the chat completion request
type Message struct {
	Role    string `json:"role" binding:"required,oneof=system user assistant developer tool"`
	Content string `json:"content" binding:"required,max=255"`
}
