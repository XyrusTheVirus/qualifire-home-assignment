package tests

import (
	"os"
	"qualifire-home-assignment/internal/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv_WithValue(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := configs.Env("TEST_KEY", "default")
	assert.Equal(t, "test_value", result)
}

func TestEnv_WithFallback(t *testing.T) {
	os.Unsetenv("NONEXISTENT_KEY")

	result := configs.Env("NONEXISTENT_KEY", "fallback_value")
	assert.Equal(t, "fallback_value", result)
}

func TestEnvInt_WithValue(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := configs.EnvInt("TEST_INT", "0")
	assert.Equal(t, 42, result)
}

func TestEnvInt_WithFallback(t *testing.T) {
	os.Unsetenv("NONEXISTENT_INT")

	result := configs.EnvInt("NONEXISTENT_INT", "10")
	assert.Equal(t, 10, result)
}

func TestEnvInt_InvalidValue(t *testing.T) {
	os.Setenv("INVALID_INT", "not_a_number")
	defer os.Unsetenv("INVALID_INT")

	result := configs.EnvInt("INVALID_INT", "5")
	assert.Equal(t, 0, result)
}

func TestIsDevelopment_Dev(t *testing.T) {
	os.Setenv("APP_ENV", "dev")
	defer os.Unsetenv("APP_ENV")

	result := configs.IsDevelopment()
	assert.True(t, result)
}

func TestIsDevelopment_NotDev(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	result := configs.IsDevelopment()
	assert.False(t, result)
}

func TestIsDevelopment_Empty(t *testing.T) {
	os.Unsetenv("APP_ENV")

	result := configs.IsDevelopment()
	assert.False(t, result)
}

func TestConfig_LoadAndRead(t *testing.T) {
	os.Setenv("IS_TEST", "1")
	defer os.Unsetenv("IS_TEST")

	configs.LoadConfig()

	value := configs.Config("virtual_keys", "")
	assert.NotNil(t, value)
}

func TestConfig_Fallback(t *testing.T) {
	result := configs.Config("nonexistent_key", "default_value")
	assert.NotNil(t, result)
}
