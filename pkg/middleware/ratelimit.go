package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"go.uber.org/zap"
)

// RateLimitMiddleware creates a rate limiter middleware
func RateLimitMiddleware(logger *zap.Logger) gin.HandlerFunc {
	// Create a rate limiter with 10 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}

	// Create a new store
	store := memory.NewStore()

	// Create a new limiter instance
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get the IP address of the client
		ip := c.ClientIP()

		// Get the context for the IP
		context, err := instance.Get(c, ip)
		if err != nil {
			logger.Error("Failed to get rate limit context",
				zap.Error(err),
				zap.String("ip", ip))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Check if the request is allowed
		if context.Reached {
			logger.Warn("Rate limit exceeded",
				zap.String("ip", ip),
				zap.Int64("limit", context.Limit),
				zap.Int64("remaining", context.Remaining),
				zap.Int64("reset", context.Reset))
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		// Add rate limit headers to the response
		c.Header("X-RateLimit-Limit", string(context.Limit))
		c.Header("X-RateLimit-Remaining", string(context.Remaining))
		c.Header("X-RateLimit-Reset", string(context.Reset))

		c.Next()
	}
}
