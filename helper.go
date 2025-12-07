package echo

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
)

var workDir string

func SetWorkDir(dir string) {
	if len(dir) == 0 {
		if len(workDir) == 0 {
			setWorkDir()
		}
		return
	}
	if !strings.HasSuffix(dir, FilePathSeparator) {
		dir += FilePathSeparator
	}
	workDir = dir
}

func setWorkDir() {
	workDir, _ = os.Getwd()
	workDir = workDir + FilePathSeparator
}

func init() {
	if len(workDir) == 0 {
		setWorkDir()
	}
}

func Wd() string {
	if len(workDir) == 0 {
		setWorkDir()
	}
	return workDir
}

// HandlerName returns the handler name
func HandlerName(h interface{}) string {
	if h == nil {
		return `<nil>`
	}
	v := reflect.ValueOf(h)
	t := v.Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(v.Pointer()).Name()
	}
	return t.String()
}

// HandlerPath returns the handler path
func HandlerPath(h interface{}) string {
	v := reflect.ValueOf(h)
	t := v.Type()
	switch t.Kind() {
	case reflect.Func:
		return runtime.FuncForPC(v.Pointer()).Name()
	case reflect.Pointer:
		t = t.Elem()
		fallthrough
	case reflect.Struct:
		return t.PkgPath() + `.` + t.Name()
	}
	return ``
}

func HandlerTmpl(handlerPath string) string {
	name := path.Base(handlerPath)
	var r []string
	var u []rune
	for _, b := range name {
		switch b {
		case '*', '(', ')':
			continue
		case '-':
			goto END
		case '.':
			r = append(r, string(u))
			u = []rune{}
		default:
			u = append(u, b)
		}
	}

END:
	if len(u) > 0 {
		r = append(r, string(u))
	}
	for i, s := range r {
		r[i] = com.SnakeCase(s)
	}
	return `/` + strings.Join(r, `/`)
}

// Methods returns methods
func Methods() []string {
	return methods
}

// ContentTypeByExtension returns the MIME type associated with the file based on
// its extension. It returns `application/octet-stream` incase MIME type is not
// found.
func ContentTypeByExtension(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); len(t) == 0 {
		t = MIMEOctetStream
	}
	return
}

func static(r RouteRegister, prefix, root string) {
	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	h := func(c Context) error {
		name := filepath.Join(root, c.Param("*"))
		if !strings.HasPrefix(name, root) {
			return ErrNotFound
		}
		return c.File(name)
	}
	if prefix == "/" {
		r.Get(prefix+"*", h)
	} else {
		r.Get(prefix+"/*", h)
	}
}

// Clear removes the specified items from the slice.
func Clear[T comparable](old []T, clears ...T) []T {
	if len(clears) == 0 {
		return nil
	}
	if len(old) == 0 {
		return old
	}
	result := []T{}
	for _, el := range old {
		var exists bool
		for _, d := range clears {
			if d == el {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, el)
		}
	}
	return result
}

// Dump 输出对象和数组的结构信息
func Dump(m interface{}, printOrNot ...bool) (r string) {
	v, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	r = string(v)
	l := len(printOrNot)
	if l < 1 || printOrNot[0] {
		fmt.Println(r)
	}
	return
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func LogIf(err error, types ...string) {
	if err == nil {
		return
	}
	var typ string
	if len(types) > 0 {
		typ = types[0]
	}
	typ = strings.Title(typ)
	switch typ {
	case `Fatal`:
		log.Fatal(err)
	case `Warn`:
		log.Debug(err)
	case `Debug`:
		log.Debug(err)
	case `Info`:
		log.Info(err)
	default:
		log.Error(err)
	}
}

func URLEncode(s string, rfc ...bool) string {
	encoded := url.QueryEscape(s)
	if len(rfc) > 0 && rfc[0] { // RFC 3986
		encoded = strings.Replace(encoded, `+`, `%20`, -1)
	}
	return encoded
}

func URLDecode(encoded string, rfc ...bool) (string, error) {
	if len(rfc) > 0 && rfc[0] {
		encoded = strings.Replace(encoded, `%20`, `+`, -1)
	}
	return url.QueryUnescape(encoded)
}

func SetAttachmentHeader(c Context, name string, setContentType bool, inline ...bool) {
	var typ string
	if len(inline) > 0 && inline[0] {
		typ = `inline`
	} else {
		typ = `attachment`
	}
	if setContentType {
		c.Response().Header().Set(HeaderContentType, ContentTypeByExtension(name))
	}
	encodedName := URLEncode(name, true)
	c.Response().Header().Set(HeaderContentDisposition, typ+"; filename="+encodedName+"; filename*=utf-8''"+encodedName)
}

func InSliceFold(value string, items []string) bool {
	for _, item := range items {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}

type HandlerFuncs map[string]func(Context) error

func (h *HandlerFuncs) Register(key string, fn func(Context) error) {
	(*h)[key] = fn
}

func (h *HandlerFuncs) Unregister(keys ...string) {
	for _, key := range keys {
		delete(*h, key)
	}
}

func (h HandlerFuncs) Call(c Context, key string) error {
	fn, ok := h[key]
	if !ok {
		return ErrNotFound
	}
	return fn(c)
}

// https://developers.google.cn/search/docs/crawling-indexing/block-indexing?hl=zh-cn#http-response-header
func SearchEngineNoindex(c Context) {
	c.Response().Header().Set(`X-Robots-Tag`, `noindex`)
}

var DefaultNextURLVarName = `next`

func GetNextURL(ctx Context, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	next := ctx.Form(varName)
	if strings.HasPrefix(next, `/`) {
		if next == ctx.Request().URL().Path() {
			next = ``
		} else if pos := strings.LastIndex(next, varName+`=`); pos > -1 {
			next = next[pos+len(varName+`=`):]
		}
	}
	return next
}

func ReturnToCurrentURL(ctx Context, varNames ...string) string {
	next := GetNextURL(ctx, varNames...)
	if len(next) == 0 {
		next = ctx.Request().URI()
	}
	return next
}

func WithNextURL(ctx Context, urlStr string, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	withVarName := varName
	if len(varNames) > 1 && len(varNames[1]) > 0 {
		withVarName = varNames[1]
	}

	next := GetNextURL(ctx, varName)
	if len(next) == 0 || next == urlStr {
		return urlStr
	}
	if next[0] == '/' {
		if len(urlStr) > 8 {
			var urlCopy string
			switch strings.ToLower(urlStr[0:7]) {
			case `https:/`:
				urlCopy = urlStr[8:]
			case `http://`:
				urlCopy = urlStr[7:]
			}
			if len(urlCopy) > 0 {
				p := strings.Index(urlCopy, `/`)
				if p > 0 && urlCopy[p:] == next {
					return urlStr
				}
			}
		}
	}

	return com.WithURLParams(urlStr, withVarName, next)
}

func GetOtherURL(ctx Context, next string) string {
	if len(next) == 0 {
		return next
	}
	urlInfo, _ := url.Parse(next)
	if urlInfo == nil || (urlInfo.Hostname() == ctx.Host() && urlInfo.Path == ctx.Request().URL().Path()) {
		next = ``
	}
	return next
}

var regErrorTemplateFile = regexp.MustCompile(`template: ([^:]+)\:([\d]+)\:(?:([\d]+)\:)? `)

func ParseTemplateError(err error, sourceContent string) *PanicError {
	content := err.Error()
	p := NewPanicError(content, err)
	matches := regErrorTemplateFile.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		line, _ := strconv.Atoi(match[2])
		t := &Trace{
			File:   match[1],
			Line:   line,
			Func:   ``,
			HasErr: true,
		}
		p.AddTrace(t, sourceContent)
	}
	return p
}

// AddExtension adds the default extension to the URI if it does not already have it.
func AddExtension(c Context, uri string) string {
	if len(uri) == 0 || strings.HasSuffix(uri, `/`) || len(c.DefaultExtension()) == 0 {
		return uri
	}
	parts := strings.SplitN(uri, `?`, 2)
	if !strings.HasSuffix(parts[0], c.DefaultExtension()) {
		parts[0] += c.DefaultExtension()
		uri = strings.Join(parts, `?`)
	}
	return uri
}

func CleanPath(ppath string) string {
	if !strings.HasPrefix(ppath, `/`) {
		ppath = `/` + ppath
	}
	return path.Clean(ppath)
}

func CleanFilePath(ppath string) string {
	if !strings.HasPrefix(ppath, FilePathSeparator) {
		ppath = FilePathSeparator + ppath
	}
	return filepath.Clean(ppath)
}

var pathWithDots = regexp.MustCompile(`(?:^\.\.[/\\]|[/\\]\.\.[/\\]|[/\\]\.\.$|^\.\.$)`)
var (
	ErrInvalidPathTraversal = errors.New("invalid path traversal")
	ErrPathEscapesBase      = errors.New("path escapes base")
)

func PathHasDots(ppath string) bool {
	return pathWithDots.MatchString(ppath)
}

func PathSafely(reqPath string) (string, error) {
	cleaned := CleanPath(reqPath)
	if pathWithDots.MatchString(cleaned) {
		return cleaned, fmt.Errorf(`%w: %s`, ErrInvalidPathTraversal, reqPath)
	}
	return cleaned, nil
}

func FilePathSafely(reqPath string) (string, error) {
	cleaned := CleanFilePath(reqPath)
	if pathWithDots.MatchString(cleaned) {
		return cleaned, fmt.Errorf(`%w: %s`, ErrInvalidPathTraversal, reqPath)
	}
	return cleaned, nil
}

func FilePathJoin(base, reqPath string) (string, error) {
	cleaned, err := FilePathSafely(reqPath)
	if err != nil {
		return ``, err
	}
	full := filepath.Join(base, cleaned)
	cleanedBase := filepath.Clean(base)
	// Ensure the resolved path is under the base directory
	if !strings.HasPrefix(full, cleanedBase+FilePathSeparator) {
		return ``, fmt.Errorf(`%w(%s): %s`, ErrPathEscapesBase, base, full)
	}
	return full, nil
}

func CreateInRoot(dir, name string) (*os.File, error) {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return r.Create(name)
}

func WriteFileInRoot(dir, name string, data []byte, perm os.FileMode) error {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer r.Close()
	return r.WriteFile(name, data, perm)
}

func ReadFileInRoot(dir, name string) ([]byte, error) {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return r.ReadFile(name)
}

func RemoveInRoot(dir, name string) error {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer r.Close()
	return r.Remove(name)
}

func RemoveAllInRoot(dir, name string) error {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer r.Close()
	return r.RemoveAll(name)
}

func MkdirInRoot(dir, name string, perm os.FileMode) error {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer r.Close()
	return r.Mkdir(name, perm)
}

func MkdirAllInRoot(dir, name string, perm os.FileMode) error {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer r.Close()
	return r.MkdirAll(name, perm)
}

func StatInRoot(dir, name string) (os.FileInfo, error) {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return r.Stat(name)
}

func RenameInRoot(dir string, oldName string, newName string) error {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer r.Close()
	return r.Rename(oldName, newName)
}
