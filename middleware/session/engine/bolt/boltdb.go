package bolt

import (
	"runtime"
	"time"

	"github.com/admpub/boltstore/reaper"
	"github.com/admpub/boltstore/store"
	"github.com/admpub/sessions"
	"github.com/boltdb/bolt"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func New(opts *BoltOptions) BoltStore {
	store, err := NewBoltStore(opts)
	if err != nil {
		panic(err.Error())
	}
	return store
}

func Reg(store BoltStore, args ...string) {
	name := `bolt`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *BoltOptions, args ...string) {
	Reg(New(opts), args...)
}

type BoltStore interface {
	ss.Store
}

type BoltOptions struct {
	File           string               `json:"file"`
	KeyPairs       [][]byte             `json:"keyPairs"`
	BucketName     string               `json:"bucketName"`
	SessionOptions *echo.SessionOptions `json:"session"`
}

// NewBoltStore ./sessions.db
func NewBoltStore(opts *BoltOptions) (BoltStore, error) {
	db, err := bolt.Open(opts.File, 0666, nil)
	if err != nil {
		return nil, err
	}
	config := store.Config{
		SessionOptions: sessions.Options{
			Path:     opts.SessionOptions.Path,
			Domain:   opts.SessionOptions.Domain,
			MaxAge:   opts.SessionOptions.MaxAge,
			Secure:   opts.SessionOptions.Secure,
			HttpOnly: opts.SessionOptions.HttpOnly,
		},
		DBOptions: store.Options{BucketName: []byte(opts.BucketName)},
	}
	stor, err := store.New(db, config, opts.KeyPairs...)
	if err != nil {
		return nil, err
	}
	b := &boltStore{Store: stor, db: db, config: &config, keyPairs: opts.KeyPairs}
	b.quiteC, b.doneC = reaper.Run(db, reaper.Options{
		BucketName:    []byte(opts.BucketName),
		CheckInterval: time.Duration(int64(opts.SessionOptions.MaxAge)) * time.Second,
	})
	runtime.SetFinalizer(b, func(b *boltStore) {
		b.Close()
	})
	return b, nil
}

type boltStore struct {
	*store.Store
	db       *bolt.DB
	config   *store.Config
	keyPairs [][]byte
	quiteC   chan<- struct{}
	doneC    <-chan struct{}
}

func (c *boltStore) Options(options echo.SessionOptions) {
	c.config.SessionOptions = sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	stor, err := store.New(c.db, *c.config, c.keyPairs...)
	if err != nil {
		panic(err.Error())
	}
	c.Store = stor
}

func (c *boltStore) Close() {
	// Invoke a reaper which checks and removes expired sessions periodically.
	reaper.Quit(c.quiteC, c.doneC)
	c.db.Close()
}
