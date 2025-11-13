package tests

import (
	"qualifire-home-assignment/internal/services"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetQuotaService_Singleton(t *testing.T) {
	service1 := services.GetQuotaService()
	service2 := services.GetQuotaService()

	assert.Same(t, service1, service2, "GetQuotaService should return the same instance")
}

func TestQuotaService_CheckQuota_Success(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()
	service.SetLimits(100, 100000, time.Hour)

	allowed, reason := service.CheckQuota("test-key-1")

	assert.True(t, allowed)
	assert.Empty(t, reason)
}

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

func TestQuotaService_GetUsage_NonExistentKey(t *testing.T) {
	service := services.GetQuotaService()
	service.Reset()

	requests, tokens := service.GetUsage("nonexistent-key")

	assert.Equal(t, 0, requests)
	assert.Equal(t, 0, tokens)
}

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
