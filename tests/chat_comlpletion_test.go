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

func setupTest() {
	os.Setenv("IS_TEST", "1")
	gin.SetMode(gin.TestMode)
	configs.LoadConfig()
}

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
	req.Header.Set("Authorization", "Bearer test-key-1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}
	result := cc.Validate(c)

	assert.NotNil(t, result)
	assert.Equal(t, 200, w.Code)
}

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
	req.Header.Set("Authorization", "Bearer test-key-1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cc := validators.ChatCompletion{}
	result := cc.Validate(c)

	assert.Nil(t, result)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

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

func TestGetKeyInfo_Success(t *testing.T) {
	setupTest()

	cc := &validators.ChatCompletion{
		AuthToken: "Bearer test-key-1",
	}

	keyInfo, virtualKey := validators.GetKeyInfo(cc)

	assert.NotNil(t, keyInfo)
	assert.Equal(t, "test-key-1", virtualKey)
	assert.Contains(t, keyInfo, "provider")
	assert.Contains(t, keyInfo, "api_key")
}

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

func TestGetProxyRequest_Success(t *testing.T) {
	setupTest()

	cc := &validators.ChatCompletion{
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
		Model:     "gpt-3.5-turbo",
		AuthToken: "Bearer test-key-1",
	}

	result := validators.GetProxyRequest(cc)

	assert.NotNil(t, result)
}
