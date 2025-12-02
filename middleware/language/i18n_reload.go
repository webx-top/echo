package language

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"golang.org/x/sync/singleflight"
)

// safeReload safely reloads the specified language file with panic recovery.
// It logs the reload operation and returns any error encountered during reload,
// including recovered panics converted to errors.
// The file parameter specifies the path to the language file to reload.
func (a *I18n) safeReload(file string) error {
	log.Info("reload language: ", file)
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(`%v`, e)
		}
	}()
	a.Reload(file)
	return err
}

// Reload reloads the translator for the specified language code.
// If the langCode ends with ".yaml", it will be trimmed and only the base name will be used.
// This method also removes the cached translator for the language.
func (a *I18n) Reload(langCode string) {
	if strings.HasSuffix(langCode, `.yaml`) {
		langCode = strings.TrimSuffix(langCode, `.yaml`)
		langCode = filepath.Base(langCode)
	}
	a.TranslatorFactory.Reload(langCode)

	a.lock.Lock()
	delete(a.translators, langCode)
	a.lock.Unlock()
}

// Monitor starts watching language files for changes and automatically reloads them when modified.
// It uses a singleflight group to prevent concurrent reloads of the same file.
// The method watches for modify, delete and rename events on .yaml files in configured MessagesPath directories.
// Returns the I18n instance for method chaining.
func (a *I18n) Monitor() *I18n {
	reload := func(file string) error {
		err := a.safeReload(file)
		if err == nil {
			return err
		}
		log.Warnf(`failed to reload language %s: %v`, file, err)
		log.Infof(`start retrying to load the language file %s`, file)
		time.Sleep(time.Second)
		err = a.safeReload(file)
		if err != nil {
			log.Errorf(`failed to reload language %s: %v`, file, err)
		}
		return err
	}

	sg := singleflight.Group{}

	onchange := func(file string) {
		sg.Do(file, func() (interface{}, error) {
			return nil, reload(file)
		})
	}
	if a.monitor != nil {
		a.monitor.Close()
	}
	a.monitor = &com.MonitorEvent{
		Modify: onchange,
		Delete: onchange,
		Rename: onchange,
	}
	a.monitor.Watch(func(f string) bool {
		log.Info("changed language: ", f)
		return strings.HasSuffix(f, `.yaml`)
	})
	for _, mp := range a.config.MessagesPath {
		if len(mp) == 0 {
			continue
		}
		if err := a.monitor.AddDir(mp); err != nil {
			log.Debugf(`failed to I18n.Monitor.AddDir(%q): %v`, mp, err)
		}
	}
	return a
}
