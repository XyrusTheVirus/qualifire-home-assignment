package tests

import (
	"net/http/httptest"
	"qualifire-home-assignment/internal/http/validators"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestConvertDto_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	result := validators.ConvertDto(c, testData)

	assert.NotNil(t, result)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
}

func TestConvertDto_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := TestStruct{
		Name:  "",
		Email: "invalid-email",
	}

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		}
	}()

	validators.ConvertDto(c, testData)
}

func TestConvertDto_NotAStruct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := "not a struct"

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		}
	}()

	validators.ConvertDto(c, testData)
}
