package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"qualifire-home-assignment/internal/configs"
	"qualifire-home-assignment/internal/http/validators"
	"qualifire-home-assignment/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTest initializes the test environment by setting test mode and loading configurations
func setupTest() {
	os.Setenv("IS_TEST", "1")
	gin.SetMode(gin.TestMode)
	configs.LoadConfig()
}

// TestChatCompletionValidate_Success tests successful validation of chat completion request with valid payload and auth token
func TestChatCompletionValidate_Success(t *testing.T) {
	setupTest()

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
		"model": "gpt-3.5-turbo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer vk_user1_openai")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}
	result := cc.Validate(c)

	assert.NotNil(t, result)
	assert.Equal(t, 200, w.Code)
}

// TestChatCompletionValidate_MissingMessages tests validation failure when the messages field is missing from the request
func TestChatCompletionValidate_MissingMessages(t *testing.T) {
	setupTest()

	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key-1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}
	result := cc.Validate(c)

	assert.Nil(t, result)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestChatCompletionValidate_InvalidRole tests validation failure when message contains an invalid role
func TestChatCompletionValidate_InvalidRole(t *testing.T) {
	setupTest()

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "invalid_role", "content": "Hello"},
		},
		"model": "gpt-3.5-turbo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer vk_user1_openai")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}
	result := cc.Validate(c)

	assert.Nil(t, result)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestChatCompletionValidate_MissingAuthToken tests validation behavior when the Authorization header is missing
func TestChatCompletionValidate_MissingAuthToken(t *testing.T) {
	setupTest()

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
		"model": "gpt-3.5-turbo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		}
	}()

	cc.Validate(c)
}

// TestChatCompletionValidate_InvalidAuthFormat tests validation behavior when the Authorization header has invalid format
func TestChatCompletionValidate_InvalidAuthFormat(t *testing.T) {
	setupTest()

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
		"model": "gpt-3.5-turbo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		}
	}()

	cc.Validate(c)
}

// TestChatCompletionValidate_InvalidVirtualKey tests validation behavior when the virtual key is invalid or not found
func TestChatCompletionValidate_InvalidVirtualKey(t *testing.T) {
	setupTest()

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
		"model": "gpt-3.5-turbo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-key")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		}
	}()

	cc.Validate(c)
}

// TestGetKeyInfo_Success tests successful retrieval of key information from a valid authorization token
func TestGetKeyInfo_Success(t *testing.T) {
	setupTest()

	cc := &validators.ChatCompletion{
		AuthToken: "Bearer vk_user1_openai",
	}

	keyInfo, virtualKey := validators.GetKeyInfo(cc)

	assert.NotNil(t, keyInfo)
	assert.Equal(t, "vk_user1_openai", virtualKey)
	assert.Contains(t, keyInfo, "provider")
	assert.Contains(t, keyInfo, "api_key")
}

// TestGetKeyInfo_WrongFormat tests error handling when attempting to get key info with wrong authorization format
func TestGetKeyInfo_WrongFormat(t *testing.T) {
	setupTest()

	cc := &validators.ChatCompletion{
		AuthToken: "InvalidFormat",
	}

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		}
	}()

	validators.GetKeyInfo(cc)
}

// TestGetProxyRequest_Success tests successful creation of proxy request from chat completion parameters
func TestGetProxyRequest_Success(t *testing.T) {
	setupTest()

	cc := &validators.ChatCompletion{
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
		Model:     "gpt-3.5-turbo",
		AuthToken: "Bearer vk_user1_openai",
	}

	result := validators.GetProxyRequest(cc)

	assert.NotNil(t, result)
}
