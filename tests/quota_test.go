package tests

import (
	"qualifire-home-assignment/internal/services"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGetQuotaService_Singleton ensures that the GetQuotaService function implements
// the singleton pattern by returning the same instance for multiple calls
func TestGetQuotaService_Singleton(t *testing.T) {
	service1 := services.GetQuotaService()
	service2 := services.GetQuotaService()

	assert.Same(t, service1, service2, "GetQuotaService should return the same instance")
}

// TestQuotaService_CheckQuota_Success verifies that quota checks pass when
// the usage is within the configured request and token limits
func TestQuotaService_CheckQuota_Success(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(100, 100000, time.Hour)

	allowed, reason := service.CheckQuota("test-key-1")

	assert.True(t, allowed)
	assert.Empty(t, reason)
}

// TestQuotaService_CheckQuota_RequestLimitExceeded verifies that quota checks fail
// when the number of requests exceeds the configured limit
func TestQuotaService_CheckQuota_RequestLimitExceeded(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(5, 100000, time.Hour)

	// Increment to limit
	for i := 0; i < 5; i++ {
		service.IncrementRequest("test-key-1", 10)
	}

	allowed, reason := service.CheckQuota("test-key-1")

	assert.False(t, allowed)
	assert.Equal(t, "request quota exceeded", reason)
}

// TestQuotaService_CheckQuota_TokenLimitExceeded verifies that quota checks fail
// when the number of tokens exceeds the configured limit
func TestQuotaService_CheckQuota_TokenLimitExceeded(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(100, 1000, time.Hour)

	// Increment tokens to limit
	service.IncrementRequest("test-key-1", 1001)

	allowed, reason := service.CheckQuota("test-key-1")

	assert.False(t, allowed)
	assert.Equal(t, "token quota exceeded", reason)
}

// TestQuotaService_IncrementRequest verifies that the service correctly tracks
// both request count and token usage when incrementing requests
func TestQuotaService_IncrementRequest(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(100, 100000, time.Hour)

	service.IncrementRequest("test-key-1", 50)
	service.IncrementRequest("test-key-1", 30)

	requests, tokens := service.GetUsage("test-key-1")

	assert.Equal(t, 2, requests)
	assert.Equal(t, 80, tokens)
}

// TestQuotaService_GetUsage_NonExistentKey verifies that GetUsage returns
// zero values for both requests and tokens when using a non-existent key
func TestQuotaService_GetUsage_NonExistentKey(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()

	requests, tokens := service.GetUsage("nonexistent-key")

	assert.Equal(t, 0, requests)
	assert.Equal(t, 0, tokens)
}

// TestQuotaService_WindowReset verifies that usage counters are automatically
// reset to zero after the configured time window expires
func TestQuotaService_WindowReset(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(100, 100000, 100*time.Millisecond)

	service.IncrementRequest("test-key-1", 50)

	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)

	requests, tokens := service.GetUsage("test-key-1")

	assert.Equal(t, 0, requests)
	assert.Equal(t, 0, tokens)
}

// TestQuotaService_ConcurrentAccess verifies that the service handles concurrent
// requests correctly by maintaining accurate counts under high concurrency
func TestQuotaService_ConcurrentAccess(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(1000, 100000, time.Hour)

	var wg sync.WaitGroup
	concurrency := 100

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.IncrementRequest("test-key-concurrent", 10)
		}()
	}

	wg.Wait()

	requests, tokens := service.GetUsage("test-key-concurrent")

	assert.Equal(t, concurrency, requests)
	assert.Equal(t, concurrency*10, tokens)
}

// TestQuotaService_MultipleKeys verifies that the service correctly maintains
// separate usage tracking for different virtual keys
func TestQuotaService_MultipleKeys(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(100, 100000, time.Hour)

	service.IncrementRequest("key-1", 10)
	service.IncrementRequest("key-2", 20)
	service.IncrementRequest("key-1", 15)

	requests1, tokens1 := service.GetUsage("key-1")
	requests2, tokens2 := service.GetUsage("key-2")

	assert.Equal(t, 2, requests1)
	assert.Equal(t, 25, tokens1)
	assert.Equal(t, 1, requests2)
	assert.Equal(t, 20, tokens2)
}

// TestQuotaService_Reset verifies that the Reset function completely clears
// all usage data from the service
func TestQuotaService_Reset(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(100, 100000, time.Hour)

	service.IncrementRequest("test-key-1", 50)
	service.Reset()

	requests, tokens := service.GetUsage("test-key-1")

	assert.Equal(t, 0, requests)
	assert.Equal(t, 0, tokens)
}
