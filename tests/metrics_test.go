package tests

import (
	"qualifire-home-assignment/internal/services"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricsService_Singleton(t *testing.T) {
	service1 := services.GetMetricsService()
	service2 := services.GetMetricsService()

	assert.Same(t, service1, service2, "GetMetricsService should return the same instance")
}

func TestMetricsService_RecordRequest(t *testing.T) {
	service := services.GetMetricsService()
	service.Reset()

	metric := services.RequestMetric{
		Provider:   "openai",
		VirtualKey: "test-key",
		Duration:   100 * time.Millisecond,
		Status:     200,
		Timestamp:  time.Now(),
	}

	service.RecordRequest(metric)

	stats := service.GetStats()
	assert.Equal(t, 1, stats["total_requests"])
}

func TestMetricsService_GetStats_Empty(t *testing.T) {
	service := services.GetMetricsService()
	service.Reset()

	stats := service.GetStats()

	assert.Equal(t, 0, stats["total_requests"])
	assert.Equal(t, float64(0), stats["average_response_time_ms"])
	assert.NotNil(t, stats["requests_per_provider"])
}

func TestMetricsService_GetStats_MultipleProviders(t *testing.T) {
	service := services.GetMetricsService()
	service.Reset()

	service.RecordRequest(services.RequestMetric{
		Provider:  "openai",
		Duration:  100 * time.Millisecond,
		Status:    200,
		Timestamp: time.Now(),
	})

	service.RecordRequest(services.RequestMetric{
		Provider:  "anthropic",
		Duration:  200 * time.Millisecond,
		Status:    200,
		Timestamp: time.Now(),
	})

	service.RecordRequest(services.RequestMetric{
		Provider:  "openai",
		Duration:  150 * time.Millisecond,
		Status:    200,
		Timestamp: time.Now(),
	})

	stats := service.GetStats()
	providerCounts := stats["requests_per_provider"].(map[string]int)

	assert.Equal(t, 3, stats["total_requests"])
	assert.Equal(t, 2, providerCounts["openai"])
	assert.Equal(t, 1, providerCounts["anthropic"])
}

func TestMetricsService_GetStats_AverageResponseTime(t *testing.T) {
	service := services.GetMetricsService()
	service.Reset()

	service.RecordRequest(services.RequestMetric{
		Provider:  "openai",
		Duration:  100 * time.Millisecond,
		Status:    200,
		Timestamp: time.Now(),
	})

	service.RecordRequest(services.RequestMetric{
		Provider:  "openai",
		Duration:  200 * time.Millisecond,
		Status:    200,
		Timestamp: time.Now(),
	})

	stats := service.GetStats()
	avgTime := stats["average_response_time_ms"].(float64)

	assert.InDelta(t, 150.0, avgTime, 1.0)
}

func TestMetricsService_ConcurrentRecording(t *testing.T) {
	service := services.GetMetricsService()
	service.Reset()

	var wg sync.WaitGroup
	concurrency := 100

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			provider := "openai"
			if index%2 == 0 {
				provider = "anthropic"
			}
			service.RecordRequest(services.RequestMetric{
				Provider:  provider,
				Duration:  100 * time.Millisecond,
				Status:    200,
				Timestamp: time.Now(),
			})
		}(i)
	}

	wg.Wait()

	stats := service.GetStats()
	assert.Equal(t, concurrency, stats["total_requests"])
}

func TestMetricsService_Reset(t *testing.T) {
	service := services.GetMetricsService()
	service.Reset()

	service.RecordRequest(services.RequestMetric{
		Provider:  "openai",
		Duration:  100 * time.Millisecond,
		Status:    200,
		Timestamp: time.Now(),
	})

	service.Reset()

	stats := service.GetStats()
	assert.Equal(t, 0, stats["total_requests"])
	assert.Equal(t, float64(0), stats["average_response_time_ms"])
}
