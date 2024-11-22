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

package engine

import (
	"encoding/base32"
	"strings"
	"sync"

	"github.com/admpub/log"
	"github.com/admpub/securecookie"
	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
)

// Defaults for sessions.Options
const (
	DefaultMaxAge = 86400 * 30 // 30days
)

func NewSession(ctx echo.Context) echo.Sessioner {
	options := ctx.SessionOptions()
	store := StoreEngine(options)
	return NewMySession(store, options.Name, ctx)
}

func NewMySession(store sessions.Store, name string, ctx echo.Context) echo.Sessioner {
	return &Session{
		name:    name,
		context: ctx,
		store:   store,
		session: nil,
		written: false,
	}
}

func StoreEngine(options *echo.SessionOptions) (store sessions.Store) {
	if options == nil {
		return nil
	}
	store = Get(options.Engine)
	if store == nil {
		log.Errorf(`the session storage engine named %q does not exist`, options.Engine)
		if options.Engine != `cookie` {
			store = Get(`cookie`)
			log.Warn(`session uses default storage engine: cookie`)
		}
	}
	return
}

func GenerateSessionID(prefix ...string) string {
	var _prefix string
	if len(prefix) > 0 {
		_prefix = prefix[0]
	}
	return _prefix + strings.TrimRight(
		base32.StdEncoding.EncodeToString(
			securecookie.GenerateRandomKey(32),
		),
		"=",
	)
}

type Stores struct {
	m map[string]sessions.Store
	l sync.RWMutex
}

func (s *Stores) Get(name string) sessions.Store {
	s.l.RLock()
	store, ok := s.m[name]
	s.l.RUnlock()
	if ok {
		return store
	}
	return nil
}

func (s *Stores) Exists(name string) bool {
	s.l.RLock()
	_, ok := s.m[name]
	s.l.RUnlock()
	return ok
}

func (s *Stores) Set(name string, store sessions.Store) {
	s.l.Lock()
	defer s.l.Unlock()
	if old, ok := s.m[name]; ok {
		if c, ok := old.(Closer); ok {
			c.Close()
		}
	}
	s.m[name] = store
}

func (s *Stores) Del(name string) {
	s.l.Lock()
	old, ok := s.m[name]
	if ok {
		if c, ok := old.(Closer); ok {
			c.Close()
		}
		delete(s.m, name)
	}
	s.l.Unlock()
}

var stores = &Stores{
	m: map[string]sessions.Store{},
}

type Closer interface {
	Close() error
}

func Reg(name string, store sessions.Store) {
	stores.Set(name, store)
}

func Get(name string) sessions.Store {
	return stores.Get(name)
}

func Exists(name string) bool {
	return stores.Exists(name)
}

func Del(name string) {
	stores.Del(name)
}
