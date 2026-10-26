package middleware

import (
	"context"
	"fmt"
	"go-starterkit/internal/package/ratelimiter"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitMiddleware struct {
	limiterEngine RateLimitEngine
}

type RateLimitEngine interface {
	CheckLimit(context.Context, string, string) (*ratelimiter.LimitResult, error)
}

func NewLimiteMiddleware(limiterEngine RateLimitEngine) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiterEngine: limiterEngine,
	}
}

func (m *RateLimitMiddleware) LimitByIP(policyName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.ClientIP()
		res, err := m.CheckLimit(c.Request.Context(), policyName, identifier)
		if err != nil {
			c.Error(fmt.Errorf("rate limit check failed: %w", err))
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code":    http.StatusServiceUnavailable,
				"message": "Rate limiting is temporarily unavailable",
			})
			return
		}
		setRateLimitHeaders(c, policyName, res)

		if !res.Allowed {
			abortRateLimited(c, res)
			return
		}

		c.Next()
	}
}

// LimitByUser digunakan untuk endpoint privat berdasarkan User ID (biasanya diambil dari context JWT)
func (m *RateLimitMiddleware) LimitByUser(policyName string, getUserIDFunc func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserIDFunc(c)
		if userID == "" {
			// Fallback ke IP jika User ID tidak ditemukan di context
			userID = c.ClientIP()
		}
		res, err := m.CheckLimit(c.Request.Context(), policyName, userID)
		if err != nil {
			c.Error(fmt.Errorf("rate limit check failed: %w", err))
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code":    http.StatusServiceUnavailable,
				"message": "Rate limiting is temporarily unavailable",
			})
			return
		}

		setRateLimitHeaders(c, policyName, res)

		if !res.Allowed {
			abortRateLimited(c, res)
			return
		}
		c.Next()
	}
}

func (m *RateLimitMiddleware) CheckLimit(ctx context.Context, policyName, identifier string) (*ratelimiter.LimitResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return m.limiterEngine.CheckLimit(ctx, policyName, identifier)
}

func setRateLimitHeaders(c *gin.Context, policyName string, result *ratelimiter.LimitResult) {
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", ratelimiter.Policies[policyName].MaxRequests))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
}

func abortRateLimited(c *gin.Context, result *ratelimiter.LimitResult) {
	retryAfter := result.RetryAfter
	if retryAfter < 1 {
		retryAfter = 1
	}
	c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"code":    http.StatusTooManyRequests,
		"message": "Too many requests. Please try again later.",
	})
}
