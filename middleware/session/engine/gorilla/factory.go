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
package session

import (
	"fmt"
	"strconv"

	"github.com/webx-top/echo"
)

func NewSession(options *echo.SessionOptions, setting interface{}, ctx echo.Context) echo.Session {
	store := StoreEngine(options, setting)
	return NewMySession(store, options.Name, ctx)
}

func NewMySession(store Store, name string, ctx echo.Context) echo.Session {
	return &Session{name, ctx, store, nil, false}
}

func StoreEngine(options *echo.SessionOptions, setting interface{}) (store Store) {
	switch options.Engine {
	case `file`:
		s := setting.(map[string]string)
		path, _ := s["path"]
		key, _ := s["key"]
		store = NewFilesystemStore(path, []byte(key))
		store.Options(*options)
	case `redis`:
		s := setting.(map[string]string)
		sizeStr, _ := s["size"]
		network, _ := s["network"]
		address, _ := s["address"]
		password, _ := s["password"]
		key, _ := s["key"]
		size, _ := strconv.Atoi(sizeStr)
		if size < 1 {
			size = 10
		}
		var err error
		store, err = NewRedisStore(size, network, address, password, []byte(key))
		if err != nil {
			fmt.Println(err)
		}
		store.Options(*options)
	case `bolt`:
		s := setting.(map[string]string)
		dbFile, _ := s["file"]
		key, _ := s["key"]
		name, _ := s["name"]
		var bucketName []byte
		if name != `` {
			bucketName = []byte(name)
		}
		var err error
		store, err = NewBoltStore(dbFile, *options, bucketName, []byte(key))
		if err != nil {
			fmt.Println(err)
		}
	case `cookie`:
		fallthrough
	default:
		s := setting.(string)
		store = NewCookieStore([]byte(s))
		store.Options(*options)
	}
	return
}
