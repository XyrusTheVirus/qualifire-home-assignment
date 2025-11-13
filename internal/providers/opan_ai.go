package providers

import (
	"context"
	"net/http"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/errors"
	"qualifire-home-assignment/internal/models"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared/constant"
)

// OpenAI provider implementation
type OpenAI struct {
	ProviderBase
}

func (o OpenAI) SendRequest() *Response {
	var err error
	var resp *openai.ChatCompletion

	httpClient := o.GetHttpClient()
	// Create OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(o.Request.ApiKey),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(configs.EnvInt("PROVIDER_MAX_RETRIES", "0")),
	)

	messages := toOpenAIMessages(o.Request)
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    o.Request.Model,
	}

	resp, err = client.Chat.Completions.New(context.Background(), params)
	if err != nil {
		panic(errors.ApiProvider{}.GetError(err.Error(), http.StatusInternalServerError))
	}

	return &Response{
		Choices: []Message{
			{Role: "user", Content: resp.Choices[0].Message.Content},
		},
	}
}

// toOpenAIMessages converts ProxyRequest messages to OpenAI ChatCompletionMessageParamUnion format
func toOpenAIMessages(req models.ProxyRequest) []openai.ChatCompletionMessageParamUnion {
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages))
	for _, m := range req.Messages {
		// Map roles to OpenAI message types
		switch m.Role {
		case "system":
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Role: constant.System(m.Role),
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String(m.Content),
					},
				},
			})
		case "user":
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Role: constant.User(m.Role),
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(m.Content),
					},
				},
			})
		case "assistant":
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Role: constant.Assistant(m.Role),
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String(m.Content),
					},
				},
			})
		case "developer":
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfDeveloper: &openai.ChatCompletionDeveloperMessageParam{
					Role: constant.Developer(m.Role),
					Content: openai.ChatCompletionDeveloperMessageParamContentUnion{
						OfString: openai.String(m.Content),
					},
				},
			})
		case "tool":
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfTool: &openai.ChatCompletionToolMessageParam{
					Role: constant.Tool(m.Role),
					Content: openai.ChatCompletionToolMessageParamContentUnion{
						OfString: openai.String(m.Content),
					},
				},
			})
		default:
			continue
		}
	}

	return messages
}
