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
	callback := &com.MonitorEvent{
		Modify: onchange,
		Delete: onchange,
		Rename: onchange,
	}
	callback.Watch(func(f string) bool {
		log.Info("changed language: ", f)
		return strings.HasSuffix(f, `.yaml`)
	})
	for _, mp := range a.config.MessagesPath {
		if len(mp) == 0 {
			continue
		}
		if err := callback.AddDir(mp); err != nil {
			log.Debugf(`failed to I18n.Monitor.AddDir(%q): %v`, mp, err)
		}
	}
	return a
}
