package tests

import (
	"qualifire-home-assignment/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExtractVirtualKey_ValidBearer verifies that ExtractVirtualKey correctly extracts a valid virtual key from the Bearer token.
func TestExtractVirtualKey_ValidBearer(t *testing.T) {
	authHeader := "Bearer test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.True(t, ok)
	assert.Equal(t, "test-key-123", virtualKey)
}

// TestExtractVirtualKey_ValidBearerWithHyphens
func TestExtractVirtualKey_ValidBearerWithHyphens(t *testing.T) {
	authHeader := "Bearer test-key-with-multiple-hyphens"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.True(t, ok)
	assert.Equal(t, "test-key-with-multiple-hyphens", virtualKey)
}

// TestExtractVirtualKey_ValidBearerWithUnderscores
func TestExtractVirtualKey_ValidBearerWithUnderscores(t *testing.T) {
	authHeader := "Bearer test_key_123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.True(t, ok)
	assert.Equal(t, "test_key_123", virtualKey)
}

// TestExtractVirtualKey_EmptyHeader verifies that ExtractVirtualKey returns false when the Authorization header is empty.
func TestExtractVirtualKey_EmptyHeader(t *testing.T) {
	authHeader := ""

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_InvalidFormat verifies that ExtractVirtualKey returns false when the Authorization header is invalid.
func TestExtractVirtualKey_InvalidFormat_NoBearer(t *testing.T) {
	authHeader := "test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_InvalidFormat_WrongPrefix verifies that ExtractVirtualKey returns false when the Authorization header has the wrong prefix.
func TestExtractVirtualKey_InvalidFormat_WrongPrefix(t *testing.T) {
	authHeader := "Basic test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_InvalidFormat_NoKey verifies that ExtractVirtualKey returns false when the Authorization header does not contain a key.
func TestExtractVirtualKey_InvalidFormat_NoSpace(t *testing.T) {
	authHeader := "Bearertest-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_InvalidFormat_ExtraSpaces
func TestExtractVirtualKey_InvalidFormat_ExtraSpaces(t *testing.T) {
	authHeader := "Bearer  test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_InvalidFormat_SpecialChars
func TestExtractVirtualKey_InvalidFormat_SpecialChars(t *testing.T) {
	authHeader := "Bearer test@key#123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_InvalidFormat_OnlyKey
func TestExtractVirtualKey_InvalidFormat_OnlyBearer(t *testing.T) {
	authHeader := "Bearer "

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_InvalidFormat_OnlyKey
func TestExtractVirtualKey_InvalidFormat_OnlyKey(t *testing.T) {
	authHeader := " test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}

// TestExtractVirtualKey_LowercaseBearer
func TestExtractVirtualKey_LowercaseBearer(t *testing.T) {
	authHeader := "bearer test-key-123"

	virtualKey, ok := utils.ExtractVirtualKey(authHeader)

	assert.False(t, ok)
	assert.Empty(t, virtualKey)
}
