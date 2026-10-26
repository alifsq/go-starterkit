package middleware

import (
	"go-starterkit/internal/package/ratelimiter"

	"github.com/google/wire"
)

var RateLimitSet = wire.NewSet(
	ratelimiter.NewLimiterEngine,
	NewLimiteMiddleware,
	wire.Bind(new(RateLimitEngine), new(*ratelimiter.LimiterEngine)),
)
