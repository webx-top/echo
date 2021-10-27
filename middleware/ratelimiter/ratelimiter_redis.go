package ratelimiter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Redis limiter imp here
type redisLimiter struct {
	sha1, max, duration string
	rc                  RedisClient
}

func newRedisLimiter(options *RateLimiterConfig) *limiter {
	sha1, err := options.Client.LuaScriptLoad(LuaScriptForRedis)
	if err != nil {
		fmt.Println("redis is not working properly. use; docker run -it -p 6379:6379 --name my-redis -d redis")
		panic(err)
	}
	r := &redisLimiter{
		rc:       options.Client,
		sha1:     sha1,
		max:      strconv.FormatInt(int64(options.Max), 10),
		duration: strconv.FormatInt(int64(options.Duration/time.Millisecond), 10),
	}
	return &limiter{r, options.Prefix}
}

func (r *redisLimiter) removeLimit(key string) error {
	return r.rc.DeleteKey(key)
}

func (r *redisLimiter) getLimit(key string, policy ...int) ([]interface{}, error) {
	keys := []string{key, fmt.Sprintf("{%s}:S", key)}
	capacity := 3
	length := len(policy)
	if length > 2 {
		capacity = length + 1
	}

	//fmt.Printf("redis max limit (%s) for (%s)",r.max,key)
	args := make([]interface{}, capacity, capacity)
	args[0] = genTimestamp()
	if length == 0 {
		args[1] = r.max
		args[2] = r.duration
	} else {
		for i, val := range policy {
			if val <= 0 {
				return nil, errors.New("ratelimiter: must be positive integer")
			}
			args[i+1] = strconv.FormatInt(int64(val), 10)
		}
	}

	res, err := r.rc.EvalulateSha(r.sha1, keys, args...)
	if err != nil && isNoScriptErr(err) {
		// try to load lua for cluster client and ring client for nodes changing.
		_, err = r.rc.LuaScriptLoad(LuaScriptForRedis)
		if err == nil {
			res, err = r.rc.EvalulateSha(r.sha1, keys, args...)
		}
	}

	if err == nil {
		arr, ok := res.([]interface{})
		if ok && len(arr) == 4 {
			return arr, nil
		}
		err = errors.New("Invalid result")
	}
	return nil, err
}

func genTimestamp() string {
	time := time.Now().UnixNano() / 1e6
	return strconv.FormatInt(time, 10)
}
func isNoScriptErr(err error) bool {
	return strings.HasPrefix(err.Error(), "NOSCRIPT ")
}

//LuaScriptForRedis script loading for cluster client and ring client for nodes changing. based on links below
//https://github.com/thunks/thunk-ratelimiter
//https://github.com/thunks/thunk-ratelimiter/blob/master/ratelimiter.lua
const LuaScriptForRedis string = `
-- KEYS[1] target hash key
-- KEYS[2] target hash key
-- ARGV[n >= 3] current timestamp, max count, duration, max count, duration, ...

-- HASH: KEYS[1]
--   field:ct(count)
--   field:lt(limit)
--   field:dn(duration)
--   field:rt(reset)

local res = {}
local policyCount = (#ARGV - 1) / 2
local limit = redis.call('hmget', KEYS[2], 'ct', 'lt', 'dn', 'rt')

if limit[1] then

  res[1] = tonumber(limit[1]) - 1
  res[2] = tonumber(limit[2])
  res[3] = tonumber(limit[3]) or ARGV[3]
  res[4] = tonumber(limit[4])

  if policyCount > 1 and res[1] == -1 then
    redis.call('incr', KEYS[1])
    redis.call('pexpire', KEYS[1], res[3] * 2)
    local index = tonumber(redis.call('get', KEYS[1]))
    if index == 1 then
      redis.call('incr', KEYS[1])
    end
  end

  if res[1] >= -1 then
    redis.call('hincrby', KEYS[2], 'ct', -1)
  else
    res[1] = -1
  end

else

  local index = 1
  if policyCount > 1 then
    index = tonumber(redis.call('get', KEYS[1])) or 1
    if index > policyCount then
      index = policyCount
    end
  end

  local total = tonumber(ARGV[index * 2])
  res[1] = total - 1
  res[2] = total
  res[3] = tonumber(ARGV[index * 2 + 1])
  res[4] = tonumber(ARGV[1]) + res[3]

  redis.call('hmset', KEYS[2], 'ct', res[1], 'lt', res[2], 'dn', res[3], 'rt', res[4])
  redis.call('pexpire', KEYS[2], res[3])

end

return res
`
