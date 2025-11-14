package utils

import (
	"qualifire-home-assignment/internal/models"
	"unicode/utf8"
)

// EstimateTokens estimates the number of tokens in a text string
// Uses the approximation: 1 token ≈ 4 characters (OpenAI/Anthropic standard)
func EstimateTokens(text string) int {
	charCount := utf8.RuneCountInString(text)
	tokens := charCount / 4
	if tokens == 0 && charCount > 0 {
		tokens = 1 // Minimum 1 token for non-empty text
	}
	return tokens
}

// CalculateRequestTokens calculates total tokens for a chat completion request
// Includes tokens from all messages plus overhead for formatting
func CalculateRequestTokens(messages []models.Message, model string) int {
	totalTokens := 0

	// Count tokens in each message
	for _, msg := range messages {
		// Role tokens (typically 1 token for role)
		totalTokens += 1
		// Content tokens
		totalTokens += EstimateTokens(msg.Content)
		// Message formatting overhead (approximately 3-4 tokens per message)
		totalTokens += 3
	}

	// Model-specific base overhead
	// OpenAI models typically add 3 tokens, Anthropic similar
	totalTokens += 3

	return totalTokens
}

// CalculateResponseTokens calculates tokens in a response
func CalculateResponseTokens(content string) int {
	// Response content tokens
	tokens := EstimateTokens(content)
	// Response formatting overhead
	tokens += 3
	return tokens
}
