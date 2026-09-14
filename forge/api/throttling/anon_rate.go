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
	store    Store
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

// NewAnonRateThrottleWithStore creates a new anonymous rate throttle with a custom store.
func NewAnonRateThrottleWithStore(rate string, store Store) *AnonRateThrottle {
	throttle := NewAnonRateThrottle(rate)
	throttle.store = store
	return throttle
}

// WithStore sets the rate limit store.
func (t *AnonRateThrottle) WithStore(store Store) *AnonRateThrottle {
	t.store = store
	return t
}

// AllowRequest checks whether the anonymous request should be allowed.
func (t *AnonRateThrottle) AllowRequest(r *http.Request, view interface{}) (bool, time.Duration, error) {
	if t.parseErr != nil {
		return true, 0, t.parseErr
	}
	key := "throttle_anon_" + t.GetScope(r, view)
	if t.store != nil {
		allowed, retryAfter := t.store.Allow(key)
		return allowed, retryAfter, nil
	}
	if t.limiter == nil {
		return true, 0, nil
	}
	allowed, retryAfter := t.limiter.Reserve(key)
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
