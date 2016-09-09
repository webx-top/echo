package session

import (
	"runtime"

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

var boltDB *bolt.DB
var onCloseBolt func() error

type BoltStore interface {
	ss.Store
}

func CloseBolt(boltDB *bolt.DB) {
	if boltDB == nil {
		return
	}
	boltDB.Close()
	if onCloseBolt != nil {
		onCloseBolt()
	}
}

type BoltOptions struct {
	File           string               `json:"file"`
	KeyPairs       [][]byte             `json:"keyPairs"`
	BucketName     string               `json:"bucketName"`
	SessionOptions *echo.SessionOptions `json:"session"`
}

//./sessions.db
func NewBoltStore(opts *BoltOptions) (BoltStore, error) {
	var err error
	if boltDB == nil {
		boltDB, err = bolt.Open(opts.File, 0666, nil)
		if err != nil {
			panic(err)
		}
		quiteC, doneC := reaper.Run(boltDB, reaper.Options{})
		onCloseBolt = func() error {
			// Invoke a reaper which checks and removes expired sessions periodically.
			reaper.Quit(quiteC, doneC)
			return nil
		}
		runtime.SetFinalizer(boltDB, CloseBolt)
	}
	config := store.Config{
		SessionOptions: sessions.Options{
			Path:     opts.SessionOptions.Path,
			Domain:   opts.SessionOptions.Domain,
			MaxAge:   opts.SessionOptions.MaxAge,
			Secure:   opts.SessionOptions.Secure,
			HttpOnly: opts.SessionOptions.HttpOnly,
		},
		DBOptions: store.Options{opts.BucketName},
	}
	stor, err := store.New(boltDB, config, opts.KeyPairs...)
	if err != nil {
		return nil, err
	}
	return &boltStore{Store: stor, config: &config, keyPairs: opts.KeyPairs}, nil
}

type boltStore struct {
	*store.Store
	config   *store.Config
	keyPairs [][]byte
}

func (c *boltStore) Options(options echo.SessionOptions) {
	c.config.SessionOptions = sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	stor, err := store.New(boltDB, *c.config, c.keyPairs...)
	if err != nil {
		panic(err.Error())
	}
	c.Store = stor
}
