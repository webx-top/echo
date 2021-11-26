/*

   Copyright 2017 Wenhui Shen <www.webx.top>

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
package jet

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/admpub/jet/v6"
	"github.com/admpub/log"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/poolx/bufferpool"
)

func init() {
	render.Reg(`jet`, func(tmplDir string) driver.Driver {
		return New(tmplDir)
	})
}

func New(templateDir string, args ...logger.Logger) driver.Driver {
	var err error
	templateDir, err = filepath.Abs(templateDir)
	if err != nil {
		panic(err.Error())
	}
	a := &Jet{
		NopRenderer: &driver.NopRenderer{},
		templateDir: templateDir,
		set:         jet.NewSet(jet.NewOSFileSystemLoader(templateDir)),
	}
	if len(args) > 0 {
		a.logger = args[0]
	} else {
		a.logger = log.New("render-jet")
	}
	return a
}

type Jet struct {
	*driver.NopRenderer
	set           *jet.Set
	templateDir   string
	logger        logger.Logger
	debug         bool
	tmplPathFixer func(echo.Context, string) string
}

func (a *Jet) Debug() bool {
	return a.debug
}

func (a *Jet) SetDebug(on bool) {
	a.debug = on
	a.set.SetDevelopmentMode(on)
}

func (a *Jet) SetLogger(l logger.Logger) {
	a.logger = l
}

func (a *Jet) Logger() logger.Logger {
	return a.logger
}

func (a *Jet) TmplDir() string {
	return a.templateDir
}

func (a *Jet) SetTmplPathFixer(fn func(echo.Context, string) string) {
	a.tmplPathFixer = fn
}

func (a *Jet) TmplPath(c echo.Context, tmpl string) string {
	if a.tmplPathFixer != nil {
		tmpl = a.tmplPathFixer(c, tmpl)
	}
	return tmpl
}

func (a *Jet) SetFuncMap(fn func() map[string]interface{}) {
	for name, fn := range fn() {
		a.set.AddGlobal(name, fn)
	}
}

func (a *Jet) Render(w io.Writer, tmpl string, data interface{}, c echo.Context) error {
	tmpl = a.TmplPath(c, tmpl)
	t, err := a.set.GetTemplate(tmpl)
	if err != nil {
		return err
	}
	vars := make(jet.VarMap)
	for name, fn := range c.Funcs() {
		vars.Set(name, fn)
	}
	return t.Execute(w, vars, data)
}

func (a *Jet) Fetch(tmpl string, data interface{}, c echo.Context) string {
	tmpl = a.TmplPath(c, tmpl)
	w := bufferpool.Get()
	defer bufferpool.Release(w)
	t, err := a.set.GetTemplate(tmpl)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl, err)
	}
	vars := make(jet.VarMap)
	for name, fn := range c.Funcs() {
		vars.Set(name, fn)
	}
	err = t.Execute(w, vars, data)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl, err)
	}
	return com.Bytes2str(w.Bytes())
	//return w.String()
}

func (a *Jet) RawContent(tmpl string) (b []byte, e error) {
	return nil, errors.New(`unsupported`)
}
