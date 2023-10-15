package bolt

import (
	"runtime"
	"sync"
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
	KeyPairs      [][]byte      `json:"-"`
	BucketName    string        `json:"bucketName"`
	MaxLength     int           `json:"maxLength"`
	CheckInterval time.Duration `json:"checkInterval"`
}

// NewBoltStore ./sessions.db
func NewBoltStore(opts *BoltOptions) (sessions.Store, error) {
	config := store.Config{
		DBOptions: store.Options{
			BucketName: []byte(opts.BucketName),
		},
		MaxLength: opts.MaxLength,
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
	db *bolt.DB
	b  *boltStore
}

func (s *Storex) Get(ctx echo.Context, name string) (*sessions.Session, error) {
	err := s.b.Init()
	if err != nil {
		return nil, err
	}
	return s.Store.Get(ctx, name)
}

func (s *Storex) New(ctx echo.Context, name string) (*sessions.Session, error) {
	return s.Store.New(ctx, name)
}

func (s *Storex) Reload(ctx echo.Context, session *sessions.Session) error {
	return s.Store.Reload(ctx, session)
}

func (s *Storex) Save(ctx echo.Context, session *sessions.Session) error {
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
	once          sync.Once
}

func (c *boltStore) Close() (err error) {
	// Invoke a reaper which checks and removes expired sessions periodically.
	if c.quiteC != nil && c.doneC != nil {
		reaper.Quit(c.quiteC, c.doneC)
	}

	if c.Storex.db != nil {
		err = c.Storex.db.Close()
	}

	return
}

func (b *boltStore) Init() (err error) {
	b.once.Do(func() {
		err = b.init()
	})
	return
}

func (b *boltStore) init() (err error) {
	if b.Storex.db != nil {
		b.Close()
	}
	b.Storex.db, err = bolt.Open(b.dbFile, 0666, nil)
	if err != nil {
		return
	}
	b.Storex.Store, err = store.New(b.Storex.db, *b.config, b.keyPairs...)
	if err != nil {
		return
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
	return
}
