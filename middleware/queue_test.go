package middleware

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	myTesting "github.com/webx-top/echo/testing"
)

func TestQueueWithConfig(t *testing.T) {
	e := echo.New()

	handler := echo.HandlerFunc(func(c echo.Context) error {
		pts := c.Get("procTime")
		procTime, _ := pts.(int)

		time.Sleep(time.Duration(procTime) * time.Millisecond)

		return c.NoContent(http.StatusOK)
	})

	mw := QueueWithConfig(QueueConfig{
		QueueSize:     2,
		Workers:       1,
		QueueTimeout:  200 * time.Millisecond,
		WorkerTimeout: 100 * time.Millisecond,
	})

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

	var errQueueFull, errQueueTimeout bool

	for i := 0; i < len(testCases); i++ {
		c := <-ch

		if c == http.StatusTooManyRequests {
			errQueueFull = true
		}

		if c == http.StatusRequestTimeout {
			errQueueTimeout = true
		}
	}

	assert.Equal(t, true, errQueueFull)
	assert.Equal(t, true, errQueueTimeout)
}

func TestQueueWithConfig_panic(t *testing.T) {
	e := echo.New()

	handler := echo.HandlerFunc(func(c echo.Context) error {
		panic(`panic should release semaphore resources`)
	})

	mw := QueueWithConfig(QueueConfig{
		QueueSize:     2,
		Workers:       1,
		QueueTimeout:  200 * time.Millisecond,
		WorkerTimeout: 100 * time.Millisecond,
	})

	recoverMw := RecoverWithConfig(RecoverConfig{
		DisableStackAll:   true,
		DisablePrintStack: true,
	})

	expectedCalls := 5
	actualCallsChan := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			req, rec := myTesting.NewRequestAndResponse(http.MethodGet, "/")

			c := e.NewContext(req, rec)

			err := recoverMw(mw(handler)).Handle(c)
			if err != nil {
				return
			}

			actualCallsChan <- struct{}{}

		}()

	}

	wg.Wait()

	actualCalls := len(actualCallsChan)

	assert.Equal(t, expectedCalls, actualCalls)
}

func TestQueueWithConfig_skipperNoSkip(t *testing.T) {
	e := echo.New()

	var beforeFuncRan bool
	handler := echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req, rec := myTesting.NewRequestAndResponse(http.MethodGet, "/")

	c := e.NewContext(req, rec)

	mw := QueueWithConfig(QueueConfig{
		Skipper: func(c echo.Context) bool {
			beforeFuncRan = true
			return false
		},
	})

	_ = mw(handler).Handle(c)

	assert.Equal(t, true, beforeFuncRan)
}
