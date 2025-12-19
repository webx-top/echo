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

package tplfunc

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/decimal"
	"github.com/gosimple/slug"

	"github.com/webx-top/captcha"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

const (
	EmptyString = ``
)

// New creates and returns a new template.FuncMap initialized with all functions
// from TplFuncMap. The returned FuncMap can be used to register template functions.
func New() (r template.FuncMap) {
	r = template.FuncMap{}
	for name, function := range TplFuncMap {
		r[name] = function
	}
	return
}

var TplFuncMap template.FuncMap = template.FuncMap{
	// ======================
	// time
	// ======================
	"Now":             Now,
	"UnixTime":        UnixTime,
	"ElapsedMemory":   com.ElapsedMemory, //内存消耗
	"TotalRunTime":    com.TotalRunTime,  //运行时长(从启动服务时算起)
	"CaptchaForm":     CaptchaForm,       //验证码图片
	"FormatByte":      com.FormatBytes,   //字节转为适合理解的格式
	"FormatBytes":     com.FormatBytes,   //字节转为适合理解的格式
	"FriendlyTime":    FriendlyTime,
	"FormatPastTime":  com.FormatPastTime, //以前距离现在多长时间
	"DateFormat":      com.DateFormat,
	"DateFormatShort": com.DateFormatShort,
	"Ts2time":         TsToTime, // 时间戳数字转time.Time
	"Ts2date":         TsToDate, // 时间戳数字转日期字符串

	// ======================
	// compare
	// ======================
	"Eq":       Eq,
	"Add":      Add,
	"Sub":      Sub,
	"Div":      Div,
	"Mul":      Mul,
	"IsNil":    IsNil,
	"IsEmpty":  IsEmpty,
	"NotEmpty": NotEmpty,
	"IsNaN":    IsNaN,
	"IsInf":    IsInf,

	// ======================
	// conversion type
	// ======================
	"Html":           ToHTML,
	"Js":             ToJS,
	"Css":            ToCSS,
	"ToJS":           ToJS,
	"ToCSS":          ToCSS,
	"ToURL":          ToURL,
	"ToHTML":         ToHTML,
	"ToHTMLAttr":     ToHTMLAttr,
	"ToHTMLAttrs":    ToHTMLAttrs,
	"ToStrSlice":     ToStrSlice,
	"ToDuration":     ToDuration,
	"Str":            com.Str,
	"Int":            com.Int,
	"Int32":          com.Int32,
	"Int64":          com.Int64,
	"Uint":           com.Uint,
	"Uint32":         com.Uint32,
	"Uint64":         com.Uint64,
	"Float32":        com.Float32,
	"Float64":        com.Float64,
	"Float2int":      com.Float2int,
	"Float2uint":     com.Float2uint,
	"Float2int64":    com.Float2int64,
	"Float2uint64":   com.Float2uint64,
	"ToFloat64":      ToFloat64,
	"ToFixed":        ToFixed,
	"ToDecimal":      ToDecimal,
	"NumberMore":     NumberMore, // {{ 1000 | NumberMore 99 }} 99+
	"Math":           Math,
	"NumberFormat":   NumberFormat,
	"NumFormat":      com.NumFormat,
	"NumberTrim":     NumberTrim,
	"DurationFormat": DurationFormat,
	"DelimLeft":      DelimLeft,
	"DelimRight":     DelimRight,
	"TemplateTag":    TemplateTag,

	// ======================
	// string
	// ======================
	"Contains":   strings.Contains,
	"HasPrefix":  strings.HasPrefix,
	"HasSuffix":  strings.HasSuffix,
	"Trim":       Trim,
	"TrimLeft":   strings.TrimLeft,
	"TrimRight":  strings.TrimRight,
	"TrimPrefix": strings.TrimPrefix,
	"TrimSuffix": strings.TrimSuffix,

	"ToLower":        strings.ToLower,
	"ToUpper":        strings.ToUpper,
	"Title":          com.Title,
	"LowerCaseFirst": com.LowerCaseFirst,
	"UpperCaseFirst": com.UpperCaseFirst,
	"CamelCase":      com.CamelCase,
	"PascalCase":     com.PascalCase,
	"SnakeCase":      com.SnakeCase,
	"Reverse":        com.Reverse,
	"Dir":            filepath.Dir,
	"Base":           filepath.Base,
	"Ext":            filepath.Ext,
	"Dirname":        path.Dir,
	"Basename":       path.Base,
	"Extension":      path.Ext,
	"InExt":          InExt,

	"Concat":          Concat,
	"Replace":         strings.Replace, //strings.Replace(s, old, new, n)
	"Split":           strings.Split,
	"Join":            strings.Join,
	"Substr":          com.Substr,
	"StripTags":       com.StripTags,
	"Nl2br":           NlToBr, // \n替换为<br>
	"AddSuffix":       AddSuffix,
	"RandomString":    RandomString,
	"Slugify":         Slugify,
	"SlugifyMaxWidth": SlugifyMaxWidth,

	// ======================
	// encode & decode
	// ======================
	"JSONEncode":       JSONEncode,
	"JSONDecode":       JSONDecode,
	"JSONDecodeSlice":  JSONDecodeSlice,
	"URLEncode":        com.URLEncode,
	"URLDecode":        URLDecode,
	"RawURLEncode":     com.RawURLEncode,
	"RawURLDecode":     URLDecode,
	"Base64Encode":     com.Base64Encode,
	"Base64Decode":     Base64Decode,
	"UnicodeDecode":    UnicodeDecode,
	"SafeBase64Encode": com.SafeBase64Encode,
	"SafeBase64Decode": SafeBase64Decode,
	"Hash":             Hash,
	"Unquote":          Unquote,
	"Quote":            strconv.Quote,

	// ======================
	// map & slice
	// ======================
	"MakeMap":        MakeMap,
	"MakeSlice":      MakeSlice,
	"InSet":          com.InSet,
	"InSlice":        com.InSlice,
	"InSlicex":       com.InSliceIface,
	"Set":            Set,
	"Append":         Append,
	"InStrSlice":     InStrSlice,
	"SearchStrSlice": SearchStrSlice,
	"URLValues":      URLValues,
	"ToSlice":        ToSlice,
	"StrToSlice":     StrToSlice,
	"GetByIndex":     param.GetByIndex,
	"ToParamString":  func(v string) param.String { return param.String(v) },

	// ======================
	// regexp
	// ======================
	"Regexp":      regexp.MustCompile,
	"RegexpPOSIX": regexp.MustCompilePOSIX,

	// ======================
	// other
	// ======================
	"Ignore":        Ignore,
	"Default":       Default,
	"WithURLParams": com.WithURLParams,
	"FullURL":       com.FullURL,
	"IsFullURL":     com.IsFullURL,
	"If":            If,
}

var (
	HashSalt          = time.Now().Format(time.RFC3339)
	HashClipPositions = []uint{1, 3, 8, 9}
	NumberFormat      = com.NumberFormat
)

// RandomString generates a random alphanumeric string of specified length.
// If no length is provided, it defaults to 8 characters.
func RandomString(length ...uint) string {
	if len(length) > 0 && length[0] > 0 {
		return com.RandomAlphanumeric(length[0])
	}
	return com.RandomAlphanumeric(8)
}

// Slugify 将字符串转换为URL友好的slug格式
func Slugify(v string) string {
	return slug.Make(v)
}

// SlugifyMaxWidth 将字符串转换为slug格式并限制最大长度
// v: 需要转换的原始字符串
// maxWidth: 返回字符串的最大长度限制
// 返回: 转换后的slug字符串，长度不超过maxWidth
func SlugifyMaxWidth(v string, maxWidth int) string {
	return com.Substr(slug.Make(v), ``, maxWidth)
}

// Hash generates a hashed string from the given text using the provided salt and positions.
// If salt is empty, it uses the default HashSalt. If no positions are provided, it uses the default HashClipPositions.
// The resulting hash is created using com.MakePassword with the specified parameters.
func Hash(text, salt string, positions ...uint) string {
	if len(salt) < 1 {
		salt = HashSalt
	}
	if len(positions) < 1 {
		positions = HashClipPositions
	}
	return com.MakePassword(text, salt, positions...)
}

// Trim removes leading and trailing whitespace from the string s if no cutset is provided.
// If cutset is provided, it removes all leading and trailing characters contained in the first string of cutset.
// The function returns the trimmed string.
func Trim(s string, cutset ...string) string {
	if len(cutset) == 0 || len(cutset[0]) == 0 {
		return strings.TrimSpace(s)
	}
	return strings.Trim(s, cutset[0])
}

// Unquote removes the surrounding quotes from a string if present.
// It handles the string as if it were quoted with HTML entity &quot;.
func Unquote(s string) string {
	r, _ := strconv.Unquote(`"` + s + `"`)
	return r
}

// NumberMore compares the input number n with max value and returns max+"+" if n is greater than max,
// otherwise returns n. Supports uint, uint32, uint64, int, int32, int64, float32 and float64 types.
func NumberMore(max interface{}, n interface{}) interface{} {
	var more bool
	switch v := n.(type) {
	case uint:
		more = v > com.Uint(max)
	case uint32:
		more = v > com.Uint32(max)
	case uint64:
		more = v > com.Uint64(max)
	case int:
		more = v > com.Int(max)
	case int32:
		more = v > com.Int32(max)
	case int64:
		more = v > com.Int64(max)
	case float64:
		more = v > com.Float64(max)
	case float32:
		more = v > com.Float32(max)
	}
	if more {
		return com.String(max) + `+`
	}
	return n
}

// UnicodeDecode converts Unicode escape sequences (like \uXXXX) in the input string to their corresponding Unicode characters.
// It processes the string sequentially, handling both escaped Unicode sequences and regular characters.
// Invalid escape sequences are preserved as-is in the output.
func UnicodeDecode(str string) string {
	buf := bytes.NewBuffer(nil)
	i, j := 0, len(str)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(str[i:])
			break
		}
		if str[i] == '\\' && str[i+1] == 'u' {
			hex := str[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(str[i:x])
			}
			i = x
		} else {
			buf.WriteByte(str[i])
			i++
		}
	}
	return buf.String()
}

// JSONEncode encodes the given value to a JSON string with optional indentation.
// It returns the JSON string representation of the value.
// The indents parameter specifies the indentation string to use (e.g., "  " for two spaces).
func JSONEncode(s interface{}, indents ...string) string {
	r, _ := com.JSONEncode(s, indents...)
	return string(r)
}

// JSONDecode decodes a JSON string into a map[string]interface{}.
// If decoding fails, it logs the error and returns an empty map.
func JSONDecode(s string) map[string]interface{} {
	r := map[string]interface{}{}
	e := com.JSONDecode([]byte(s), &r)
	if e != nil {
		log.Println(e)
	}
	return r
}

// JSONDecodeSlice decodes a JSON string into a slice of interfaces.
// If decoding fails, logs the error and returns an empty slice.
func JSONDecodeSlice(s string) []interface{} {
	r := []interface{}{}
	e := com.JSONDecode([]byte(s), &r)
	if e != nil {
		log.Println(e)
	}
	return r
}

// URLDecode decodes a URL-encoded string and returns the decoded result.
// If decoding fails, logs the error and returns the original string.
func URLDecode(s string) string {
	r, e := com.URLDecode(s)
	if e != nil {
		log.Println(e)
	}
	return r
}

// Base64Decode decodes a base64 encoded string and returns the result.
// If decoding fails, it logs the error and returns an empty string.
func Base64Decode(s string) string {
	r, e := com.Base64Decode(s)
	if e != nil {
		log.Println(e)
	}
	return r
}

// SafeBase64Decode decodes a base64 encoded string safely, returning the decoded string.
// If decoding fails, it logs the error and returns an empty string.
func SafeBase64Decode(s string) string {
	r, e := com.SafeBase64Decode(s)
	if e != nil {
		log.Println(e)
	}
	return r
}

// Ignore returns nil for any input value, effectively ignoring it.
func Ignore(_ interface{}) interface{} {
	return nil
}

// URLValues creates new url.Values and adds the provided values to it.
// It accepts variadic arguments of key-value pairs or maps to populate the values.
// Returns the populated url.Values.
func URLValues(values ...interface{}) url.Values {
	v := url.Values{}
	return AddURLValues(v, values...)
}

// AddURLValues adds key-value pairs from the values slice to the url.Values.
// The values slice should contain alternating keys and values (key1, value1, key2, value2, ...).
// If an odd number of arguments is provided, the last key will be added with an empty value.
// Returns the modified url.Values.
func AddURLValues(v url.Values, values ...interface{}) url.Values {
	var k string
	for i, j := 0, len(values); i < j; i++ {
		if i%2 == 0 {
			k = fmt.Sprint(values[i])
			continue
		}
		v.Add(k, fmt.Sprint(values[i]))
		k = ``
	}
	if len(k) > 0 {
		v.Add(k, ``)
		k = ``
	}
	return v
}

// ToStrSlice converts variadic string arguments into a string slice.
func ToStrSlice(s ...string) []string {
	return s
}

// ToSlice converts variadic arguments into a slice of interfaces.
func ToSlice(s ...interface{}) []interface{} {
	return s
}

// StrToSlice converts a string into a slice of interfaces by splitting it with the specified separator.
// Each substring becomes an element in the returned slice.
func StrToSlice(s string, sep string) []interface{} {
	ss := strings.Split(s, sep)
	r := make([]interface{}, len(ss))
	for i, s := range ss {
		r[i] = s
	}
	return r
}

// Concat joins multiple strings together without any separator.
// It takes a variadic number of string arguments and returns their concatenation.
func Concat(s ...string) string {
	return strings.Join(s, ``)
}

// If returns yesValue if condition is true, otherwise returns noValue.
func If(condition bool, yesValue interface{}, noValue interface{}) interface{} {
	if condition {
		return yesValue
	}
	return noValue
}

// InExt checks if the file extension of fileName matches any of the provided extensions (case-insensitive).
// fileName: The filename to check
// exts: List of extensions to match against (e.g. ".jpg", ".png")
// Returns true if the file extension matches any of the provided extensions, false otherwise
func InExt(fileName string, exts ...string) bool {
	var max int
	for _, _ext := range exts {
		l := len(_ext)
		if max < l {
			max = l
		}
	}
	var ext string
	for i, j := len(fileName)-1, 0; i >= 0; i-- {
		j++
		if fileName[i] == '.' {
			ext = fileName[i:]
			break
		}
		if j >= max {
			return false
		}
	}
	for _, _ext := range exts {
		if strings.EqualFold(ext, _ext) {
			return true
		}
	}
	return false
}

// Default returns defaultV if v is nil, empty, zero or converts to an empty string.
// Otherwise, it returns v. It handles various primitive types including strings,
// numeric types (int, float), and converts other types to string for empty check.
func Default(defaultV interface{}, v interface{}) interface{} {
	switch val := v.(type) {
	case nil:
		return defaultV
	case string:
		if len(val) == 0 {
			return defaultV
		}
	case uint8, int8, uint, int, uint32, int32, int64, uint64:
		if val == 0 {
			return defaultV
		}
	case float32, float64:
		if val == 0.0 {
			return defaultV
		}
	default:
		if len(com.Str(v)) == 0 {
			return defaultV
		}
	}
	return v
}

// Set adds or updates a key-value pair in the renderArgs map and returns an empty string.
func Set(renderArgs map[string]interface{}, key string, value interface{}) string {
	renderArgs[key] = value
	return EmptyString
}

// Append adds a value to a slice in the renderArgs map under the specified key.
// If the key doesn't exist, it creates a new slice with the value.
// Returns EmptyString as a placeholder (no meaningful return value).
func Append(renderArgs map[string]interface{}, key string, value interface{}) string {
	if renderArgs[key] == nil {
		renderArgs[key] = []interface{}{value}
	} else {
		renderArgs[key] = append(renderArgs[key].([]interface{}), value)
	}
	return EmptyString
}

// NlToBr Replaces newlines with <br />
func NlToBr(text string) template.HTML {
	return template.HTML(Nl2br(text))
}

// CaptchaForm 验证码表单域
func CaptchaForm(args ...interface{}) template.HTML {
	return CaptchaFormWithURLPrefix(``, args...)
}

// CaptchaFormWithURLPrefix 验证码表单域
func CaptchaFormWithURLPrefix(urlPrefix string, args ...interface{}) template.HTML {
	id := "captcha"
	msg := "页面验证码已经失效，必须重新请求当前页面。确定要刷新本页面吗？"
	onErr := "if(this.src.indexOf('?reload=')!=-1 && confirm('%s')) window.location.reload();"
	format := `<img id="%[2]sImage" src="` + urlPrefix + `/captcha/%[1]s.png" alt="Captcha image" onclick="this.src=this.src.split('?')[0]+'?reload='+Math.random();" onerror="%[3]s" style="cursor:pointer" /><input type="hidden" name="captchaId" id="%[2]sId" value="%[1]s" />`
	var (
		customOnErr bool
		cid         string
	)
	switch len(args) {
	case 3:
		switch v := args[2].(type) {
		case template.JS:
			onErr = string(v)
			customOnErr = true
		case string:
			msg = v
		}
		fallthrough
	case 2:
		if args[1] != nil {
			v := fmt.Sprint(args[1])
			format = v
		}
		fallthrough
	case 1:
		switch v := args[0].(type) {
		case template.JS:
			onErr = string(v)
			customOnErr = true
		case template.HTML:
			format = string(v)
		case string:
			id = v
		case param.Store:
			cid = v.String(`captchaId`)
			if v.Has(`onErr`) {
				onErr = v.String(`onErr`)
			}
			if v.Has(`format`) {
				format = v.String(`format`)
			}
			if v.Has(`id`) {
				id = v.String(`id`)
			}
		case map[string]interface{}:
			h := param.Store(v)
			cid = h.String(`captchaId`)
			if h.Has(`onErr`) {
				onErr = h.String(`onErr`)
			}
			if h.Has(`format`) {
				format = h.String(`format`)
			}
			if h.Has(`id`) {
				id = h.String(`id`)
			}
		}
	}
	if len(cid) == 0 {
		cid = captcha.New()
	}
	if !customOnErr {
		onErr = fmt.Sprintf(onErr, msg)
	}
	return template.HTML(fmt.Sprintf(format, cid, id, onErr))
}

// CaptchaVerify 验证码验证
func CaptchaVerify(captchaSolution string, idGet func(string, ...string) string) bool {
	//id := r.FormValue("captchaId")
	id := idGet("captchaId")
	return captcha.VerifyString(id, captchaSolution)
}

// Nl2br 将换行符替换为<br />
func Nl2br(text string) string {
	return com.Nl2br(template.HTMLEscapeString(text))
}

// IsNil checks if the given interface value is nil. Returns true if the value is nil, false otherwise.
func IsNil(a interface{}) bool {
	switch a.(type) {
	case nil:
		return true
	default:
		//return reflect.ValueOf(a).IsNil()
	}
	return false
}

func interface2Int64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}

func interface2Float64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// ToFloat64 converts the given value to float64. It first attempts to convert to int64,
// then to float64, and finally falls back to a general conversion if the previous attempts fail.
func ToFloat64(value interface{}) float64 {
	if v, ok := interface2Int64(value); ok {
		return float64(v)
	}
	if v, ok := interface2Float64(value); ok {
		return v
	}
	return com.Float64(value)
}

// Add returns the sum of two numeric values (int64 or float64) after converting them to compatible types.
// It handles mixed type inputs by converting both values to float64 if either is not an integer.
// left: first value to add (int64 or float64)
// right: second value to add (int64 or float64)
// Returns: sum of left and right as int64 if both are integers, otherwise float64
func Add(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool
	rleft, isInt = interface2Int64(left)
	if !isInt {
		fleft, _ = interface2Float64(left)
	}
	rright, isInt = interface2Int64(right)
	if !isInt {
		fright, _ = interface2Float64(right)
	}
	intSum := rleft + rright

	if isInt {
		return intSum
	}
	return fleft + fright + float64(intSum)
}

// Div returns the division result of left divided by right after converting both to float64.
func Div(left interface{}, right interface{}) interface{} {
	return ToFloat64(left) / ToFloat64(right)
}

// Mul returns the product of two values after converting them to float64.
// left: the first value to multiply
// right: the second value to multiply
func Mul(left interface{}, right interface{}) interface{} {
	return ToFloat64(left) * ToFloat64(right)
}

// Math performs various mathematical operations based on the given operation string.
// Supported operations: mod, abs, acos, acosh, asin, asinh, atan, atan2, atanh, cbrt,
// ceil, copysign, cos, cosh, dim, erf, erfc, exp, exp2, floor, max, min, pow, sqrt,
// sin, log, log2, log10, tan, tanh, add, sub, mul, div.
// Returns the result as interface{} which can be converted to appropriate numeric type.
// Returns 0 if insufficient arguments are provided for the operation.
func Math(op string, args ...interface{}) interface{} {
	length := len(args)
	if length < 1 {
		return float64(0)
	}
	switch op {
	case `mod`: //模
		if length < 2 {
			return float64(0)
		}
		return math.Mod(ToFloat64(args[0]), ToFloat64(args[1]))
	case `abs`:
		return math.Abs(ToFloat64(args[0]))
	case `acos`:
		return math.Acos(ToFloat64(args[0]))
	case `acosh`:
		return math.Acosh(ToFloat64(args[0]))
	case `asin`:
		return math.Asin(ToFloat64(args[0]))
	case `asinh`:
		return math.Asinh(ToFloat64(args[0]))
	case `atan`:
		return math.Atan(ToFloat64(args[0]))
	case `atan2`:
		if length < 2 {
			return float64(0)
		}
		return math.Atan2(ToFloat64(args[0]), ToFloat64(args[1]))
	case `atanh`:
		return math.Atanh(ToFloat64(args[0]))
	case `cbrt`:
		return math.Cbrt(ToFloat64(args[0]))
	case `ceil`:
		return math.Ceil(ToFloat64(args[0]))
	case `copysign`:
		if length < 2 {
			return float64(0)
		}
		return math.Copysign(ToFloat64(args[0]), ToFloat64(args[1]))
	case `cos`:
		return math.Cos(ToFloat64(args[0]))
	case `cosh`:
		return math.Cosh(ToFloat64(args[0]))
	case `dim`:
		if length < 2 {
			return float64(0)
		}
		return math.Dim(ToFloat64(args[0]), ToFloat64(args[1]))
	case `erf`:
		return math.Erf(ToFloat64(args[0]))
	case `erfc`:
		return math.Erfc(ToFloat64(args[0]))
	case `exp`:
		return math.Exp(ToFloat64(args[0]))
	case `exp2`:
		return math.Exp2(ToFloat64(args[0]))
	case `floor`:
		return math.Floor(ToFloat64(args[0]))
	case `max`:
		if length < 2 {
			return float64(0)
		}
		return math.Max(ToFloat64(args[0]), ToFloat64(args[1]))
	case `min`:
		if length < 2 {
			return float64(0)
		}
		return math.Min(ToFloat64(args[0]), ToFloat64(args[1]))
	case `pow`: //幂
		if length < 2 {
			return float64(0)
		}
		return math.Pow(ToFloat64(args[0]), ToFloat64(args[1]))
	case `sqrt`: //平方根
		return math.Sqrt(ToFloat64(args[0]))
	case `sin`:
		return math.Sin(ToFloat64(args[0]))
	case `log`:
		return math.Log(ToFloat64(args[0]))
	case `log2`:
		return math.Log2(ToFloat64(args[0]))
	case `log10`:
		return math.Log10(ToFloat64(args[0]))
	case `tan`:
		return math.Tan(ToFloat64(args[0]))
	case `tanh`:
		return math.Tanh(ToFloat64(args[0]))
	case `add`: //加
		if length < 2 {
			return float64(0)
		}
		return Add(ToFloat64(args[0]), ToFloat64(args[1]))
	case `sub`: //减
		if length < 2 {
			return float64(0)
		}
		return Sub(ToFloat64(args[0]), ToFloat64(args[1]))
	case `mul`: //乘
		if length < 2 {
			return float64(0)
		}
		return Mul(ToFloat64(args[0]), ToFloat64(args[1]))
	case `div`: //除
		if length < 2 {
			return float64(0)
		}
		return Div(ToFloat64(args[0]), ToFloat64(args[1]))
	}
	return nil
}

// IsNaN reports whether v is a NaN (Not a Number) value after converting it to float64.
func IsNaN(v interface{}) bool {
	return math.IsNaN(ToFloat64(v))
}

// IsInf reports whether v is an infinity according to s.
// v is converted to float64 and s is converted to int before comparison.
// Returns true if v is positive or negative infinity.
func IsInf(v interface{}, s interface{}) bool {
	return math.IsInf(ToFloat64(v), com.Int(s))
}

// Sub returns the result of subtracting right from left. It supports both integer
// and floating-point numbers by automatically converting the inputs to the
// appropriate numeric type. If either operand is a float, the result will be
// a float; otherwise, it returns an integer result.
func Sub(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool
	rleft, isInt = interface2Int64(left)
	if !isInt {
		fleft, _ = interface2Float64(left)
	}
	rright, isInt = interface2Int64(right)
	if !isInt {
		fright, _ = interface2Float64(right)
	}
	if isInt {
		return rleft - rright
	}
	return fleft + float64(rleft) - (fright + float64(rright))
}

// ToFixed converts a value to a fixed-point string representation with specified precision.
// value: the input value to convert (can be any numeric type or string representation of a number)
// precision: the number of decimal places to round to (must be convertible to int)
func ToFixed(value interface{}, precision interface{}) string {
	return fmt.Sprintf("%.*f", com.Int(precision), ToFloat64(value))
}

// Now returns the current local time.
func Now() time.Time {
	return time.Now()
}

// UnixTime returns the current time as Unix timestamp (seconds since epoch)
func UnixTime() int64 {
	return time.Now().Unix()
}

// Eq compares two values for equality, handling nil cases properly.
// Returns true if both values are nil or their string representations are equal.
func Eq(left interface{}, right interface{}) bool {
	leftIsNil := (left == nil)
	rightIsNil := (right == nil)
	if leftIsNil || rightIsNil {
		if leftIsNil && rightIsNil {
			return true
		}
		return false
	}
	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

// ToHTML converts raw value to template.HTML type.
// If input is already template.HTML, returns it directly.
// For strings, converts to template.HTML without escaping.
// For other types, converts to string first using com.String.
func ToHTML(raw interface{}) template.HTML {
	switch v := raw.(type) {
	case template.HTML:
		return v
	case string:
		return template.HTML(v)
	default:
		return template.HTML(com.String(raw))
	}
}

// ToHTMLAttr converts various input types to template.HTMLAttr.
// It accepts template.HTML, template.HTMLAttr, string, or any type that can be converted to string.
// Returns the input as template.HTMLAttr, converting non-string types using com.String.
func ToHTMLAttr(raw interface{}) template.HTMLAttr {
	switch v := raw.(type) {
	case template.HTML:
		return template.HTMLAttr(string(v))
	case template.HTMLAttr:
		return v
	case string:
		return template.HTMLAttr(v)
	default:
		return template.HTMLAttr(com.String(raw))
	}
}

// ToHTMLAttrs converts a map of string keys to interface values into a map of HTMLAttr keys.
// The keys are converted using ToHTMLAttr function while preserving the original values.
func ToHTMLAttrs(raw map[string]interface{}) (r map[template.HTMLAttr]interface{}) {
	r = make(map[template.HTMLAttr]interface{})
	for k, v := range raw {
		r[ToHTMLAttr(k)] = v
	}
	return
}

// ToJS converts various input types to template.JS type for safe JavaScript embedding.
// It handles template.HTML, template.JS, string, and other types (converted via com.String).
// The conversion ensures the output is properly escaped for JavaScript contexts.
func ToJS(raw interface{}) template.JS {
	switch v := raw.(type) {
	case template.HTML:
		return template.JS(string(v))
	case template.JS:
		return v
	case string:
		return template.JS(v)
	default:
		return template.JS(com.String(raw))
	}
}

// ToCSS converts various input types to template.CSS type. It handles conversion from
// template.HTML, template.CSS, string, and other types (using com.String for conversion).
// The function ensures the output is always of type template.CSS for safe HTML/CSS rendering.
func ToCSS(raw interface{}) template.CSS {
	switch v := raw.(type) {
	case template.HTML:
		return template.CSS(string(v))
	case template.CSS:
		return v
	case string:
		return template.CSS(v)
	default:
		return template.CSS(com.String(raw))
	}
}

// ToURL converts various input types to template.URL type.
// It handles conversion from template.HTML, template.URL, string, and other types (using com.String).
// Returns the converted template.URL value.
func ToURL(raw interface{}) template.URL {
	switch v := raw.(type) {
	case template.HTML:
		return template.URL(string(v))
	case template.URL:
		return v
	case string:
		return template.URL(v)
	default:
		return template.URL(com.String(raw))
	}
}

// AddSuffix adds the given suffix to the string before the last occurrence of the specified character.
// If no character is specified, it defaults to '.'.
// If the character is empty or not found in the string, it simply appends the suffix to the end.
// Additional args can be provided to specify the character before which to add the suffix.
func AddSuffix(s string, suffix string, args ...string) string {
	beforeChar := `.`
	if len(args) > 0 {
		beforeChar = args[0]
		if beforeChar == `` {
			return s + suffix
		}
	}
	p := strings.LastIndex(s, beforeChar)
	if p < 0 {
		return s
	}
	return s[0:p] + suffix + s[p:]
}

// IsEmpty checks if the given interface value is empty.
// It returns true for nil, empty string, empty slice, or when the string representation is "<nil>", "", or "[]".
func IsEmpty(a interface{}) bool {
	switch v := a.(type) {
	case nil:
		return true
	case string:
		return len(v) == 0
	case []interface{}:
		return len(v) < 1
	default:
		switch fmt.Sprintf(`%v`, a) {
		case `<nil>`, ``, `[]`:
			return true
		}
	}
	return false
}

// NotEmpty reports whether the given value is not empty.
// It returns the inverse of IsEmpty(a).
func NotEmpty(a interface{}) bool {
	return !IsEmpty(a)
}

// InStrSlice checks if a string value exists in a string slice.
// Returns true if the value is found, false otherwise.
func InStrSlice(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// SearchStrSlice searches for the given value in a string slice and returns its index.
// Returns -1 if the value is not found.
func SearchStrSlice(values []string, value string) int {
	for i, v := range values {
		if v == value {
			return i
		}
	}
	return -1
}

// DurationFormat converts the given time duration to a formatted string representation based on the specified language.
// The 't' parameter can be a time.Duration, string, or numeric value representing duration.
// The 'lang' parameter specifies the language for formatting (e.g., "en" for English).
// Optional 'args' can provide additional formatting parameters.
// Returns a pointer to com.Durafmt containing the formatted duration.
func DurationFormat(lang interface{}, t interface{}, args ...string) *com.Durafmt {
	duration := ToDuration(t, args...)
	return com.ParseDuration(duration, lang)
}

// ToTime converts various input types to time.Time.
// Supports time.Time, string (format: "2006-01-02 15:04:05"), and other types via TsToTime.
// Panics if string parsing fails.
func ToTime(t interface{}) time.Time {
	switch v := t.(type) {
	case time.Time:
		return v
	case string:
		t, err := time.ParseInLocation(`2006-01-02 15:04:05`, v, time.Local)
		if err != nil {
			panic(err)
		}
		return t
	default:
		return TsToTime(t)
	}
}

// ToDuration converts various input types to time.Duration with optional unit specification.
// Accepts numeric types (int, int64, uint, etc.) and time.Duration as input.
// Optional args[0] specifies the unit: "ns", "us", "ms", "s", "m", "h" (default: "s").
// Returns the converted duration value.
func ToDuration(t interface{}, args ...string) time.Duration {
	td := time.Second
	if len(args) > 0 {
		switch args[0] {
		case `ns`:
			td = time.Nanosecond
		case `us`:
			td = time.Microsecond
		case `s`:
			td = time.Second
		case `ms`:
			td = time.Millisecond
		case `h`:
			td = time.Hour
		case `m`:
			td = time.Minute
		}
	}
	switch v := t.(type) {
	case time.Duration:
		return v
	case int64:
		td = time.Duration(v) * td
	case int:
		td = time.Duration(v) * td
	case uint:
		td = time.Duration(v) * td
	case int32:
		td = time.Duration(v) * td
	case uint32:
		td = time.Duration(v) * td
	case uint64:
		td = time.Duration(v) * td
	default:
		td = time.Duration(com.Int64(t)) * td
	}
	return td
}

// FriendlyTime converts various time representations to a human-friendly duration string.
// It accepts time.Duration, int, int64, uint, int32, uint32, uint64, or any type that can be converted to int64.
// The optional args parameter allows for customizing the output format.
// Returns a formatted string representation of the duration.
func FriendlyTime(t interface{}, args ...interface{}) string {
	var td time.Duration
	switch v := t.(type) {
	case time.Duration:
		td = v
	case int64:
		td = time.Duration(v)
	case int:
		td = time.Duration(v)
	case uint:
		td = time.Duration(v)
	case int32:
		td = time.Duration(v)
	case uint32:
		td = time.Duration(v)
	case uint64:
		td = time.Duration(v)
	default:
		td = time.Duration(com.Int64(t))
	}
	return com.FriendlyTime(td, args...)
}

// TsToTime converts a timestamp of various types to time.Time by delegating to TimestampToTime.
func TsToTime(timestamp interface{}) time.Time {
	return TimestampToTime(timestamp)
}

// TsToDate converts a timestamp to a formatted date string.
// The format parameter follows the standard Go time format layout.
// If the timestamp is invalid or zero, returns an empty string.
func TsToDate(format string, timestamp interface{}) string {
	t := TimestampToTime(timestamp)
	if t.IsZero() {
		return EmptyString
	}
	return t.Format(format)
}

// TimestampToTime converts various timestamp formats to time.Time.
// It accepts int, uint, int64, uint64, int32, uint32 or string representations of timestamps.
// For string inputs, it attempts to parse them as base-10 integers.
// Returns the corresponding time.Time value or zero time if parsing fails.
func TimestampToTime(timestamp interface{}) time.Time {
	var ts int64
	switch v := timestamp.(type) {
	case int64:
		ts = v
	case uint:
		ts = int64(v)
	case int:
		ts = int64(v)
	case uint32:
		ts = int64(v)
	case int32:
		ts = int64(v)
	case uint64:
		ts = int64(v)
	default:
		i, e := strconv.ParseInt(fmt.Sprint(timestamp), 10, 64)
		if e != nil {
			log.Println(e)
		}
		ts = i
	}
	return time.Unix(ts, 0)
}

// ToDecimal converts any numeric type to a decimal.Decimal.
// It first converts the input to float64 using ToFloat64, then creates a decimal from the float value.
func ToDecimal(number interface{}) decimal.Decimal {
	money := ToFloat64(number)
	return decimal.NewFromFloat(money)
}

// NumberTrim converts a number to float64, truncates it to the specified precision,
// and formats it with optional separator. Returns the formatted string representation.
//
// Parameters:
//   - number: the input number to be formatted (can be any numeric type)
//   - precision: the number of decimal places to keep
//   - separator: optional thousand separator (default is none)
//
// Returns: formatted string representation of the number
func NumberTrim(number interface{}, precision int, separator ...string) string {
	money := ToFloat64(number)
	s := decimal.NewFromFloat(money).Truncate(int32(precision)).String()
	return com.NumberTrim(s, precision, separator...)
}

// MakeMap creates a param.Store from alternating key-value pairs.
// It accepts either a flat list of arguments or a single slice of values.
// Keys are converted to strings using fmt.Sprint. If an odd number of arguments
// is provided, the last key will be set with a nil value.
func MakeMap(values ...interface{}) param.Store {
	h := param.Store{}
	length := len(values)
	if length == 0 {
		return h
	}
	if length == 1 {
		if vals, ok := values[0].([]interface{}); ok {
			length = len(vals)
			if length == 0 {
				return h
			}
			values = vals
		}
	}
	var k string
	for i, j := 0, length; i < j; i++ {
		if i%2 == 0 {
			k = fmt.Sprint(values[i])
			continue
		}
		h.Set(k, values[i])
		k = ``
	}
	if len(k) > 0 {
		h.Set(k, nil)
	}
	return h
}

type iSlice []interface{}

// Add appends the given elements to the slice and returns an empty string.
func (i *iSlice) Add(sl ...interface{}) string {
	*i = append(*i, sl...)
	return EmptyString
}

// MakeSlice converts variadic arguments into an iSlice type.
func MakeSlice(values ...interface{}) iSlice {
	return iSlice(values)
}

// DelimLeft returns the left delimiter used in templates.
func DelimLeft() string {
	return `{{`
}

// DelimRight returns the right delimiter used in templates.
func DelimRight() string {
	return `}}`
}

// TemplateTag returns a template tag string by combining the given name with delimiters.
func TemplateTag(name string) string {
	return DelimLeft() + name + DelimRight()
}
