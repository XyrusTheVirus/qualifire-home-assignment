package providers

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	ProviderBase
}

func (o OpenAI) SendRequest() *Response {
	client := openai.NewClient(o.Request.ApiKey)
	log.Println("Hi:", client)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o, // or GPT4Turbo / GPT3Dot5Turbo
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: "Write a short haiku about Go and AI."},
			},
		},
	)
	if err != nil {
		log.Fatalf("OpenAI request failed: %v", err)
	}

	fmt.Println("OpenAI Response:", resp.Choices[0].Message.Content)
	return &Response{
		Choices: []Message{
			{Role: "user", Content: resp.Choices[0].Message.Content},
		},
	}
}
