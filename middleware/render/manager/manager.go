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

package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/admpub/fsnotify"
	"github.com/admpub/log"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/render/driver"
)

var Default driver.Manager = New()

func New() *Manager {
	m := &Manager{
		caches:   com.InitSafeMap[string, []byte](),
		ignores:  com.InitSafeMap[string, bool](),
		allows:   com.InitSafeMap[string, bool](),
		callback: com.InitSafeMap[string, func(string, string, string)](),
		Logger:   log.GetLogger(`watcher`),
		done:     make(chan bool),
	}
	m.ignores.Set("*.tmp", false)
	m.ignores.Set("*.TMP", false)
	return m
}

// Manager Tempate manager
type Manager struct {
	caches   com.SafeMap[string, []byte] //map[string][]byte
	firstDir atomic.Value
	ignores  com.SafeMap[string, bool] // map[string]bool
	allows   com.SafeMap[string, bool] //map[string]bool
	Logger   logger.Logger
	callback com.SafeMap[string, func(string, string, string)] //map[string]func(string, string, string) //参数为：目标名称，类型(file/dir)，事件名(create/delete/modify/rename)
	done     chan bool
	isClosed atomic.Bool
	watcher  atomic.Value
	once     sync.Once
}

func (m *Manager) closeMoniter() {
	if m.isClosed.Load() {
		return
	}
	m.isClosed.Store(true)
	m.firstDir.Store(``)
	t := time.NewTimer(time.Second * 2)
	defer t.Stop()
	for {
		select {
		case m.done <- true:
			return
		case <-t.C:
			go func() {
				<-m.done
			}()
			return
		}
	}
}

func (m *Manager) getWatcher() (wt *fsnotify.Watcher, err error) {
	if m.isClosed.Load() {
		m.isClosed.Store(false)
		m.once = sync.Once{}
	}
	m.once.Do(func() {
		wt, err = fsnotify.NewWatcher()
		if err != nil {
			m.Logger.Error(err)
		}
		m.watcher.Store(wt)
	})
	wt = m.watcher.Load().(*fsnotify.Watcher)
	return
}

func (m *Manager) AddCallback(rootDir string, callback func(name, typ, event string)) {
	m.callback.Set(rootDir, callback)
}

func (m *Manager) ClearCallback() {
	m.callback.Reset()
}

func (m *Manager) DelCallback(rootDir string) {
	m.callback.Remove(rootDir)
}

func (m *Manager) ClearAllows() {
	m.allows.Reset()
}

func (m *Manager) AddAllow(allows ...string) {
	for _, allow := range allows {
		m.allows.Set(allow, true)
	}
}

func (m *Manager) DelAllow(allow string) {
	m.allows.Remove(allow)
}

func (m *Manager) ClearIgnores() {
	m.ignores.Reset()
}

func (m *Manager) AddIgnore(ignores ...string) {
	for _, ignore := range ignores {
		m.allows.Set(ignore, false)
	}
}

func (m *Manager) DelIgnore(ignore string) {
	m.ignores.Remove(ignore)
}

func (m *Manager) SetLogger(logger logger.Logger) {
	m.Logger = logger
}

func (m *Manager) allowCached(name string) bool {
	ok := m.allows.Size() == 0
	if !ok {
		_, ok = m.allows.GetOk(`*` + filepath.Ext(name))
		if !ok {
			ok = m.allows.Get(filepath.Base(name))
		}
	}
	return ok
}

func (m *Manager) AddWatchDir(ppath string) (err error) {
	if !com.FileExists(ppath) {
		return
	}
	ppath, err = filepath.Abs(ppath)
	if err != nil {
		return
	}
	if v, y := m.firstDir.Load().(string); !y || len(v) == 0 {
		m.firstDir.Store(ppath)
	}
	var watcher *fsnotify.Watcher
	watcher, err = m.getWatcher()
	if err != nil {
		return err
	}
	err = watcher.Add(ppath)
	if err != nil {
		m.Logger.Error(err.Error())
		return
	}

	err = filepath.Walk(ppath, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Add(f)
		}
		return nil
	})

	//err = m.cacheAll(ppath)
	return
}

func (m *Manager) CancelWatchDir(oldDir string) (err error) {
	if !com.FileExists(oldDir) {
		return
	}
	oldDir, err = filepath.Abs(oldDir)
	if err != nil {
		return
	}

	m.caches.ClearEmpty(func(tmpl string, _ []byte) bool {
		return strings.HasPrefix(tmpl, oldDir)
	})

	var watcher *fsnotify.Watcher
	watcher, err = m.getWatcher()
	if err != nil {
		return
	}
	filepath.Walk(oldDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			watcher.Remove(f)
			return nil
		}
		return nil
	})
	watcher.Remove(oldDir)
	return
}

func (m *Manager) ChangeWatchDir(oldDir string, newDir string) (err error) {
	err = m.CancelWatchDir(oldDir)
	if err != nil {
		return err
	}
	err = m.AddWatchDir(newDir)
	return
}

func (m *Manager) Start() error {
	go m.watch()
	return nil
}

func (m *Manager) watch() error {
	watcher, err := m.getWatcher()
	if err != nil {
		return err
	}
	var logSuffix string
	if v, y := m.firstDir.Load().(string); y && len(v) > 0 {
		logSuffix = ": " + v + " etc"
	}
	m.Logger.Debug("TemplateMgr watcher is start" + logSuffix + ".")
	defer watcher.Close()

	for {
		select {
		case ev := <-watcher.Events:
			if _, ok := m.ignores.GetOk(filepath.Base(ev.Name)); ok {
				continue
			}
			if _, ok := m.ignores.GetOk(`*` + filepath.Ext(ev.Name)); ok {
				continue
			}
			d, err := os.Stat(ev.Name)
			if err != nil {
				continue
			}
			if ev.Op&fsnotify.Create == fsnotify.Create {
				if d.IsDir() {
					watcher.Add(ev.Name)
					m.onChange(ev.Name, "dir", "create")
					continue
				}
				m.onChange(ev.Name, "file", "create")
				if m.allowCached(ev.Name) {
					content, err := os.ReadFile(ev.Name)
					if err != nil {
						m.Logger.Infof("loaded template %v failed: %v", ev.Name, err)
						continue
					}
					m.Logger.Infof("loaded template file %v success", ev.Name)
					m.CacheTemplate(ev.Name, content)
				}
			} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
				if d.IsDir() {
					watcher.Remove(ev.Name)
					m.onChange(ev.Name, "dir", "delete")
					continue
				}
				m.onChange(ev.Name, "file", "delete")
				if m.allowCached(ev.Name) {
					m.CacheDelete(ev.Name)
				}
			} else if ev.Op&fsnotify.Write == fsnotify.Write {
				if d.IsDir() {
					m.onChange(ev.Name, "dir", "modify")
					continue
				}
				m.onChange(ev.Name, "file", "modify")
				if m.allowCached(ev.Name) {
					content, err := os.ReadFile(ev.Name)
					if err != nil {
						m.Logger.Errorf("reloaded template %v failed: %v", ev.Name, err)
						continue
					}
					m.CacheTemplate(ev.Name, content)
					m.Logger.Infof("reloaded template %v success", ev.Name)
				}
			} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
				if d.IsDir() {
					watcher.Remove(ev.Name)
					m.onChange(ev.Name, "dir", "rename")
					continue
				}
				m.onChange(ev.Name, "file", "rename")
				if m.allowCached(ev.Name) {
					m.CacheDelete(ev.Name)
				}
			}
		case err := <-watcher.Errors:
			if err != nil {
				m.Logger.Error("error:", err)
			}
		case <-m.done:
			goto END
		}
	}

END:
	m.Logger.Debug("TemplateMgr watcher is closed" + logSuffix + ".")
	return nil
}

func (m *Manager) onChange(name, typ, event string) {
	m.callback.Range(func(key string, callback func(string, string, string)) bool {
		callback(name, typ, event)
		return true
	})
}

func (m *Manager) cacheAll(rootDir string) error {
	fmt.Print(rootDir + ": Reading the contents of the template files, please wait... ")
	err := filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, ok := m.ignores.GetOk(filepath.Base(f)); !ok {
			content, err := os.ReadFile(f)
			if err != nil {
				m.Logger.Debugf("load template %s error: %v", f, err)
				return err
			}
			m.Logger.Debugf("loaded template", f)
			m.caches.Set(f, content)
		}
		return nil
	})
	fmt.Println(rootDir + ": Complete.")
	return err
}

func (m *Manager) Close() {
	m.closeMoniter()
}

func (m *Manager) GetTemplate(tmpl string) ([]byte, error) {
	tmplPath, err := filepath.Abs(tmpl)
	if err != nil {
		return nil, err
	}

	if !m.allowCached(tmplPath) {
		return os.ReadFile(tmplPath)
	}
	if content, ok := m.caches.GetOk(tmplPath); ok {
		m.Logger.Debugf("load template %v from cache", tmplPath)
		return content, nil
	}
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}
	m.Logger.Debugf("load template %v from the file", tmplPath)
	m.caches.Set(tmplPath, content)
	return content, err
}

func (m *Manager) SetTemplate(tmpl string, content []byte) error {
	tmplPath, err := filepath.Abs(tmpl)
	if err != nil {
		return err
	}

	err = os.WriteFile(tmplPath, content, 0666)
	if err != nil {
		return err
	}
	if chmodErr := os.Chmod(tmplPath, 0666); chmodErr != nil {
		m.Logger.Error(`%s: %s`, tmplPath, chmodErr.Error())
	}
	if m.allowCached(tmplPath) {
		m.Logger.Debugf("load template %v from the file", tmplPath)
		m.caches.Set(tmplPath, content)
	}
	return err
}

func (m *Manager) CacheTemplate(tmpl string, content []byte) {
	m.Logger.Debugf("update template %v on cache", tmpl)
	m.caches.Set(tmpl, content)
}

func (m *Manager) CacheDelete(tmpl string) {
	if m.caches.Exists(tmpl) {
		m.Logger.Infof("delete template %v from cache", tmpl)
		m.caches.Delete(tmpl)
	}
}

func (m *Manager) ClearCache() {
	m.caches.Reset()
}
