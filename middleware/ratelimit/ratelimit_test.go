package ratelimit

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	myTesting "github.com/webx-top/echo/testing"
)

func TestRateLimit(t *testing.T) {
	e := echo.New()
	var touched []int

	handler := echo.HandlerFunc(func(c echo.Context) error {
		pts := c.Get("procTime")
		fmt.Println(`--------------------`, pts, c.Request().RemoteAddress())
		procTime, _ := pts.(int)
		touched = append(touched, procTime)
		//time.Sleep(time.Duration(procTime) * time.Millisecond)

		return c.NoContent(http.StatusOK)
	})

	// Create a limiter struct.
	limiter := New(1, time.Second)
	mw := LimitHandler(limiter)

	testCases := []struct {
		procTime int // in Milliseconds
	}{
		{50},
		{99},
		{99},
		{99},
		{99},
		{150},
		{150},
		{150},
		{300},
		{300},
		{300},
	}

	ch := make(chan int, len(testCases))
	var wg sync.WaitGroup

	for _, tc := range testCases {

		wg.Add(1)

		go func(pt int) {

			defer wg.Done()
			req, rec := myTesting.NewRequestAndResponse(http.MethodGet, "/")
			req.StdRequest().RemoteAddr = `127.0.0.1`
			c := e.NewContext(req, rec)
			c.Set("procTime", pt)

			err := mw(handler).Handle(c)
			if err != nil {
				ch <- err.(*echo.HTTPError).Code
				return
			}

			ch <- rec.Status()

		}(tc.procTime)

	}

	wg.Wait()

	var errTooManyRequests, errOK bool

	for i := 0; i < len(testCases); i++ {
		c := <-ch
		if c == http.StatusTooManyRequests {
			errTooManyRequests = true
		}
		if c == http.StatusOK {
			errOK = true
		}
	}

	assert.Equal(t, 1, len(touched))
	assert.Equal(t, true, errTooManyRequests)
	assert.Equal(t, true, errOK)
}
