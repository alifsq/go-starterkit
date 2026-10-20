package auth

import (
	"sync"
	"time"
)

type inMemoryRateLimiter struct {
	mu     sync.Mutex
	entries map[string][]time.Time
}

func newInMemoryRateLimiter() *inMemoryRateLimiter {
	return &inMemoryRateLimiter{entries: make(map[string][]time.Time)}
}

func (l *inMemoryRateLimiter) Allow(key string, window time.Duration, max int) bool {
	if max <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-window)
	filtered := l.entries[key][:0]
	for _, ts := range l.entries[key] {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	if len(filtered) >= max {
		l.entries[key] = filtered
		return false
	}
	filtered = append(filtered, now)
	l.entries[key] = filtered
	return true
}
