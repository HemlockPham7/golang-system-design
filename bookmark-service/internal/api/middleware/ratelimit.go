package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HemlockPham7/golang-system-design/internal/repository/ratelimit"
	"github.com/HemlockPham7/golang-system-design/pkg/requestutils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type RateLimit interface {
	RateLimit() gin.HandlerFunc
}

type rateLimit struct {
	repository ratelimit.Repository
}

func NewRateLimit(repository ratelimit.Repository) RateLimit {
	return &rateLimit{repository: repository}
}

const (
	rateLimitInterval  = 1 * time.Minute // sliding window ton tai trong bao lau, moi mot phut check bao nhieu request
	rateLimitCount     = 10              // so max request thuc hien trong 1 phut
	rateLimitKeyFormat = "rate_limit:%s"
)

func (r *rateLimit) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user-id from request
		uid, err := requestutils.GetUserIDFromRequest(c)
		if err != nil {
			return
		}

		// create rate limit key
		rateLimitKey := fmt.Sprintf(rateLimitKeyFormat, uid)

		// get current rate limit
		currentRate, err := r.repository.GetCurrentRateLimit(c, rateLimitKey)
		if err != nil {
			log.Error().Err(err).Msg("failed to get current rate limit")
		}

		// check if rate limit exceeded
		if currentRate >= rateLimitCount {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		// increase rate limit
		r.repository.IncreaseRateLimit(c, rateLimitKey, rateLimitInterval)
		c.Next()
	}
}
