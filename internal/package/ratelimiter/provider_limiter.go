package ratelimiter

import "github.com/google/wire"

var RateLimitSet = wire.NewSet(
	NewLimiterEngine,
)
