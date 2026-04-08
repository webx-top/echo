package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	goredislib "github.com/redis/go-redis/v9"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/pprof"
	"github.com/webx-top/echo/middleware/ratelimiter"
	"github.com/webx-top/echo/testing"
)

// RedisClient Implements RedisClient for redis.Client
type RedisClient struct {
	*goredislib.Client
}

func (c *RedisClient) DeleteKey(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *RedisClient) EvalulateSha(ctx context.Context, sha1 string, keys []string, args ...any) (any, error) {
	return c.Client.EvalSha(ctx, sha1, keys, args...).Result()
}

func (c *RedisClient) LuaScriptLoad(ctx context.Context, script string) (string, error) {
	return c.Client.ScriptLoad(ctx, script).Result()
}

func main() {
	e := echo.New()
	e.Get("/", echo.HandlerFunc(func(c echo.Context) error {
		return c.String("Hello, World!")
	}), ratelimiter.RateLimiterWithConfig(ratelimiter.RateLimiterConfig{
		Max:      5,
		Duration: time.Second,
	}))
	redisDisabled := false
	if !redisDisabled {

		var client = goredislib.NewClient(&goredislib.Options{
			Addr: `127.0.0.1:6379`,
		})
		defer client.Close()
		e.Get("/redis", echo.HandlerFunc(func(c echo.Context) error {
			return c.String("Hello, World!")
		}), ratelimiter.RateLimiterWithConfig(ratelimiter.RateLimiterConfig{
			Max:      5,
			Duration: time.Second,
			Client:   &RedisClient{Client: client},
		}))

	}

	pprof.Wrap(e)

	go e.Run(standard.New(":4444"))

	time.Sleep(time.Second * 2)
	var codes []int
	wg := sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rec := testing.Request(`GET`, `/redis`, e)
			after := rec.Header().Get(`X-Retry-After`)
			if len(after) > 0 {
				fmt.Println(`X-Retry-After:`, after)
				m, err := time.ParseDuration(after)
				if err != nil {
					fmt.Println(err)
				} else {
					//time.Sleep(m + time.Millisecond)
					time.Sleep(m + time.Second)
				}
			}
			codes = append(codes, rec.Code)
		}()
	}
	wg.Wait()
	echo.Dump(codes)

	<-make(chan struct{})
}
