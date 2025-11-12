package models

type ProxyRequest struct {
	Provider string    `json:"provider"`
	ApiKey   string    `json:"api_key"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
