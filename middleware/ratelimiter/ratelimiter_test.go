package ratelimiter

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/webx-top/echo"
	te "github.com/webx-top/echo/testing"
)

func genID() string {
	buf := make([]byte, 12)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func TestRateLimiter(t *testing.T) {

	t.Run("ratelimiter middleware", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		t.Run("should return ok with x-remaining and x-limit value", func(t *testing.T) {

			rec := httptest.NewRecorder()
			c := e.NewContext(te.WrapRequest(req), te.WrapResponse(req, rec))
			rateLimit := RateLimiter()

			h := rateLimit(echo.HandlerFunc(func(c echo.Context) error {
				return c.String("test")
			}))
			h.Handle(c)
			assert.Contains(t, rec.Header().Get("X-Ratelimit-Remaining"), "99")
			assert.Contains(t, rec.Header().Get("X-Ratelimit-Limit"), "100")

		})

		t.Run("should throw too many request", func(t *testing.T) {

			//ratelimit with config
			rateLimitWithConfig := RateLimiterWithConfig(RateLimiterConfig{
				Max: 2,
			})

			rec := httptest.NewRecorder()
			c := e.NewContext(te.WrapRequest(req), te.WrapResponse(req, rec))
			hx := rateLimitWithConfig(echo.HandlerFunc(func(c echo.Context) error {
				return c.String("test")
			}))
			hx.Handle(c)
			hx.Handle(c)
			expectedErrorStatus := hx.Handle(c).(*echo.HTTPError)

			assert.Contains(t, rec.Header().Get("X-Ratelimit-Remaining"), "-1")
			assert.Equal(t, http.StatusTooManyRequests, expectedErrorStatus.Code)

		})

		t.Run("should return status ok after to many request status expired", func(t *testing.T) {

			expectedDuration := time.Millisecond * 5

			//ratelimit with config; expected result getting 429 after 5 second it should return 200
			rateLimitWithConfig := RateLimiterWithConfig(RateLimiterConfig{
				Max:      2,
				Duration: expectedDuration,
			})

			rec := httptest.NewRecorder()
			c := e.NewContext(te.WrapRequest(req), te.WrapResponse(req, rec))
			hx := rateLimitWithConfig(echo.HandlerFunc(func(c echo.Context) error {
				return c.String("test")
			}))
			hx.Handle(c)
			hx.Handle(c)
			expectedErrorStatus := hx.Handle(c).(*echo.HTTPError)
			assert.Equal(t, http.StatusTooManyRequests, expectedErrorStatus.Code)
			time.Sleep(expectedDuration)
			exceptedHTTPStatusOk, ok := hx.Handle(c).(*echo.HTTPError)

			if ok {
				assert.Equal(t, http.StatusOK, exceptedHTTPStatusOk.Code)
			}

		})

		t.Run("should return ok even limiter throw an exception", func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/t", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(te.WrapRequest(req), te.WrapResponse(req, rec))
			rateLimitWithConfig := RateLimiterWithConfig(RateLimiterConfig{
				SkipRateLimiterInternalError: true,
			})

			h := rateLimitWithConfig(echo.HandlerFunc(func(c echo.Context) error {
				return c.String("test")
			}))
			h.Handle(c)
			fmt.Println(rec.Result().Status)
			assert.Contains(t, rec.Result().Status, "200 OK")

		})
	})
}
