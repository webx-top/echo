package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func testConnect(addr string) {

	redisOptions := &RedisOptions{
		Size:     10,
		Network:  `tcp`,
		Address:  addr,
		Password: ``,
		DB:       0,
		KeyPairs: [][]byte{
			[]byte(`12345678901234567890123456789012`),
			[]byte(`12345678901234567890123456789012`),
		},
		MaxAge:       86400,
		MaxReconnect: 30,
	}
	RegWithOptions(redisOptions)
}

func TestReconnect(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	testConnect(s.Addr())

	store := ss.Get(`redis`)
	ctx := defaults.NewMockContext()
	sess, err := store.New(ctx, `TEST_REDIS_SID`)
	assert.NoError(t, err)
	sess.AddFlash(`TEST1`)
	err = store.Save(ctx, sess)
	assert.NoError(t, err)

	sess2, err := store.Get(ctx, `TEST_REDIS_SID`)
	assert.NoError(t, err)
	t.Logf(`%+v`, sess.Values)
	t.Logf(`%+v`, sess2.Values)
}

func TestReconnect2(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	testConnect(s.Addr())

	store := ss.Get(`redis`)
	ctx := defaults.NewMockContext()
	sess, err := store.New(ctx, `TEST_REDIS_SID`)
	assert.NoError(t, err)
	sess.AddFlash(`TEST2`)
	err = store.Save(ctx, sess)
	assert.NoError(t, err)
}
