package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func testConnect() {

	redisOptions := &RedisOptions{
		Size:     10,
		Network:  `tcp`,
		Address:  `127.0.0.1:6379`,
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
	testConnect()

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
	testConnect()

	store := ss.Get(`redis`)
	ctx := defaults.NewMockContext()
	sess, err := store.New(ctx, `TEST_REDIS_SID`)
	assert.NoError(t, err)
	sess.AddFlash(`TEST2`)
	err = store.Save(ctx, sess)
	assert.NoError(t, err)
}
