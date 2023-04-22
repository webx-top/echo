/*
Copyright 2016 Wenhui Shen <www.webx.top>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package ratelimit

import (
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/middleware/ratelimit/config"
)

func LimitMiddleware(limiter *config.Limiter, bucketExpireTTL ...time.Duration) echo.MiddlewareFunc {
	var limit func(engine.Request) *echo.HTTPError
	if len(bucketExpireTTL) > 0 {
		var cleanUpInterval time.Duration
		if len(bucketExpireTTL) > 1 {
			cleanUpInterval = bucketExpireTTL[1]
		} else {
			cleanUpInterval = time.Minute * 5
		}
		limiter.SetExpiring(bucketExpireTTL[0], cleanUpInterval)
		limit = func(req engine.Request) *echo.HTTPError {
			return LimitByRequestWithExpiring(limiter, req, bucketExpireTTL[0])
		}
	} else {
		limit = func(req engine.Request) *echo.HTTPError {
			return LimitByRequest(limiter, req)
		}
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			SetResponseHeaders(limiter, c.Response())
			httpError := limit(c.Request())
			if httpError != nil {
				return httpError
			}
			return h.Handle(c)
		})
	}
}

func LimitHandler(limiter *config.Limiter, bucketExpireTTL ...time.Duration) echo.MiddlewareFunc {
	return LimitMiddleware(limiter, bucketExpireTTL...)
}

// New is a convenience function to config.NewLimiter.
func New(max int64, ttl time.Duration) *config.Limiter {
	return config.NewLimiter(max, ttl)
}

// NewExpiringBuckets 带过期时间
// @param max 每个周期的最大连接数
// @param ttl 周期
// @param bucketDefaultExpirationTTL 默认有效时长
// @param bucketExpireJobInterval 默认清理周期
func NewExpiringBuckets(max int64, ttl, bucketDefaultExpirationTTL, bucketExpireJobInterval time.Duration) *config.Limiter {
	return config.NewLimiterExpiringBuckets(max, ttl, bucketDefaultExpirationTTL, bucketExpireJobInterval)
}

// LimitByKeys keeps track number of request made by keys separated by pipe.
// It returns HTTPError when limit is exceeded.
func LimitByKeys(limiter *config.Limiter, keys []string) *echo.HTTPError {
	if limiter.LimitReached(strings.Join(keys, "|")) {
		return echo.NewHTTPError(limiter.StatusCode, limiter.Message)
	}

	return nil
}

// LimitByKeysWithCustomTokenBucketTTL keeps track number of request made by keys separated by pipe.
// It returns HTTPError when limit is exceeded.
// User can define a TTL for the key to expire
func LimitByKeysWithCustomTokenBucketTTL(limiter *config.Limiter, keys []string, bucketExpireTTL time.Duration) *echo.HTTPError {
	if limiter.LimitReachedWithCustomTokenBucketTTL(strings.Join(keys, "|"), bucketExpireTTL) {
		return echo.NewHTTPError(limiter.StatusCode, limiter.Message)
	}

	return nil
}

const HTTPTimeLayout = "Mon, 02 Jan 2006 15:04:05 GMT"

// SetResponseHeaders configures X-Rate-Limit-Limit and X-Rate-Limit-Duration
func SetResponseHeaders(limiter *config.Limiter, w engine.Response) {
	w.Header().Add("X-Rate-Limit-Limit", strconv.FormatInt(limiter.Max, 10))
	w.Header().Add("X-Rate-Limit-Duration", limiter.TTL.String())
	w.Header().Add("Retry-After", time.Now().Add(limiter.TTL).Format(HTTPTimeLayout))
}

// LimitByRequest builds keys based on http.Request struct,
// loops through all the keys, and check if any one of them returns HTTPError.
func LimitByRequest(limiter *config.Limiter, r engine.Request) *echo.HTTPError {
	sliceKeys := BuildKeys(limiter, r)

	// Loop sliceKeys and check if one of them has error.
	for _, keys := range sliceKeys {
		httpError := LimitByKeys(limiter, keys)
		if httpError != nil {
			return httpError
		}
	}

	return nil
}

func LimitByRequestWithExpiring(limiter *config.Limiter, r engine.Request, bucketExpireTTL time.Duration) *echo.HTTPError {
	sliceKeys := BuildKeys(limiter, r)

	// Loop sliceKeys and check if one of them has error.
	for _, keys := range sliceKeys {
		httpError := LimitByKeysWithCustomTokenBucketTTL(limiter, keys, bucketExpireTTL)
		if httpError != nil {
			return httpError
		}
	}

	return nil
}

// StringInSlice finds needle in a slice of strings.
func StringInSlice(sliceString []string, needle string) bool {
	for _, b := range sliceString {
		if b == needle {
			return true
		}
	}
	return false
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// RemoteIP finds IP Address given http.Request struct.
func RemoteIP(ipLookups []string, r engine.Request) string {

	for _, lookup := range ipLookups {
		if lookup == "RemoteAddr" {
			return ipAddrFromRemoteAddr(r.RemoteAddress())
		}
		if lookup == "X-Forwarded-For" {
			if forwardedFor := r.Header().Get("X-Forwarded-For"); len(forwardedFor) > 0 {
				// X-Forwarded-For is potentially a list of addresses separated with ","
				return strings.TrimSpace(strings.SplitN(forwardedFor, ",", 2)[0])
			}
		}
		if lookup == "X-Real-IP" {
			if realIP := r.Header().Get("X-Real-IP"); len(realIP) > 0 {
				return realIP
			}
		}
	}

	return ""
}

// BuildKeys generates a slice of keys to rate-limit by given config and request structs.
func BuildKeys(limiter *config.Limiter, r engine.Request) [][]string {
	remoteIP := RemoteIP(limiter.IPLookups, r)
	path := r.URL().Path()
	sliceKeys := make([][]string, 0)

	// Don't BuildKeys if remoteIP is blank.
	if remoteIP == "" {
		return sliceKeys
	}

	if limiter.Methods != nil && limiter.Headers != nil && limiter.BasicAuthUsers != nil {
		// Limit by HTTP methods and HTTP headers+values and Basic Auth credentials.
		if StringInSlice(limiter.Methods, r.Method()) {
			for headerKey, headerValues := range limiter.Headers {
				if (headerValues == nil || len(headerValues) <= 0) && r.Header().Get(headerKey) != "" {
					// If header values are empty, rate-limit all request with headerKey.
					username, _, ok := r.BasicAuth()
					if ok && StringInSlice(limiter.BasicAuthUsers, username) {
						sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method(), headerKey, username})
					}

				} else if len(headerValues) > 0 && r.Header().Get(headerKey) != "" {
					// If header values are not empty, rate-limit all request with headerKey and headerValues.
					for _, headerValue := range headerValues {
						username, _, ok := r.BasicAuth()
						if ok && StringInSlice(limiter.BasicAuthUsers, username) {
							sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method(), headerKey, headerValue, username})
						}
					}
				}
			}
		}

	} else if limiter.Methods != nil && limiter.Headers != nil {
		// Limit by HTTP methods and HTTP headers+values.
		if StringInSlice(limiter.Methods, r.Method()) {
			for headerKey, headerValues := range limiter.Headers {
				if (headerValues == nil || len(headerValues) <= 0) && r.Header().Get(headerKey) != "" {
					// If header values are empty, rate-limit all request with headerKey.
					sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method(), headerKey})

				} else if len(headerValues) > 0 && r.Header().Get(headerKey) != "" {
					// If header values are not empty, rate-limit all request with headerKey and headerValues.
					for _, headerValue := range headerValues {
						sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method(), headerKey, headerValue})
					}
				}
			}
		}

	} else if limiter.Methods != nil && limiter.BasicAuthUsers != nil {
		// Limit by HTTP methods and Basic Auth credentials.
		if StringInSlice(limiter.Methods, r.Method()) {
			username, _, ok := r.BasicAuth()
			if ok && StringInSlice(limiter.BasicAuthUsers, username) {
				sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method(), username})
			}
		}

	} else if limiter.Methods != nil {
		// Limit by HTTP methods.
		if StringInSlice(limiter.Methods, r.Method()) {
			sliceKeys = append(sliceKeys, []string{remoteIP, path, r.Method()})
		}

	} else if limiter.Headers != nil {
		// Limit by HTTP headers+values.
		for headerKey, headerValues := range limiter.Headers {
			if (headerValues == nil || len(headerValues) <= 0) && r.Header().Get(headerKey) != "" {
				// If header values are empty, rate-limit all request with headerKey.
				sliceKeys = append(sliceKeys, []string{remoteIP, path, headerKey})

			} else if len(headerValues) > 0 && r.Header().Get(headerKey) != "" {
				// If header values are not empty, rate-limit all request with headerKey and headerValues.
				for _, headerValue := range headerValues {
					sliceKeys = append(sliceKeys, []string{remoteIP, path, headerKey, headerValue})
				}
			}
		}

	} else if limiter.BasicAuthUsers != nil {
		// Limit by Basic Auth credentials.
		username, _, ok := r.BasicAuth()
		if ok && StringInSlice(limiter.BasicAuthUsers, username) {
			sliceKeys = append(sliceKeys, []string{remoteIP, path, username})
		}
	} else {
		// Default: Limit by remoteIP and path.
		sliceKeys = append(sliceKeys, []string{remoteIP, path})
	}

	return sliceKeys
}
