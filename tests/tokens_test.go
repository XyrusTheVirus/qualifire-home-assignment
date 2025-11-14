package tests

import (
	"qualifire-home-assignment/internal/models"
	"qualifire-home-assignment/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEstimateTokens_Empty verifies token estimation for empty input
// Expected behavior: returns 0 tokens for empty string
func TestEstimateTokens_Empty(t *testing.T) {
	tokens := utils.EstimateTokens("")
	assert.Equal(t, 0, tokens)
}

// TestEstimateTokens_ShortText verifies token estimation for short input
// Tests that a 5-character word "Hello" results in 1 token (5/4 = 1.25 rounded down)
func TestEstimateTokens_ShortText(t *testing.T) {
	// "Hello" = 5 chars / 4 = 1.25 = 1 token
	tokens := utils.EstimateTokens("Hello")
	assert.Equal(t, 1, tokens)
}

// TestEstimateTokens_MediumText verifies token estimation for medium-length input
// Tests that a 19-character sentence results in 4 tokens (19/4 = 4.75 rounded down)
func TestEstimateTokens_MediumText(t *testing.T) {
	// "Hello, how are you?" = 19 chars / 4 = 4.75 = 4 tokens
	tokens := utils.EstimateTokens("Hello, how are you?")
	assert.Equal(t, 4, tokens)
}

// TestEstimateTokens_LongText verifies token estimation for long input
// Tests that a 100-character text results in approximately 25 tokens with delta of 2
func TestEstimateTokens_LongText(t *testing.T) {
	// 100 characters = 25 tokens
	text := "This is a longer text that contains exactly one hundred characters to test the token estimation."
	tokens := utils.EstimateTokens(text)
	assert.InDelta(t, 25, tokens, 2, "Should be approximately 25 tokens")
}

// TestEstimateTokens_Unicode verifies token estimation for Unicode text
// Tests that text containing Unicode characters produces non-zero token count
func TestEstimateTokens_Unicode(t *testing.T) {
	// Unicode characters should be counted correctly
	tokens := utils.EstimateTokens("Hello 世界")
	assert.Greater(t, tokens, 0)
}

// TestCalculateRequestTokens_SingleMessage verifies token calculation for single message
// Tests that token components (role=1, content=1, formatting=3, base=3) sum to 8 tokens
func TestCalculateRequestTokens_SingleMessage(t *testing.T) {
	messages := []models.Message{
		{Role: "user", Content: "Hello"},
	}

	tokens := utils.CalculateRequestTokens(messages, "gpt-3.5-turbo")

	// 1 (role) + 1 (content) + 3 (formatting) + 3 (base) = 8 tokens
	assert.Equal(t, 8, tokens)
}

// TestCalculateRequestTokens_MultipleMessages verifies token calculation for multiple messages
// Tests that three messages plus base tokens sum to 31 tokens total
func TestCalculateRequestTokens_MultipleMessages(t *testing.T) {
	messages := []models.Message{
		{Role: "system", Content: "You are a helpful assistant"},
		{Role: "user", Content: "Hello, how are you?"},
		{Role: "assistant", Content: "I'm doing great, thank you!"},
	}

	tokens := utils.CalculateRequestTokens(messages, "gpt-3.5-turbo")

	// Message 1: 1 + 6 + 3 = 10
	// Message 2: 1 + 4 + 3 = 8
	// Message 3: 1 + 6 + 3 = 10
	// Base: 3
	// Total: 31 tokens
	assert.Equal(t, 31, tokens)
}

// TestCalculateRequestTokens_EmptyMessage verifies token calculation for empty message
// Tests that empty message components (role=1, content=0, formatting=3, base=3) sum to 7 tokens
func TestCalculateRequestTokens_EmptyMessage(t *testing.T) {
	messages := []models.Message{
		{Role: "user", Content: ""},
	}

	tokens := utils.CalculateRequestTokens(messages, "gpt-3.5-turbo")

	// 1 (role) + 0 (content) + 3 (formatting) + 3 (base) = 7 tokens
	assert.Equal(t, 7, tokens)
}

// TestCalculateRequestTokens_LongConversation verifies token calculation for multi-turn conversation
// Tests that a four-message conversation results in token count between 20 and 100
func TestCalculateRequestTokens_LongConversation(t *testing.T) {
	messages := []models.Message{
		{Role: "user", Content: "What is the capital of France?"},
		{Role: "assistant", Content: "The capital of France is Paris."},
		{Role: "user", Content: "What about Germany?"},
		{Role: "assistant", Content: "The capital of Germany is Berlin."},
	}

	tokens := utils.CalculateRequestTokens(messages, "gpt-3.5-turbo")

	// Should be a reasonable number based on content
	assert.Greater(t, tokens, 20)
	assert.Less(t, tokens, 100)
}

// TestCalculateResponseTokens_Empty verifies token calculation for empty response
// Tests that empty response results in 3 tokens for formatting overhead
func TestCalculateResponseTokens_Empty(t *testing.T) {
	tokens := utils.CalculateResponseTokens("")
	assert.Equal(t, 3, tokens) // Just formatting overhead
}

// TestCalculateResponseTokens_ShortResponse verifies token calculation for minimal response
// Tests that "Yes" results in 4 tokens total (1 content + 3 formatting)
func TestCalculateResponseTokens_ShortResponse(t *testing.T) {
	tokens := utils.CalculateResponseTokens("Yes")
	// 0 (3 chars / 4) + 3 (formatting) = 3 tokens (minimum 1 for content)
	assert.Equal(t, 4, tokens)
}

// TestCalculateResponseTokens_LongResponse verifies token calculation for detailed response
// Tests that long response results in approximately 35 tokens with delta of 5
func TestCalculateResponseTokens_LongResponse(t *testing.T) {
	response := "This is a detailed response that contains multiple sentences and provides comprehensive information about the topic at hand."
	tokens := utils.CalculateResponseTokens(response)

	// ~130 chars / 4 = ~32 tokens + 3 formatting = ~35 tokens
	assert.InDelta(t, 35, tokens, 5)
}

// TestTokenEstimation_Accuracy verifies token estimation accuracy across multiple cases
// Tests various text lengths (4, 8, 16, 32 chars) against expected token ranges
func TestTokenEstimation_Accuracy(t *testing.T) {
	testCases := []struct {
		name          string
		text          string
		expectedRange [2]int // min, max expected tokens
	}{
		{"4 chars", "test", [2]int{1, 2}},
		{"8 chars", "test123", [2]int{1, 3}},
		{"16 chars", "test1234test5678", [2]int{3, 5}},
		{"32 chars", "test1234test5678test1234test5678", [2]int{7, 9}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := utils.EstimateTokens(tc.text)
			assert.GreaterOrEqual(t, tokens, tc.expectedRange[0])
			assert.LessOrEqual(t, tokens, tc.expectedRange[1])
		})
	}
}
