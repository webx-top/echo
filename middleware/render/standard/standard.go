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

/**
 * 模板扩展
 * @author swh <swh@admpub.com>
 */
package standard

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/admpub/log"
	"golang.org/x/sync/singleflight"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/render/manager"
	"github.com/webx-top/poolx/bufferpool"
)

var Debug = false

func New(templateDir string, args ...logger.Logger) driver.Driver {
	var err error
	templateDir, err = filepath.Abs(templateDir)
	if err != nil {
		panic(err.Error())
	}
	t := &Standard{
		cache:             com.InitSafeMap[string, *CacheData](),
		TemplateDir:       templateDir,
		DelimLeft:         "{{",
		DelimRight:        "}}",
		IncludeTag:        "Include",
		SnippetTag:        "Snippet",
		ExtendTag:         "Extend",
		BlockTag:          "Block",
		SuperTag:          "Super",
		StripTag:          "Strip",
		Ext:               ".html",
		debug:             Debug,
		fileEvents:        make([]func(string), 0),
		contentProcessors: make([]func(tmpl string, content []byte) []byte, 0),
	}
	if len(args) > 0 {
		t.logger = args[0]
	} else {
		t.logger = log.New("render-standard")
	}
	t.InitRegexp()
	t.SetManager(manager.Default)
	return t
}

type Standard struct {
	cache              com.SafeMap[string, *CacheData]
	TemplateDir        string
	TemplateMgr        driver.Manager
	contentProcessors  []func(tmpl string, content []byte) []byte
	DelimLeft          string
	DelimRight         string
	incTagRegex        *regexp.Regexp
	funcTagRegex       *regexp.Regexp
	extTagRegex        *regexp.Regexp
	blkTagRegex        *regexp.Regexp
	rplTagRegex        *regexp.Regexp
	innerTagBlankRegex *regexp.Regexp
	stripTagRegex      *regexp.Regexp
	IncludeTag         string
	SnippetTag         string
	ExtendTag          string
	BlockTag           string
	SuperTag           string
	StripTag           string
	Ext                string
	tmplPathFixer      func(echo.Context, string) string
	debug              bool
	getFuncs           func() map[string]interface{}
	logger             logger.Logger
	fileEvents         []func(string)
	quotedLeft         string
	quotedRight        string
	quotedRfirst       string
	sg                 singleflight.Group
}

func (a *Standard) Debug() bool {
	return a.debug
}

func (a *Standard) SetDebug(on bool) {
	a.debug = on
}

func (a *Standard) SetLogger(l logger.Logger) {
	a.logger = l
	if a.TemplateMgr != nil {
		a.TemplateMgr.SetLogger(a.logger)
	}
}
func (a *Standard) Logger() logger.Logger {
	return a.logger
}

func (a *Standard) TmplDir() string {
	return a.TemplateDir
}

func (a *Standard) MonitorEvent(fn func(string)) {
	if fn == nil {
		return
	}
	a.fileEvents = append(a.fileEvents, fn)
}

func (a *Standard) SetContentProcessor(fn func(tmpl string, content []byte) []byte) {
	if fn == nil {
		return
	}
	a.contentProcessors = append(a.contentProcessors, fn)
}

func (a *Standard) SetFuncMap(fn func() map[string]interface{}) {
	a.getFuncs = fn
}

func (a *Standard) deleteCaches(_ string) {
	a.cache.Reset()
	a.logger.Info("remove cached template object")
}

func (a *Standard) Init() {
	a.InitRegexp()
	callback := func(name, typ, event string) {
		switch event {
		case "create":
		case "delete", "modify", "rename":
			if typ == "dir" {
				return
			}
			a.deleteCaches(name)
			for _, fn := range a.fileEvents {
				fn(name)
			}
		}
	}
	a.TemplateMgr.AddAllow("*" + a.Ext)
	a.TemplateMgr.AddWatchDir(a.TemplateDir)
	a.TemplateMgr.AddCallback(a.TemplateDir, callback)
	a.TemplateMgr.Start()
}

func (a *Standard) SetManager(mgr driver.Manager) {
	if a.TemplateMgr != nil {
		a.TemplateMgr.Close()
	}
	a.TemplateMgr = mgr
}

func (a *Standard) Manager() driver.Manager {
	return a.TemplateMgr
}

func (a *Standard) SetTmplPathFixer(fn func(echo.Context, string) string) {
	a.tmplPathFixer = fn
}

func (a *Standard) TmplPath(c echo.Context, p string) string {
	if a.tmplPathFixer != nil {
		return a.tmplPathFixer(c, p)
	}
	p = filepath.Join(a.TemplateDir, p)
	return p
}

var (
	quoteDelim = "[\"`]"
	quoteInner = "([^\"`]+)"
	quoteRegex = quoteDelim + quoteInner + quoteDelim
)

func (a *Standard) InitRegexp() {
	a.quotedLeft = regexp.QuoteMeta(a.DelimLeft)
	a.quotedRight = regexp.QuoteMeta(a.DelimRight)
	a.quotedRfirst = regexp.QuoteMeta(a.DelimRight[0:1])

	//{{Include "tmpl"}} or {{Include "tmpl" .}}
	a.incTagRegex = regexp.MustCompile(a.quotedLeft + a.IncludeTag + `[\s]+` + quoteRegex + `(?:[\s]+([^` + a.quotedRfirst + `]+))?[\s]*\/?` + a.quotedRight)

	//{{Snippet "funcName"}} or {{Snippet "funcName" .}}
	a.funcTagRegex = regexp.MustCompile(a.quotedLeft + a.SnippetTag + `[\s]+` + quoteRegex + `(?:[\s]+([^` + a.quotedRfirst + `]+))?[\s]*\/?` + a.quotedRight)

	//{{Extend "name"}}
	a.extTagRegex = regexp.MustCompile(`^[\s]*` + a.quotedLeft + a.ExtendTag + `[\s]+` + quoteRegex + `(?:[\s]+([^` + a.quotedRfirst + `]+))?[\s]*\/?` + a.quotedRight)

	//{{Block "name"}}content{{/Block}}
	a.blkTagRegex = regexp.MustCompile(`(?s)` + a.quotedLeft + a.BlockTag + `[\s]+` + quoteRegex + `[\s]*` + a.quotedRight + `(.*?)` + a.quotedLeft + `\/` + a.BlockTag + a.quotedRight)

	//{{Block "name"/}}
	a.rplTagRegex = regexp.MustCompile(a.quotedLeft + a.BlockTag + `[\s]+` + quoteRegex + `[\s]*\/` + a.quotedRight)

	//}}...{{ or >...<
	a.innerTagBlankRegex = regexp.MustCompile(`(?s)(` + a.quotedRight + `|>)[\s]{2,}(` + a.quotedLeft + `|<)`)

	//{{Strip}}...{{/Strip}}
	a.stripTagRegex = regexp.MustCompile(`(?s)` + a.quotedLeft + a.StripTag + a.quotedRight + `(.*?)` + a.quotedLeft + `\/` + a.StripTag + a.quotedRight)
}

// Render HTML
func (a *Standard) Render(w io.Writer, tmplName string, values interface{}, c echo.Context) error {
	return a.RenderBy(w, tmplName, a.RawContent, values, c)
}

// RenderBy render by content
func (a *Standard) RenderBy(w io.Writer, tmplName string, tmplContent func(string) ([]byte, error), values interface{}, c echo.Context) error {
	tmpl, err := a.parse(c, tmplName, tmplContent)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, tmpl.Name(), values)
}

func (a *Standard) parse(c echo.Context, tmplName string, tmplContent func(string) ([]byte, error)) (tmpl *template.Template, err error) {
	tmplOriginalName := tmplName
	tmplName = tmplName + a.Ext
	tmplName = a.TmplPath(c, tmplName)
	cachedKey := tmplName
	cachedData, ok := a.cache.GetOk(cachedKey)
	if ok {
		tmpl = cachedData.template
		return
	}
	var v interface{}
	v, err, _ = a.sg.Do(cachedKey, func() (interface{}, error) {
		var funcMap template.FuncMap
		if a.getFuncs != nil {
			funcMap = template.FuncMap(a.getFuncs())
		}
		if funcMap == nil {
			funcMap = template.FuncMap{}
		}
		return a.find(c, tmplOriginalName, tmplName, tmplContent, cachedKey, funcMap)
	})
	if err != nil {
		return
	}
	tmpl = v.(*template.Template)
	return
}

var bytesBOM = []byte("\xEF\xBB\xBF")

func (a *Standard) find(c echo.Context,
	tmplOriginalName string, tmplName string, tmplContent func(string) ([]byte, error),
	cachedKey string, funcMap template.FuncMap) (tmpl *template.Template, err error) {
	if a.debug {
		start := time.Now()
		a.logger.Debug(` ◐ compile template: `, tmplName)
		defer func() {
			a.logger.Debug(` ◑ finished compile: `+tmplName, ` (elapsed: `+time.Since(start).String()+`)`)
		}()
	}
	tmpl = template.New(driver.CleanTemplateName(tmplName))
	tmpl.Delims(a.DelimLeft, a.DelimRight)
	cacheData := NewCache(tmpl)
	funcMap = cacheData.setFunc(funcMap)
	tmpl.Funcs(funcMap)
	var b []byte
	b, err = tmplContent(tmplName)
	if err != nil {
		return
	}
	content := string(b)
	subcs := map[string]string{} //子模板内容
	extcs := map[string]string{} //母板内容
	m := a.extTagRegex.FindAllStringSubmatch(content, 1)
	content = a.rplTagRegex.ReplaceAllString(content, ``)
	for i := 0; i < 10 && len(m) > 0; i++ {
		a.ParseBlock(c, content, subcs, extcs)
		extFile := m[0][1] + a.Ext
		passObject := m[0][2]
		extFile = a.TmplPath(c, extFile)
		b, err = a.RawContent(extFile)
		if err != nil {
			err = parseError(err, string(b))
			return
		}
		content = string(b)
		content, m = a.ParseExtend(c, content, extcs, passObject, subcs)
	}
	content = a.ContainsSubTpl(c, content, subcs)
	clips := map[string]string{}
	content = a.ContainsSnippetResult(c, tmplOriginalName, content, clips)
	tmpl, err = tmpl.Parse(content)
	if err != nil {
		err = parseError(err, content)
		return
	}

	var defines string

	// include
	for name, subc := range subcs {
		subc = a.ContainsSnippetResult(c, tmplOriginalName, subc, clips)
		subc = a.Tag(`define "`+driver.CleanTemplateName(name)+`"`) + subc + a.Tag(`end`)
		defines += subc
	}

	// block
	for name, extc := range extcs {
		extc = a.ContainsSnippetResult(c, tmplOriginalName, extc, clips)
		extc = a.Tag(`define "`+driver.CleanTemplateName(name)+`"`) + extc + a.Tag(`end`)
		defines += extc
		cacheData.blocks[name] = struct{}{}
	}

	// parse define...
	tmpl, err = tmpl.Parse(defines)
	if err != nil {
		err = parseError(err, defines)
		return
	}
	a.cache.Set(cachedKey, cacheData)
	return
}

func (a *Standard) Fetch(tmplName string, data interface{}, c echo.Context) string {
	content, err := a.parse(c, tmplName, a.RawContent)
	if err != nil {
		return err.Error()
	}
	return a.execute(content, data)
}

func (a *Standard) execute(tmpl *template.Template, data interface{}) string {
	buf := bufferpool.Get()
	defer bufferpool.Release(buf)
	err := tmpl.ExecuteTemplate(buf, tmpl.Name(), data)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err)
	}
	return com.Bytes2str(buf.Bytes())
}

func (a *Standard) ParseBlock(c echo.Context, content string, subcs map[string]string, extcs map[string]string) {
	matches := a.blkTagRegex.FindAllStringSubmatchIndex(content, -1)
	for _, v := range matches {
		var blockName, innerData string
		com.GetMatchedByIndex(content, v, nil, &blockName, &innerData)
		innerData = trimSpaceInBlock(innerData)
		extcs[blockName] = a.ContainsSubTpl(c, innerData, subcs)
	}
}

func (a *Standard) ParseExtend(c echo.Context, content string, extcs map[string]string, passObject string, subcs map[string]string) (string, [][]string) {
	m := a.extTagRegex.FindAllStringSubmatch(content, 1)
	hasParent := len(m) > 0
	if len(passObject) == 0 {
		passObject = "."
	}
	content = a.rplTagRegex.ReplaceAllStringFunc(content, func(match string) string {
		blockName := match[strings.Index(match, `"`)+1:]
		blockName = blockName[0:strings.Index(blockName, `"`)]
		if v, ok := extcs[blockName]; ok {
			return v
		}
		return ``
	})
	matches := a.blkTagRegex.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return content, m
	}
	var superTag string
	if len(a.SuperTag) > 0 {
		superTag = a.Tag(a.SuperTag)
	}
	rec := make(map[string]uint8)
	sup := make(map[string]string)
	var replaced string
	fn := com.ReplaceByMatchedIndex(content, matches, &replaced)
	for k, v := range matches {
		var blockName, innerStr string
		com.GetMatchedByIndex(content, v, nil, &blockName, &innerStr)
		innerStr = trimSpaceInBlock(innerStr)
		if val, ok := extcs[blockName]; ok {
			var suffix string
			if idx, ok := rec[blockName]; ok {
				idx++
				rec[blockName] = idx
				suffix = fmt.Sprintf(`.%v`, idx)
			} else {
				rec[blockName] = 0
			}
			if len(superTag) > 0 {
				sv, hasSuper := sup[blockName]
				if !hasSuper {
					hasSuper = strings.Contains(val, superTag)
					if hasSuper {
						sup[blockName] = val
					}
				} else {
					val = sv
				}
				if hasSuper {
					innerStr = a.ContainsSubTpl(c, innerStr, subcs)
					val = strings.Replace(val, superTag, innerStr, 1)
					if suffix == `` {
						extcs[blockName] = val
					}
				}
			}
			if len(suffix) > 0 {
				extcs[blockName+suffix] = val
				rec[blockName+suffix] = 0
			}
			if hasParent {
				innerStr = a.DelimLeft + a.BlockTag + ` "` + blockName + `"` + a.DelimRight + val + a.DelimLeft + `/` + a.BlockTag + a.DelimRight
			} else {
				innerStr = a.Tag(`template "` + blockName + suffix + `" ` + passObject)
			}
		} else {
			if hasParent {
				fn(k, v)
				continue
			}
		}
		fn(k, v, innerStr)
	}
	//只保留layout中存在的Block
	for k := range extcs {
		if _, ok := rec[k]; !ok {
			delete(extcs, k)
		}
	}
	return replaced, m
}

func (a *Standard) ContainsSubTpl(c echo.Context, content string, subcs map[string]string) string {
	matches := a.incTagRegex.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return content
	}
	var replaced string
	fn := com.ReplaceByMatchedIndex(content, matches, &replaced)
	for k, v := range matches {
		var tmplFile, passObject string
		com.GetMatchedByIndex(content, v, nil, &tmplFile, &passObject)
		tmplFile += a.Ext
		tmplFile = a.TmplPath(c, tmplFile)
		if _, ok := subcs[tmplFile]; !ok {
			// if v, ok := a.CachedRelation[tmplFile]; ok && v.Tpl[1] != nil {
			// 	subcs[tmplFile] = ""
			// } else {
			b, err := a.RawContent(tmplFile)
			if err != nil {
				return fmt.Sprintf("RenderTemplate %v read err: %s", tmplFile, err)
			}
			str := string(b)
			subcs[tmplFile] = "" //先登记，避免死循环
			str = a.ContainsSubTpl(c, str, subcs)
			subcs[tmplFile] = str
			//}
		}
		if len(passObject) == 0 {
			passObject = "."
		}
		fn(k, v, a.Tag(`template "`+driver.CleanTemplateName(tmplFile)+`" `+passObject))
	}
	return replaced
}

func (a *Standard) ContainsSnippetResult(c echo.Context, tmplOriginalName string, content string, clips map[string]string) string {
	matches := a.funcTagRegex.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return content
	}
	var replaced string
	fn := com.ReplaceByMatchedIndex(content, matches, &replaced)
	for k, v := range matches {
		var funcName, passArg string
		com.GetMatchedByIndex(content, v, nil, &funcName, &passArg)
		key := funcName + `:` + passArg
		if _, ok := clips[key]; !ok {
			switch fn := c.GetFunc(funcName).(type) {
			case func(echo.Context, string, string) string:
				clips[key] = fn(c, tmplOriginalName, passArg)
			case func(string, string) string:
				clips[key] = fn(tmplOriginalName, passArg)
			default:
				clips[key] = ``
			}
		}
		fn(k, v, clips[key])
	}
	return replaced
}

func (a *Standard) Tag(content string) string {
	return a.DelimLeft + content + a.DelimRight
}

func (a *Standard) preprocess(tmpl string, b []byte) []byte {
	if b == nil {
		return nil
	}
	if a.contentProcessors != nil {
		for _, fn := range a.contentProcessors {
			b = fn(tmpl, b)
		}
	}
	return a.strip(b)
}

func (a *Standard) RawContent(tmpl string) (b []byte, e error) {
	if a.TemplateMgr != nil {
		b, e = a.TemplateMgr.GetTemplate(tmpl)
	} else {
		b, e = os.ReadFile(tmpl)
	}
	if e != nil {
		return
	}
	b = bytes.TrimPrefix(b, bytesBOM)
	b = a.preprocess(tmpl, b)
	return
}

func (a *Standard) strip(src []byte) []byte {
	if a.debug {
		src = bytes.ReplaceAll(src, []byte(a.DelimLeft+a.StripTag+a.DelimRight), []byte{})
		return bytes.ReplaceAll(src, []byte(a.DelimLeft+`/`+a.StripTag+a.DelimRight), []byte{})
	}
	src = a.stripTagRegex.ReplaceAllFunc(src, func(b []byte) []byte {
		b = bytes.TrimPrefix(b, []byte(a.DelimLeft+a.StripTag+a.DelimRight))
		b = bytes.TrimSuffix(b, []byte(a.DelimLeft+`/`+a.StripTag+a.DelimRight))
		var pres [][]byte
		b, pres = driver.ReplacePRE(b)
		b = a.innerTagBlankRegex.ReplaceAll(b, driver.FE)
		b = driver.RemoveMultiCRLF(b)
		b = bytes.TrimSpace(b)
		b = driver.RecoveryPRE(b, pres)
		return b
	})
	return src
}

func (a *Standard) stripSpace(b []byte) []byte {
	var pres [][]byte
	b, pres = driver.ReplacePRE(b)
	b = a.innerTagBlankRegex.ReplaceAll(b, driver.FE)
	b = bytes.TrimSpace(b)
	b = driver.RecoveryPRE(b, pres)
	return b
}

func (a *Standard) ClearCache() {
	if a.TemplateMgr != nil {
		a.TemplateMgr.ClearCache()
	}
	a.cache.Reset()
}

func (a *Standard) Close() {
	a.ClearCache()
	if a.TemplateMgr != nil {
		if a.TemplateMgr == manager.Default {
			a.TemplateMgr.CancelWatchDir(a.TemplateDir)
			a.TemplateMgr.DelCallback(a.TemplateDir)
		} else {
			a.TemplateMgr.Close()
		}
	}
}
