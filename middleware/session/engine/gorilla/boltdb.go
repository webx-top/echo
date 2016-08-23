package session

import (
	"runtime"

	"github.com/admpub/boltstore/reaper"
	"github.com/admpub/boltstore/store"
	"github.com/admpub/sessions"
	"github.com/boltdb/bolt"
	"github.com/webx-top/echo"
)

var boltDB *bolt.DB
var onCloseBolt func() error

type BoltStore interface {
	Store
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

//./sessions.db
func NewBoltStore(dbFile string, options echo.SessionOptions, bucketName []byte, keyPairs ...[]byte) (BoltStore, error) {
	var err error
	if boltDB == nil {
		boltDB, err = bolt.Open(dbFile, 0666, nil)
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
			Path:     options.Path,
			Domain:   options.Domain,
			MaxAge:   options.MaxAge,
			Secure:   options.Secure,
			HttpOnly: options.HttpOnly,
		},
		DBOptions: store.Options{bucketName},
	}
	stor, err := store.New(boltDB, config, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &boltStore{Store: stor, config: &config, keyPairs: keyPairs}, nil
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
