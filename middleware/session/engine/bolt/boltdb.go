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

var (
	DefaultCheckInterval = time.Minute * 30
)

func New(opts *BoltOptions) sessions.Store {
	store, err := NewBoltStore(opts)
	if err != nil {
		panic(err.Error())
	}
	return store
}

func Reg(store sessions.Store, args ...string) {
	name := `bolt`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *BoltOptions, args ...string) sessions.Store {
	store := New(opts)
	Reg(store, args...)
	return store
}

type BoltOptions struct {
	File          string        `json:"file"`
	KeyPairs      [][]byte      `json:"keyPairs"`
	BucketName    string        `json:"bucketName"`
	CheckInterval time.Duration `json:"checkInterval"`
}

// NewBoltStore ./sessions.db
func NewBoltStore(opts *BoltOptions) (sessions.Store, error) {
	config := store.Config{
		DBOptions: store.Options{
			BucketName: []byte(opts.BucketName),
		},
	}
	b := &boltStore{
		config:   &config,
		keyPairs: opts.KeyPairs,
		dbFile:   opts.File,
		Storex: &Storex{
			Store: &store.Store{},
		},
		checkInterval: opts.CheckInterval,
	}
	b.Storex.b = b
	return b, nil
}

type Storex struct {
	*store.Store
	db          *bolt.DB
	b           *boltStore
	initialized bool
}

func (s *Storex) Get(ctx echo.Context, name string) (*sessions.Session, error) {
	if s.initialized == false {
		err := s.b.Init()
		if err != nil {
			return nil, err
		}
	}
	return s.Store.Get(ctx, name)
}

func (s *Storex) New(ctx echo.Context, name string) (*sessions.Session, error) {
	if s.initialized == false {
		err := s.b.Init()
		if err != nil {
			return nil, err
		}
	}
	return s.Store.New(ctx, name)
}

func (s *Storex) Reload(ctx echo.Context, session *sessions.Session) error {
	return s.Store.Reload(ctx, session)
}

func (s *Storex) Save(ctx echo.Context, session *sessions.Session) error {
	if s.initialized == false {
		err := s.b.Init()
		if err != nil {
			return err
		}
	}
	return s.Store.Save(ctx, session)
}

type boltStore struct {
	*Storex
	config        *store.Config
	keyPairs      [][]byte
	quiteC        chan<- struct{}
	doneC         <-chan struct{}
	dbFile        string
	checkInterval time.Duration
}

func (c *boltStore) Close() error {
	// Invoke a reaper which checks and removes expired sessions periodically.
	if c.quiteC != nil && c.doneC != nil {
		reaper.Quit(c.quiteC, c.doneC)
	}

	if c.Storex.db != nil {
		c.Storex.db.Close()
	}

	return nil
}

func (b *boltStore) Init() error {
	if b.Storex.db == nil {
		var err error
		b.Storex.db, err = bolt.Open(b.dbFile, 0666, nil)
		if err != nil {
			return err
		}
		b.Storex.Store, err = store.New(b.Storex.db, *b.config, b.keyPairs...)
		if err != nil {
			return err
		}
		if b.checkInterval == 0 {
			b.checkInterval = DefaultCheckInterval
		}
		b.quiteC, b.doneC = reaper.Run(b.Storex.db, reaper.Options{
			BucketName:    b.config.DBOptions.BucketName,
			CheckInterval: b.checkInterval,
		})
		runtime.SetFinalizer(b, func(b *boltStore) {
			b.Close()
		})
	}
	b.Storex.initialized = true
	return nil
}
