package mysql

import (
	"database/sql"
	"encoding/base32"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/admpub/errors"
	"github.com/admpub/null"
	"github.com/admpub/securecookie"
	"github.com/admpub/sessions"
	"github.com/go-sql-driver/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/encoding/dbconfig"
	ss "github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/file"
)

func New(cfg *Options) sessions.Store {
	cfg.Config.Engine = `mysql`
	eng, err := NewMySQLStore(cfg)
	if err != nil {
		log.Println("sessions: Operation MySQL failed:", err)
		return file.NewFilesystemStore(&file.FileOptions{
			SavePath:      ``,
			KeyPairs:      cfg.KeyPairs,
			CheckInterval: cfg.CheckInterval,
		})
	}
	return eng
}

func Reg(store sessions.Store, args ...string) {
	name := `mysql`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *Options, args ...string) sessions.Store {
	store := New(opts)
	Reg(store, args...)
	return store
}

type Options struct {
	Config        dbconfig.Config `json:"-"`
	Table         string          `json:"table"`
	KeyPairs      [][]byte        `json:"-"`
	CheckInterval time.Duration   `json:"checkInterval"`
}

type MySQLStore struct {
	db         *sql.DB
	stmtInsert *sql.Stmt
	stmtDelete *sql.Stmt
	stmtUpdate *sql.Stmt
	stmtSelect *sql.Stmt

	Codecs        []securecookie.Codec
	table         string
	checkInterval time.Duration
	quiteC        chan<- struct{}
	doneC         <-chan struct{}
	once          sync.Once
}

const DDL = "CREATE TABLE IF NOT EXISTS %s (" +
	"	`id` char(64) NOT NULL," +
	"	`data` longblob NOT NULL," +
	"	`created` int(11) unsigned NOT NULL DEFAULT '0'," +
	"	`modified` int(11) unsigned NOT NULL DEFAULT '0'," +
	"	`expires` int(11) unsigned NOT NULL DEFAULT '0'," +
	"	PRIMARY KEY (`id`)" +
	"  ) ENGINE=InnoDB;"

var (
	DefaultMaxAge    = 86400
	DefaultKeyPrefix = `_`
)

type sessionRow struct {
	id       null.String
	data     null.String
	created  null.Int64
	modified null.Int64
	expires  null.Int64
}

// NewMySQLStore takes the following paramaters
// endpoint - A sql.Open style endpoint
// tableName - table where sessions are to be saved. Required fields are created automatically if the table doesnot exist.
// path - path for Set-Cookie header
// maxAge
// codecs
func NewMySQLStore(cfg *Options) (*MySQLStore, error) {
	db, err := sql.Open("mysql", cfg.Config.String())
	if err != nil {
		return nil, err
	}

	return NewMySQLStoreFromConnection(db, cfg)
}

// NewMySQLStoreFromConnection .
func NewMySQLStoreFromConnection(db *sql.DB, cfg *Options) (*MySQLStore, error) {
	// Make sure table name is enclosed.
	tableName := "`" + strings.Trim(cfg.Table, "`") + "`"

	cTableQ := fmt.Sprintf(DDL, tableName)
	if _, err := db.Exec(cTableQ); err != nil {
		switch verr := err.(type) {
		case *mysql.MySQLError:
			// Error 1142 means permission denied for create command
			if verr.Number == 1142 {
				break
			} else {
				return nil, errors.Wrap(err, cTableQ)
			}
		default:
			return nil, err
		}
	}

	insQ := "REPLACE INTO " + tableName +
		"(id, data, created, modified, expires) VALUES (?, ?, ?, ?, ?)"
	stmtInsert, stmtErr := db.Prepare(insQ)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, insQ)
	}

	delQ := "DELETE FROM " + tableName + " WHERE id = ?"
	stmtDelete, stmtErr := db.Prepare(delQ)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, delQ)
	}

	updQ := "UPDATE " + tableName + " SET data = ?, created = ?, expires = ? " +
		"WHERE id = ?"
	stmtUpdate, stmtErr := db.Prepare(updQ)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, updQ)
	}

	selQ := "SELECT id, data, created, modified, expires from " +
		tableName + " WHERE id = ?"
	stmtSelect, stmtErr := db.Prepare(selQ)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, selQ)
	}

	return &MySQLStore{
		db:            db,
		stmtInsert:    stmtInsert,
		stmtDelete:    stmtDelete,
		stmtUpdate:    stmtUpdate,
		stmtSelect:    stmtSelect,
		Codecs:        securecookie.CodecsFromPairs(cfg.KeyPairs...),
		table:         tableName,
		checkInterval: cfg.CheckInterval,
	}, nil
}

func (m *MySQLStore) Close() (err error) {
	m.stmtSelect.Close()
	m.stmtUpdate.Close()
	m.stmtDelete.Close()
	m.stmtInsert.Close()
	err = m.db.Close()
	m.closeCleanup()
	return
}

func (m *MySQLStore) Get(ctx echo.Context, name string) (*sessions.Session, error) {
	m.Init()
	return sessions.GetRegistry(ctx).Get(m, name)
}

func (m *MySQLStore) New(ctx echo.Context, name string) (*sessions.Session, error) {
	session := sessions.NewSession(m, name)
	session.IsNew = true
	var err error
	value := ctx.GetCookie(name)
	if len(value) == 0 {
		return session, err
	}
	err = securecookie.DecodeMulti(name, value, &session.ID, m.Codecs...)
	if err != nil {
		return session, err
	}
	err = m.load(ctx, session)
	if err == nil {
		session.IsNew = false
	} else {
		err = nil
	}
	return session, err
}

func (m *MySQLStore) Reload(ctx echo.Context, session *sessions.Session) error {
	err := m.load(ctx, session)
	if err == nil {
		session.IsNew = false
	} else {
		err = nil
	}
	return err
}

func (m *MySQLStore) Save(ctx echo.Context, session *sessions.Session) error {
	var err error
	// Delete if max-age is < 0
	if ctx.CookieOptions().MaxAge < 0 {
		return m.Delete(ctx, session)
	}
	if len(session.ID) == 0 {
		// generate random session ID key suitable for storage in the db
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
		if err = m.insert(ctx, session); err != nil {
			return err
		}
	} else if err = m.save(ctx, session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, m.Codecs...)
	if err != nil {
		return err
	}
	sessions.SetCookie(ctx, session.Name(), encoded)
	return nil
}

func (m *MySQLStore) Remove(sessionID string) error {
	if len(sessionID) == 0 {
		return nil
	}
	_, delErr := m.stmtDelete.Exec(sessionID)
	if delErr != nil {
		return delErr
	}
	return nil
}

func (m *MySQLStore) insert(ctx echo.Context, session *sessions.Session) error {
	var modifiedAt int64
	var createdAt int64
	var expiredAt int64
	nowTs := time.Now().Unix()
	created := session.Values[DefaultKeyPrefix+"created"]
	if created == nil {
		createdAt = nowTs
	} else {
		createdAt = created.(int64)
	}
	modifiedAt = createdAt
	expires := session.Values[DefaultKeyPrefix+"expires"]
	if expires == nil {
		expiredAt = nowTs + int64(m.maxAge(ctx))
	} else {
		expiredAt = expires.(int64)
	}
	delete(session.Values, DefaultKeyPrefix+"created")
	delete(session.Values, DefaultKeyPrefix+"expires")
	delete(session.Values, DefaultKeyPrefix+"modified")

	encoded, encErr := securecookie.EncodeMulti(session.Name(), session.Values, m.Codecs...)
	if encErr != nil {
		return encErr
	}
	_, insErr := m.stmtInsert.Exec(session.ID, encoded, createdAt, modifiedAt, expiredAt)
	return insErr
}

func (m *MySQLStore) Delete(ctx echo.Context, session *sessions.Session) error {
	sessions.SetCookie(ctx, session.Name(), ``, -1)
	// Clear session values.
	for k := range session.Values {
		delete(session.Values, k)
	}
	return m.Remove(session.ID)
}

func (n *MySQLStore) maxAge(ctx echo.Context) int {
	maxAge := ctx.CookieOptions().MaxAge
	if maxAge == 0 {
		maxAge = DefaultMaxAge
	}
	return maxAge
}

func (m *MySQLStore) save(ctx echo.Context, session *sessions.Session) error {
	if session.IsNew {
		return m.insert(ctx, session)
	}
	var createdAt int64
	var expiredAt int64
	nowTs := time.Now().Unix()
	created := session.Values[DefaultKeyPrefix+"created"]
	if created == nil {
		createdAt = nowTs
	} else {
		createdAt = created.(int64)
	}

	expires := session.Values[DefaultKeyPrefix+"expires"]
	maxAge := int64(m.maxAge(ctx))
	if expires == nil {
		expiredAt = nowTs + maxAge
	} else {
		expiredAt = expires.(int64)
		expiresTs := nowTs + maxAge
		if expiredAt < expiresTs {
			expiredAt = expiresTs
		}
	}

	delete(session.Values, DefaultKeyPrefix+"created")
	delete(session.Values, DefaultKeyPrefix+"expires")
	delete(session.Values, DefaultKeyPrefix+"modified")
	encoded, encErr := securecookie.EncodeMulti(session.Name(), session.Values, m.Codecs...)
	if encErr != nil {
		return encErr
	}
	_, updErr := m.stmtUpdate.Exec(encoded, createdAt, expiredAt, session.ID)
	if updErr != nil {
		return updErr
	}
	return nil
}

func (m *MySQLStore) load(ctx echo.Context, session *sessions.Session) error {
	row := m.stmtSelect.QueryRow(session.ID)
	sess := sessionRow{}
	scanErr := row.Scan(&sess.id, &sess.data, &sess.created, &sess.modified, &sess.expires)
	if scanErr != nil {
		return scanErr
	}
	if sess.expires.Int64 < time.Now().Unix() {
		log.Printf("Session expired on %s, but it is %s now.", time.Unix(sess.expires.Int64, 0), time.Now())
		return errors.New("Session expired")
	}
	err := securecookie.DecodeMulti(session.Name(), sess.data.String, &session.Values, m.Codecs...)
	if err != nil {
		return err
	}
	session.Values[DefaultKeyPrefix+"created"] = sess.created.Int64
	session.Values[DefaultKeyPrefix+"modified"] = sess.modified.Int64
	session.Values[DefaultKeyPrefix+"expires"] = sess.expires.Int64
	return nil

}

func (m *MySQLStore) closeCleanup() {
	// Invoke a reaper which checks and removes expired sessions periodically.
	if m.quiteC != nil && m.doneC != nil {
		m.StopCleanup(m.quiteC, m.doneC)
	}
}

func (m *MySQLStore) Init() {
	m.once.Do(m.init)
}

func (m *MySQLStore) init() {
	m.closeCleanup()
	m.quiteC, m.doneC = m.Cleanup(m.checkInterval)
}
