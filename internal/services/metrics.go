package services

import (
	"sync"
	"time"
)

// RequestMetric tracks individual request metrics
type RequestMetric struct {
	Provider     string
	VirtualKey   string
	Duration     time.Duration
	Status       int
	Timestamp    time.Time
}

// MetricsService collects and aggregates usage statistics
type MetricsService struct {
	metrics      []RequestMetric
	mu           sync.RWMutex
	totalRequests int
	providerCounts map[string]int
	totalDuration time.Duration
}

var (
	metricsServiceInstance *MetricsService
	metricsServiceOnce     sync.Once
)

// GetMetricsService returns the singleton instance of MetricsService
func GetMetricsService() *MetricsService {
	metricsServiceOnce.Do(func() {
		metricsServiceInstance = &MetricsService{
			metrics:        make([]RequestMetric, 0),
			providerCounts: make(map[string]int),
		}
	})
	return metricsServiceInstance
}

// RecordRequest records a request metric
func (m *MetricsService) RecordRequest(metric RequestMetric) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics = append(m.metrics, metric)
	m.totalRequests++
	m.providerCounts[metric.Provider]++
	m.totalDuration += metric.Duration
}

// GetStats returns aggregated statistics
func (m *MetricsService) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	avgResponseTime := float64(0)
	if m.totalRequests > 0 {
		avgResponseTime = float64(m.totalDuration.Milliseconds()) / float64(m.totalRequests)
	}

	// Copy provider counts to avoid race conditions
	providerCounts := make(map[string]int)
	for k, v := range m.providerCounts {
		providerCounts[k] = v
	}

	return map[string]interface{}{
		"total_requests":          m.totalRequests,
		"requests_per_provider":   providerCounts,
		"average_response_time_ms": avgResponseTime,
	}
}

// Reset clears all metrics (useful for testing)
func (m *MetricsService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics = make([]RequestMetric, 0)
	m.totalRequests = 0
	m.providerCounts = make(map[string]int)
	m.totalDuration = 0
}
