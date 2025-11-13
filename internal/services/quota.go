package services

import (
	"sync"
	"time"
)

// QuotaEntry tracks usage for a virtual key
type QuotaEntry struct {
	RequestCount int
	TokenUsage   int
	WindowStart  time.Time
	mu           sync.Mutex
}

// QuotaService manages quotas per virtual key
type QuotaService struct {
	quotas       sync.Map // map[string]*QuotaEntry
	maxRequests  int
	maxTokens    int
	windowPeriod time.Duration
}

var (
	quotaServiceInstance *QuotaService
	quotaServiceOnce     sync.Once
)

// GetQuotaService returns the singleton instance of QuotaService
func GetQuotaService() *QuotaService {
	quotaServiceOnce.Do(func() {
		quotaServiceInstance = &QuotaService{
			maxRequests:  100,                // 100 requests per hour
			maxTokens:    100000,             // 100k tokens per hour
			windowPeriod: time.Hour,
		}
	})
	return quotaServiceInstance
}

// CheckQuota verifies if the virtual key has quota available
func (q *QuotaService) CheckQuota(virtualKey string) (bool, string) {
	entry := q.getOrCreateEntry(virtualKey)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	now := time.Now()

	// Reset window if expired
	if now.Sub(entry.WindowStart) >= q.windowPeriod {
		entry.RequestCount = 0
		entry.TokenUsage = 0
		entry.WindowStart = now
	}

	// Check request quota
	if entry.RequestCount >= q.maxRequests {
		return false, "request quota exceeded"
	}

	// Check token quota
	if entry.TokenUsage >= q.maxTokens {
		return false, "token quota exceeded"
	}

	return true, ""
}

// IncrementRequest increments the request count for a virtual key
func (q *QuotaService) IncrementRequest(virtualKey string, tokens int) {
	entry := q.getOrCreateEntry(virtualKey)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	now := time.Now()

	// Reset window if expired
	if now.Sub(entry.WindowStart) >= q.windowPeriod {
		entry.RequestCount = 0
		entry.TokenUsage = 0
		entry.WindowStart = now
	}

	entry.RequestCount++
	entry.TokenUsage += tokens
}

// GetUsage returns the current usage for a virtual key
func (q *QuotaService) GetUsage(virtualKey string) (requests int, tokens int) {
	value, ok := q.quotas.Load(virtualKey)
	if !ok {
		return 0, 0
	}

	entry := value.(*QuotaEntry)
	entry.mu.Lock()
	defer entry.mu.Unlock()

	now := time.Now()

	// Return 0 if window expired
	if now.Sub(entry.WindowStart) >= q.windowPeriod {
		return 0, 0
	}

	return entry.RequestCount, entry.TokenUsage
}

// getOrCreateEntry gets or creates a quota entry for a virtual key
func (q *QuotaService) getOrCreateEntry(virtualKey string) *QuotaEntry {
	value, _ := q.quotas.LoadOrStore(virtualKey, &QuotaEntry{
		WindowStart: time.Now(),
	})
	return value.(*QuotaEntry)
}

// SetLimits allows updating quota limits (useful for testing)
func (q *QuotaService) SetLimits(maxRequests, maxTokens int, windowPeriod time.Duration) {
	q.maxRequests = maxRequests
	q.maxTokens = maxTokens
	q.windowPeriod = windowPeriod
}

// Reset clears all quota entries (useful for testing)
func (q *QuotaService) Reset() {
	q.quotas.Range(func(key, value interface{}) bool {
		q.quotas.Delete(key)
		return true
	})
}
