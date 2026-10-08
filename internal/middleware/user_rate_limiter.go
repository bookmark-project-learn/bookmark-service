package middleware

import (
	"net/http"
	"time"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/request_ultils"
	"github.com/bookmark-project-learn/bookmark-common-libs/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Define a rate limiter with a limit of 10 requests per second
var limitRequestMax = 10
var limitRequestDuration = 1 * time.Second // in seconds
var rateLimitKeyPrefix = "user_rate_limit_"

type UserRateLimiter struct {
	rateLimiterRepository ratelimiter.RateLimiter
}

func NewUserRateLimiter(rateLimiterRepository ratelimiter.RateLimiter) *UserRateLimiter {
	return &UserRateLimiter{
		rateLimiterRepository: rateLimiterRepository,
	}
}

// Define a handler function that enforces the rate limit
func (u *UserRateLimiter) UserRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := request_ultils.GetSubjectFromClaims(c)
		if err != nil || userId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
			return
		}

		rateLimitKey := rateLimitKeyPrefix + userId
		count, err := u.rateLimiterRepository.GetCurrentRateLimit(c, rateLimitKey)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get current rate limit")
		}

		if count >= limitRequestMax {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Rate limit exceeded"})
			return
		}

		err = u.rateLimiterRepository.IncreaseRateLimit(c, rateLimitKey, limitRequestDuration)
		if err != nil {
			log.Error().Err(err).Msg("Failed to increase rate limit")
		}
		c.Next()
	}
}
