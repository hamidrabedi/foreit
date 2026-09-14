package throttling

import (
	"net/http"
	"time"

	internalratelimit "github.com/forgego/forge/internal/ratelimit"
	"github.com/forgego/forge/netutil"
)

// AnonRateThrottle throttles anonymous requests by client IP.
type AnonRateThrottle struct {
	Rate     string
	Scope    string
	limiter  *internalratelimit.KeyedLimiter
	parseErr error
}

// NewAnonRateThrottle creates a new anonymous rate throttle.
func NewAnonRateThrottle(rate string) *AnonRateThrottle {
	limit, window, err := parseRate(rate)
	return &AnonRateThrottle{
		Rate: rate, Scope: "anon", parseErr: err,
		limiter: newKeyedLimiter(limit, window, err),
	}
}

// AllowRequest checks whether the anonymous request should be allowed.
func (t *AnonRateThrottle) AllowRequest(r *http.Request, view interface{}) (bool, time.Duration, error) {
	if t.parseErr != nil {
		return true, 0, t.parseErr
	}
	allowed, retryAfter := t.limiter.Reserve("throttle_anon_" + t.GetScope(r, view))
	return allowed, retryAfter, nil
}

// GetScope returns the client IP address for anonymous users.
func (t *AnonRateThrottle) GetScope(r *http.Request, view interface{}) string {
	return getClientIP(r)
}

func newKeyedLimiter(limit int, window time.Duration, err error) *internalratelimit.KeyedLimiter {
	if err != nil {
		return nil
	}
	return internalratelimit.NewKeyedLimiter(limit, window)
}

func getClientIP(r *http.Request) string {
	return netutil.ClientIP(r, netutil.TrustedProxies())
}
