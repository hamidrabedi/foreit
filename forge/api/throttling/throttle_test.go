package throttling

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/forgego/forge/api/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAuthUser struct{ ID string }

func TestParseRate(t *testing.T) {
	limit, window, err := parseRate("100/hour")
	require.NoError(t, err)
	assert.Equal(t, 100, limit)
	assert.Equal(t, time.Hour, window)
}

func TestAnonRateThrottle_AllowsBurstThenReturnsRetry(t *testing.T) {
	throttle := NewAnonRateThrottle("3/minute")
	defer throttle.limiter.Close()
	request := httptest.NewRequest("GET", "/test", nil)
	request.RemoteAddr = "192.0.2.1:1234"
	for i := 0; i < 3; i++ {
		allowed, retryAfter, err := throttle.AllowRequest(request, nil)
		require.NoError(t, err)
		assert.True(t, allowed)
		assert.Zero(t, retryAfter)
	}
	allowed, retryAfter, err := throttle.AllowRequest(request, nil)
	require.NoError(t, err)
	assert.False(t, allowed)
	assert.Positive(t, retryAfter)
}

func TestUserAndAnonymousThrottleKeysAreSeparate(t *testing.T) {
	anon := NewAnonRateThrottle("1/minute")
	user := NewUserRateThrottle("1/minute")
	defer anon.limiter.Close()
	defer user.limiter.Close()
	request := httptest.NewRequest("GET", "/test", nil)
	request.RemoteAddr = "192.0.2.1:1234"
	allowed, _, err := anon.AllowRequest(request, nil)
	require.NoError(t, err)
	assert.True(t, allowed)
	authentication.SetUserOnRequest(request, &mockAuthUser{ID: "42"})
	allowed, _, err = user.AllowRequest(request, nil)
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestUserRateThrottleFallsBackToClientIP(t *testing.T) {
	throttle := NewUserRateThrottle("1/minute")
	defer throttle.limiter.Close()
	request := httptest.NewRequest("GET", "/test", nil)
	allowed, _, err := throttle.AllowRequest(request, nil)
	require.NoError(t, err)
	assert.True(t, allowed)
	allowed, retryAfter, err := throttle.AllowRequest(request, nil)
	require.NoError(t, err)
	assert.False(t, allowed)
	assert.Positive(t, retryAfter)
}

func TestCheckThrottlesReturnsWaitDuration(t *testing.T) {
	throttle := NewAnonRateThrottle("1/minute")
	defer throttle.limiter.Close()
	request := httptest.NewRequest("GET", "/test", nil)
	require.NoError(t, CheckThrottles(request, nil, []Throttle{throttle}))
	err := CheckThrottles(request, nil, []Throttle{throttle})
	throttled, ok := err.(*ThrottledError)
	require.True(t, ok)
	assert.Positive(t, throttled.WaitDuration)
}
