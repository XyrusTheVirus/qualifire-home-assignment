package tests

import (
	"qualifire-home-assignment/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVirtualKey_ValidBearer(t *testing.T) {
	authHeader := "Bearer test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.True(t, ok)
	assert.Equal(t, "test-key-123", virtualKey)
}

func TestExtractVirtualKey_ValidBearerWithHyphens(t *testing.T) {
	authHeader := "Bearer test-key-with-multiple-hyphens"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.True(t, ok)
	assert.Equal(t, "test-key-with-multiple-hyphens", virtualKey)
}

func TestExtractVirtualKey_ValidBearerWithUnderscores(t *testing.T) {
	authHeader := "Bearer test_key_123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.True(t, ok)
	assert.Equal(t, "test_key_123", virtualKey)
}

func TestExtractVirtualKey_EmptyHeader(t *testing.T) {
	authHeader := ""

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_InvalidFormat_NoBearer(t *testing.T) {
	authHeader := "test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_InvalidFormat_WrongPrefix(t *testing.T) {
	authHeader := "Basic test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_InvalidFormat_NoSpace(t *testing.T) {
	authHeader := "Bearertest-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_InvalidFormat_ExtraSpaces(t *testing.T) {
	authHeader := "Bearer  test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_InvalidFormat_SpecialChars(t *testing.T) {
	authHeader := "Bearer test@key#123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_InvalidFormat_OnlyBearer(t *testing.T) {
	authHeader := "Bearer "

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_InvalidFormat_OnlyKey(t *testing.T) {
	authHeader := " test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

func TestExtractVirtualKey_LowercaseBearer(t *testing.T) {
	authHeader := "bearer test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}
