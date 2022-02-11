package ratelimiter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"gopkg.in/redis.v5"

	"github.com/webx-top/echo"
	te "github.com/webx-top/echo/testing"
)

// Implements RedisClient for redis.Client
type redisClient struct {
	*redis.Client
}

func (c *redisClient) DeleteKey(key string) error {
	return c.Del(key).Err()
}

func (c *redisClient) EvalulateSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return c.EvalSha(sha1, keys, args...).Result()
}

func (c *redisClient) LuaScriptLoad(script string) (string, error) {
	return c.ScriptLoad(script).Result()
}

// Implements RedisClient for redis.Client
type failedClient struct {
	*redis.Client
}

func (c *failedClient) DeleteKey(key string) error {
	return c.Del(key).Err()
}

func (c *failedClient) EvalulateSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return nil, errors.New("noscript mock error")
}

func (c *failedClient) LuaScriptLoad(script string) (string, error) {
	return c.ScriptLoad(script).Result()
}

func TestRedisRatelimiter(t *testing.T) {

	e := echo.New()
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	var client = redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer client.Close()

	t.Run("Redis ratelimiter middleware", func(t *testing.T) {

		t.Run("FakeRedis is running as excepted", func(t *testing.T) {
			pong, err := client.Ping().Result()
			assert.Nil(t, err)
			assert.Equal(t, "PONG", pong)
		})

		t.Run("New instrance running with redis option as excepted", func(t *testing.T) {

			rateLimitWithConfig := RateLimiterWithConfig(RateLimiterConfig{
				Client: &redisClient{client},
			})
			assert.NotNil(t, rateLimitWithConfig)
		})

		t.Run("Get method should return ok with excepted remaining and limit values", func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(te.WrapRequest(req), te.WrapResponse(req, rec))

			rateLimitWithConfig := RateLimiterWithConfig(RateLimiterConfig{
				Max:      100,
				Client:   &redisClient{client},
				Duration: time.Minute * 1,
			})

			hx := rateLimitWithConfig(echo.HandlerFunc(func(c echo.Context) error {
				return c.String("test")
			}))
			hx.Handle(c)
			assert.Contains(t, rec.Header().Get("X-Ratelimit-Remaining"), "99")
			assert.Contains(t, rec.Header().Get("X-Ratelimit-Limit"), "100")

		})

		t.Run("Get method should throw too many request", func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/alper", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(te.WrapRequest(req), te.WrapResponse(req, rec))

			xx := RateLimiterWithConfig(RateLimiterConfig{
				Max:      2,
				Client:   &redisClient{client},
				Duration: time.Minute * 1,
			})

			hx := xx.Handle(echo.HandlerFunc(func(c echo.Context) error {
				return c.String("test")
			}))
			hx.Handle(c)
			hx.Handle(c)
			expectedErrorStatus, ok := hx.Handle(c).(*echo.HTTPError)

			if ok {
				assert.Contains(t, rec.Header().Get("X-Ratelimit-Remaining"), "-1")
				assert.Equal(t, http.StatusTooManyRequests, expectedErrorStatus.Code)
			} else {
				assert.Error(t, errors.New("it should throw too many ruqest exception"))
			}
		})

	})

	t.Run("Redis ratelimiter implementation", func(t *testing.T) {
		var id = genID()
		var duration = time.Duration(60 * 1e9)
		var redisLimiter *limiter
		redisLimiter = newRedisLimiter(&RateLimiterConfig{

			Client:   &redisClient{client},
			Max:      100,
			Duration: duration,
		})

		t.Run("New instance running with failedClient should be", func(t *testing.T) {

			var limiter *limiter
			limiter = newRedisLimiter(&RateLimiterConfig{

				Client: &failedClient{client},
			})
			policy := []int{2, 100}
			res, err := limiter.Get(id, policy...)

			assert.Equal(t, "noscript mock error", err.Error())
			assert.Equal(t, 0, res.Total)
			assert.Equal(t, 0, res.Remaining)
			assert.Equal(t, time.Duration(0), res.Duration)
		})

		t.Run("Redislimiter.Get method should be", func(t *testing.T) {

			res, err := redisLimiter.Get(id)
			assert.Nil(t, err)
			assert.Equal(t, res.Total, 100)
			assert.Equal(t, res.Remaining, 99)
			assert.Equal(t, res.Duration, duration)
			assert.True(t, res.Reset.UnixNano() > time.Now().UnixNano())

			res, err = redisLimiter.Get(id)
			assert.Nil(t, err)
			assert.Equal(t, res.Total, 100)
			assert.Equal(t, res.Remaining, 98)

		})

		t.Run("Redislimiter.Get with invalid args should throw error", func(t *testing.T) {
			_, err := redisLimiter.Get(id, 10)
			assert.Equal(t, "ratelimiter: must be paired values", err.Error())

			_, err2 := redisLimiter.Get(id, -1, 10)
			assert.Equal(t, "ratelimiter: must be positive integer", err2.Error())

			_, err3 := redisLimiter.Get(id, 10, 0)
			assert.Equal(t, "ratelimiter: must be positive integer", err3.Error())
		})

		t.Run("Redislimiter.Get with policy", func(t *testing.T) {

			idx := genID()
			assert := assert.New(t)

			policy := []int{2, 100}

			res, err := redisLimiter.Get(idx, policy...)
			assert.Nil(err)
			assert.Equal(2, res.Total)
			assert.Equal(1, res.Remaining)
			assert.Equal(time.Millisecond*100, res.Duration)

			res, err = redisLimiter.Get(idx, policy...)
			assert.Nil(err)
			assert.Equal(0, res.Remaining)

			res, err = redisLimiter.Get(idx, policy...)
			assert.Nil(err)
			assert.Equal(-1, res.Remaining)
		})

		t.Run("Redislimiter.Remove method should be", func(t *testing.T) {

			err := redisLimiter.Remove(id)
			assert.Nil(t, err)

			err = redisLimiter.Remove(id)
			assert.Nil(t, err)

			res, err := redisLimiter.Get(id)
			assert.Nil(t, err)
			assert.Equal(t, res.Total, 100)
			assert.Equal(t, res.Remaining, 97)

		})

	})
}
