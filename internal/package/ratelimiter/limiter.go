package ratelimiter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Policy mendefinisikan aturan Token Bucket per endpoint/aksi
type Policy struct {
	MaxRequests int           // Kapasitas maksimal ember (Burst / Max Hits)
	Window      time.Duration // Durasi waktu untuk mengisi penuh kembali ember
}


type LimiterEngine struct {
	rdb *redis.Client
}

func NewLimiterEngine(rdb *redis.Client) *LimiterEngine {
	return &LimiterEngine{rdb: rdb}
}

type LimitResult struct {
	Allowed    bool
	Remaining  int
	RetryAfter int64
}

// CheckLimit mengevaluasi policy berdasarkan identifier (IP, User ID, atau email).
func (e *LimiterEngine) CheckLimit(ctx context.Context, policyName string, identifier string) (*LimitResult, error) {
	policy, exists := Policies[policyName]
	if !exists {
		return nil, fmt.Errorf("unknown rate limit policy %q", policyName)
	}
	if policy.MaxRequests <= 0 || policy.Window <= 0 || identifier == "" {
		return nil, fmt.Errorf("invalid rate limit policy or identifier")
	}

	identifierHash := sha256.Sum256([]byte(identifier))
	key := fmt.Sprintf("ratelimit:tb:%s:%s", policyName, hex.EncodeToString(identifierHash[:]))

	windowSecs := policy.Window.Seconds()
	capacity := float64(policy.MaxRequests)
	refillRate := capacity / windowSecs // Kecepatan refill koin per detik

	requested := 1.0
	ttl := int(windowSecs)

	res, err := tokenBucketScript.Run(ctx, e.rdb, []string{key}, capacity, refillRate, requested, ttl).Result()
	if err != nil {
		return nil, err
	}

	results, ok := res.([]any)
	if !ok || len(results) < 3 {
		return nil, fmt.Errorf("invalid response from redis lua script")
	}

	allowedValue, ok := results[0].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid allowed value from redis lua script")
	}
	remainingValue, ok := results[1].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid remaining value from redis lua script")
	}
	retryAfter, ok := results[2].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid retry-after value from redis lua script")
	}

	return &LimitResult{
		Allowed:    allowedValue == 1,
		Remaining:  int(remainingValue),
		RetryAfter: retryAfter,
	}, nil
}
