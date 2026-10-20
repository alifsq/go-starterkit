package middleware

import (
	"go-starterkit/internal/package/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()
		if log == nil {
			return
		}
		latency := time.Since(start)
		status := c.Writer.Status()

		if query != "" {
			path = path + "?" + query
		}

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Error("Gin Internal Error", zap.Error(e.Err))
			}
		}

		if status >= 500 {
			log.Error("HTTP Server Error", fields...)
		} else if status >= 400 {
			log.Warn("HTTP Client Error", fields...)
		} else {
			log.Info("HTTP Request Success", fields...)
		}
	}
}
