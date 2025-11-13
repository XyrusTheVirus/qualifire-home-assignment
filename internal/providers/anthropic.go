package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/errors"
)

const (
	VERSION      = "v1"
	MESSAGES_URI = "messages"
)

// Anthropic provider implementation
type Anthropic struct {
	ProviderBase
}

// SendRequest sends a request to the Anthropic provider and returns the response
func (a Anthropic) SendRequest() *Response {
	var err error
	var res *http.Response
	var req *http.Request

	body := map[string]interface{}{
		"model":      a.Request.Model,
		"max_tokens": 200,
		"messages":   a.Request.Messages,
	}
	b, _ := json.Marshal(body)
	httpClient := a.GetHttpClient()
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/%s/%s", configs.Env("ANTHROPIC_ENDPOINT", ""), VERSION, MESSAGES_URI), bytes.NewBuffer(b))
	if err != nil {
		return nil
	}
	req.Header.Set("x-api-key", a.Request.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", a.Request.Model)

	res, err = httpClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(fmt.Sprintf("Failed to close Anthropic response body: %s", err.Error()))
		}
	}(res.Body)
	if err != nil {
		panic(errors.ApiProvider{}.GetError(err.Error(), res.StatusCode))
	}

	var result map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode Anthropic response: %s", err.Error()))
	}

	if res.StatusCode != http.StatusOK {
		panic(errors.ApiProvider{}.GetError(result["error"].(map[string]interface{})["message"].(string), res.StatusCode))
	}

	return &Response{
		Choices: []Message{
			{Role: result["role"].(string), Content: result["content"].(string)},
		},
	}
}
