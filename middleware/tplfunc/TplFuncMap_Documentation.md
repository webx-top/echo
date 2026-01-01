# TplFuncMap Documentation

## 简介 / Introduction

`TplFuncMap` 是一个全面的模板函数映射，提供了丰富的实用函数集合用于模板渲染。它包括时间操作、类型转换、字符串处理、编码/解码、数学运算等功能。

`TplFuncMap` is a comprehensive template function map that provides a rich set of utility functions for template rendering. It includes functions for time manipulation, type conversion, string processing, encoding/decoding, mathematical operations, and more.

## 快速开始 / Quick Start

```go
import (
    "html/template"
    "github.com/webx-top/echo/middleware/tplfunc"
)

// Create template with TplFuncMap / 使用 TplFuncMap 创建模板
t := template.New("example")
t.Funcs(tplfunc.New())

// Use in template / 在模板中使用
{{ Now | DateFormat "2006-01-02" }}
{{ 100 | ToFixed 2 }}
{{ "hello world" | ToUpper }}
```

## 函数分类 / Function Categories

### 1. Time Functions / 时间相关函数

| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Now` | 当前时间 / Current time | `{{ Now }}` |
| `UnixTime` | Unix时间戳 / Unix timestamp | `{{ UnixTime }}` |
| `ElapsedMemory` | 内存消耗 / Memory consumption | `{{ ElapsedMemory }}` |
| `TotalRunTime` | 运行时长(从启动服务时算起) / Runtime since service started | `{{ TotalRunTime }}` |
| `CaptchaForm` | 验证码图片表单 / CAPTCHA image form | `{{ CaptchaForm }}` |
| `FormatByte` | 字节转为适合理解的格式 / Bytes to human-readable format | `{{ 1024 \| FormatByte }}` |
| `FormatBytes` | 字节转为适合理解的格式 / Bytes to human-readable format | `{{ 1024 \| FormatBytes }}` |
| `FriendlyTime` | 友好的时间格式 / Human-friendly time format | `{{ 3600 \| FriendlyTime }}` |
| `FormatPastTime` | 以前距离现在多长时间 / Time elapsed since past | `{{ 1609459200 \| FormatPastTime }}` |
| `DateFormat` | 日期格式化 / Date format | `{{ 1609459200 \| DateFormat "2006-01-02" }}` |
| `DateFormatShort` | 短日期格式 / Short date format | `{{ 1609459200 \| DateFormatShort }}` |
| `Ts2time` | 时间戳数字转time.Time / Timestamp to time.Time | `{{ 1609459200 \| Ts2time }}` |
| `Ts2date` | 时间戳数字转日期字符串 / Timestamp to date string | `{{ 1609459200 \| Ts2date "2006-01-02" }}` |

### 2. Comparison Functions / 比较函数

| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Eq` | 等值比较 / Equal comparison | `{{ Eq a b }}` |
| `Add` | 加法 / Addition | `{{ Add 1 2 }}` |
| `Sub` | 减法 / Subtraction | `{{ Sub 5 2 }}` |
| `Div` | 除法 / Division | `{{ Div 10 2 }}` |
| `Mul` | 乘法 / Multiplication | `{{ Mul 3 4 }}` |
| `IsNil` | 检查是否为nil / Check if nil | `{{ IsNil value }}` |
| `IsEmpty` | 检查是否为空 / Check if empty | `{{ IsEmpty value }}` |
| `NotEmpty` | 检查是否不为空 / Check if not empty | `{{ NotEmpty value }}` |
| `IsNaN` | 检查是否为非数字 / Check if NaN | `{{ IsNaN value }}` |
| `IsInf` | 检查是否为无穷大 / Check if infinity | `{{ IsInf value 1 }}` |

### 3. Type Conversion Functions / 类型转换函数

| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Html`, `ToHTML` | 转换为HTML / Convert to HTML | `{{ content \| ToHTML }}` |
| `Js`, `ToJS` | 转换为JavaScript / Convert to JS | `{{ content \| ToJS }}` |
| `Css`, `ToCSS` | 转换为CSS / Convert to CSS | `{{ content \| ToCSS }}` |
| `ToURL` | 转换为URL / Convert to URL | `{{ url \| ToURL }}` |
| `ToHTMLAttr` | 转换为HTML属性 / Convert to HTML attribute | `{{ attr \| ToHTMLAttr }}` |
| `ToHTMLAttrs` | 转换为HTML属性集合 / Convert to HTML attributes | `{{ attrs \| ToHTMLAttrs }}` |
| `ToStrSlice` | 转换为字符串切片 / Convert to string slice | `{{ ToStrSlice "a" "b" }}` |
| `ToDuration` | 转换为时间间隔 / Convert to duration | `{{ ToDuration 60 "s" }}` |
| `Str` | 转换为字符串 / Convert to string | `{{ 123 \| Str }}` |
| `Int` | 转换为整数 / Convert to int | `{{ "123" \| Int }}` |
| `Int32` | 转换为32位整数 / Convert to int32 | `{{ "123" \| Int32 }}` |
| `Int64` | 转换为64位整数 / Convert to int64 | `{{ "123" \| Int64 }}` |
| `Uint` | 转换为无符号整数 / Convert to uint | `{{ "123" \| Uint }}` |
| `Uint32` | 转换为32位无符号整数 / Convert to uint32 | `{{ "123" \| Uint32 }}` |
| `Uint64` | 转换为64位无符号整数 / Convert to uint64 | `{{ "123" \| Uint64 }}` |
| `Float32` | 转换为32位浮点数 / Convert to float32 | `{{ "3.14" \| Float32 }}` |
| `Float64` | 转换为64位浮点数 / Convert to float64 | `{{ "3.14" \| Float64 }}` |
| `Float2int` | 浮点数转整数 / Float to int | `{{ 3.9 \| Float2int }}` |
| `Float2uint` | 浮点数转无符号整数 / Float to uint | `{{ 3.9 \| Float2uint }}` |
| `Float2int64` | 浮点数转64位整数 / Float to int64 | `{{ 3.9 \| Float2int64 }}` |
| `Float2uint64` | 浮点数转64位无符号整数 / Float to uint64 | `{{ 3.9 \| Float2uint64 }}` |
| `ToFloat64` | 转换为64位浮点数 / Convert to float64 | `{{ "3.14" \| ToFloat64 }}` |
| `ToFixed` | 格式化为固定位小数 / Format to fixed decimal | `{{ ToFixed 3.14159 2 }}` |
| `ToDecimal` | 转换为 decimal.Decimal 类型 / Converts any numeric type to a decimal.Decimal | `{{ "3.14" \| ToDecimal }}` |
| `NumberMore` | 数字格式化(如：99+) / Number format with max | `{{ 1000 \| NumberMore 99 }}` 大于99时显示为99+ |
| `Math` | 数学运算 / Math operations | `{{ Math "sqrt" 16 }}` |
| `NumberFormat` | 数字格式化 / Number formatting | `{{ NumberFormat 1234567.89 2 }}` |
| `NumFormat` | 数字格式化 / Number formatting | `{{ 1234567.89 \| NumFormat 2 }}` |
| `NumberTrim` | 截断数字精度 / Trim number precision | `{{ NumberTrim 3.14159 2 }}` |
| `DurationFormat` | 时间间隔格式化 / Duration formatting | `{{ 3600 \| DurationFormat "en" }}` |
| `DelimLeft` | 模板标签左分隔符 / Template tag left delimiter | `{{ DelimLeft }}` |
| `DelimRight` | 模板标签右分隔符 / Template tag right delimiter | `{{ DelimRight }}` |
| `TemplateTag` | 模板标签 / Template tag | `{{ TemplateTag "name" }}` |

### 4. String Functions / 字符串处理函数

#### 字符串检查 / String Check
| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Contains` | 检查子串 / Check substring | `{{ Contains "hello world" "world" }}` |
| `HasPrefix` | 检查前缀 / Check prefix | `{{ HasPrefix "hello" "he" }}` |
| `HasSuffix` | 检查后缀 / Check suffix | `{{ HasSuffix "hello" "lo" }}` |
| `InExt` | 检查文件扩展名 / Check file extension | `{{ InExt "image.jpg" ".jpg" ".png" }}` |

#### 字符串修剪 / String Trimming
| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Trim` | 去除空白 / Trim whitespace | `{{ "  hello  " \| Trim }}` |
| `TrimLeft` | 去除左侧 / Trim left | `{{ "  hello" \| TrimLeft }}` |
| `TrimRight` | 去除右侧 / Trim right | `{{ "hello  " \| TrimRight }}` |
| `TrimPrefix` | 去除前缀 / Trim prefix | `{{ TrimPrefix "prefix_hello" "prefix_" }}` |
| `TrimSuffix` | 去除后缀 / Trim suffix | `{{ TrimSuffix "hello_suffix" "_suffix" }}` |

#### 大小写转换 / Case Conversion
| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `ToLower` | 转小写 / To lowercase | `{{ "HELLO" \| ToLower }}` |
| `ToUpper` | 转大写 / To uppercase | `{{ "hello" \| ToUpper }}` |
| `Title` | 标题格式 / To title case | `{{ "hello world" \| Title }}` |
| `LowerCaseFirst` | 首字母小写 / Lower first letter | `{{ "Hello" \| LowerCaseFirst }}` |
| `UpperCaseFirst` | 首字母大写 / Upper first letter | `{{ "hello" \| UpperCaseFirst }}` |
| `CamelCase` | 转驼峰命名 / To camel case | `{{ "hello_world" \| CamelCase }}` |
| `PascalCase` | 转帕斯卡命名 / To pascal case | `{{ "hello_world" \| PascalCase }}` |
| `SnakeCase` | 转蛇形命名 / To snake case | `{{ "HelloWorld" \| SnakeCase }}` |

#### 路径处理 / Path Handling
| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Dir` | 目录路径 / Directory path | `{{ "/path/to/file.txt" \| Dir }}` |
| `Base` | 基础名称 / Base name | `{{ "/path/to/file.txt" \| Base }}` |
| `Ext` | 文件扩展名 / File extension | `{{ "file.txt" \| Ext }}` |
| `Dirname` | 目录名 / Directory name | `{{ "path/to/file" \| Dirname }}` |
| `Basename` | 基础名称 / Base name | `{{ "path/to/file" \| Basename }}` |
| `Extension` | 扩展名 / Extension | `{{ "file.txt" \| Extension }}` |

#### 字符串操作 / String Operations
| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Concat` | 连接字符串 / Concatenate strings | `{{ Concat "hello" "world" }}` |
| `Replace` | 替换子串 / Replace substring | `{{ Replace "hello world" "world" "go" }}` |
| `Split` | 分割字符串 / Split string | `{{ Split "a,b,c" "," }}` |
| `Join` | 连接字符串 / Join strings | `{{ Join slice "," }}` |
| `Substr` | 子字符串 / Substring | `{{ Substr "hello" 0 3 }}` |
| `StripTags` | 去除HTML标签 / Strip HTML tags | `{{ "<p>hello</p>" \| StripTags }}` |
| `Nl2br` | 换行转<br> / Newline to <br> | `{{ "line1\nline2" \| Nl2br }}` |
| `AddSuffix` | 添加后缀 / Add suffix | `{{ AddSuffix "file.txt" "_new" }}` 输出 file_new.txt |
| `RandomString` | 随机字符串 / Random string | `{{ RandomString 8 }}` |
| `Slugify` | 转URL友好格式 / To slug | `{{ "Hello World" \| Slugify }}` 输出 hello-world |
| `SlugifyMaxWidth` | 转URL友好格式(限制长度) / To slug with max width | `{{ "Hello World" \| SlugifyMaxWidth 20 }}` |
| `Reverse` | 反转字符串 / Reverse string | `{{ "hello" \| Reverse }}` 输出 olleh |

### 5. Encode & Decode Functions / 编码解码函数

| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `JSONEncode` | JSON编码 / Encode to JSON | `{{ JSONEncode data "  " }}` |
| `JSONDecode` | JSON解码 / Decode from JSON | `{{ ``{"A":1}`` \| JSONDecode }}` |
| `JSONDecodeSlice` | JSON解码为切片 / Decode JSON to slice | `{{ "[1,2,3]" \| JSONDecodeSlice }}` |
| `URLEncode` | URL编码 / URL encode | `{{ "hello world" \| URLEncode }}` |
| `URLDecode` | URL解码 / URL decode | `{{ "hello%20world" \| URLDecode }}` |
| `RawURLEncode` | 原始URL编码 / Raw URL encode | `{{ "hello world" \| RawURLEncode }}` |
| `RawURLDecode` | 原始URL解码 / Raw URL decode | `{{ url \| RawURLDecode }}` |
| `Base64Encode` | Base64编码 / Base64 encode | `{{ "hello" \| Base64Encode }}` |
| `Base64Decode` | Base64解码 / Base64 decode | `{{ base64Str \| Base64Decode }}` |
| `UnicodeDecode` | Unicode解码 / Unicode decode | `{{ "\\u4f60\\u597d" \| UnicodeDecode }}` |
| `SafeBase64Encode` | 安全Base64编码 / Safe Base64 encode | `{{ "hello" \| SafeBase64Encode }}` |
| `SafeBase64Decode` | 安全Base64解码 / Safe Base64 decode | `{{ base64Str \| SafeBase64Decode }}` |
| `Hash` | 哈希字符串 / Hash string | `{{ Hash "password" "salt" }}` |
| `Unquote` | 去除引号 / Unquote string | `{{ "'hello'" \| Unquote }}` |
| `Quote` | 添加引号 / Quote string | `{{ hello \| Quote }}` |

### 6. Map & Slice Functions / 映射和切片函数

| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `MakeMap` | 创建映射 / Create map | `{{ MakeMap "key1" "value1" "key2" "value2" }}` |
| `MakeSlice` | 创建切片 / Create slice | `{{ MakeSlice "a" "b" "c" }}` |
| `InSet` | 检查是否在集合中 / Check in set | `{{ InSet "value" set }}` |
| `InSlice` | 检查是否在切片中 / Check in slice | `{{ InSlice "value" slice }}` |
| `InSlicex` | 检查是否在切片中(interface) / Check in slice (interface) | `{{ InSlicex "value" slice }}` |
| `Set` | 设置值 / Set value | `{{ Set renderArgs "key" "value" }}` |
| `Append` | 追加到切片 / Append to slice | `{{ Append renderArgs "key" "value" }}` |
| `InStrSlice` | 检查是否在字符串切片中 / Check in string slice | `{{ InStrSlice slice "value" }}` |
| `SearchStrSlice` | 在字符串切片中搜索 / Search in string slice | `{{ SearchStrSlice slice "value" }}` |
| `URLValues` | 创建URL值 / Create URL values | `{{ URLValues "key1" "value1" }}` |
| `ToSlice` | 转换为切片 / Convert to slice | `{{ ToSlice 1 2 3 }}` |
| `StrToSlice` | 字符串转切片 / String to slice | `{{ StrToSlice "a,b,c" "," }}` |
| `GetByIndex` | 按索引获取 / Get by index | `{{ GetByIndex list 0 }}` |
| `ToParamString` | 转换为参数字符串 / Convert to param string | `{{ "value" \| ToParamString }}` |

### 7. RegExp Functions / 正则表达式函数

| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Regexp` | 编译正则表达式 / Compile regexp | `{{ $r := Regexp "^[a-z]+$" }}` |
| `RegexpPOSIX` | 编译POSIX正则表达式 / Compile POSIX regexp | `{{ $r := RegexpPOSIX "^[[:alpha:]]+$" }}` |

### 8. Other Functions / 其他函数

| 函数名 Function | 说明 Description | 示例 Example |
|---|---|---|
| `Ignore` | 忽略值 / Ignore value | `{{ value \| Ignore }}` |
| `Default` | 默认值 / Default value | `{{ value \| Default "default" }}` |
| `WithURLParams` | 带参数的URL / URL with params | `{{ WithURLParams "/path" "key" "value" }}` |
| `FullURL` | 完整URL / Full URL | `{{ FullURL "https://example.com" "/path" }}` 输出 https://example.com/path |
| `IsFullURL` | 检查是否为完整URL / Check if full URL | `{{ "https://example.com" \| IsFullURL }}` |
| `If` | 条件判断 / Conditional | `{{ If condition yes no }}` |

## 使用示例 / Usage Examples

### 时间处理 / Time Handling

```go
// 获取当前时间 / Get current time
{{ Now.Format "2006-01-02 15:04:05" }}

// 时间戳转换 / Timestamp conversion
{{ 1609459200 | DateFormat "2006-01-02" }}

// 友好时间格式 / Friendly time format
{{ 86400 | FriendlyTime }}  // 1天 / 1 day

// 时间间隔格式化 / Duration formatting
{{ 3600 | DurationFormat "en" }}  // 1 hour
```

### 数学运算 / Mathematical Operations

```go
// 基本运算 / Basic operations
{{ Add 10 20 }}        // 30
{{ Sub 50 20 }}        // 30
{{ Mul 6 7 }}         // 42
{{ Div 100 4 }}        // 25

// 数学函数 / Math functions
{{ Math "sqrt" 16 }}   // 4
{{ Math "pow" 2 10 }}  // 1024
{{ Math "abs" -5 }}    // 5

// 数字格式化 / Number formatting
{{ ToFixed 3.14159 2 }}           // "3.14"
{{ NumberFormat 1234567.89 2 }}    // "1,234,567.89"
{{ 1000 | NumberMore 99 }}           // "99+"
```

### 字符串处理 / String Processing

```go
// 大小写转换 / Case conversion
{{ "hello world" | Title }}          // "Hello World"
{{ "hello" | ToUpper }}              // "HELLO"
{{ "HELLO" | ToLower }}              // "hello"

// 命名转换 / Naming convention conversion
{{ "hello_world" | CamelCase }}      // "helloWorld"
{{ "hello_world" | PascalCase }}     // "HelloWorld"
{{ "HelloWorld" | SnakeCase }}       // "hello_world"

// 字符串操作 / String operations
{{ Contains "hello world" "world" }}   // true
{{ "file.txt" | Ext }}                   // ".txt"
{{ Concat "hello" "world" }}          // "helloworld"

// 路径处理 / Path handling
{{ "/path/to/file.txt" | Dir }}      // "/path/to"
{{ "/path/to/file.txt" | Base }}     // "file.txt"
{{ "file.txt" | Ext }}              // ".txt"
```

### 编码解码 / Encoding & Decoding

```go
// JSON / JSON
{{ JSONEncode mapData "  " }}
{{ jsonStr | JSONDecode }}

// Base64 / Base64
{{ "hello" | Base64Encode }}         // "aGVsbG8="
{{ "aGVsbG8=" | Base64Decode }}     // "hello"

// URL / URL
{{ "hello world" | URLEncode }}      // "hello%20world"
{{ "hello%20world" | URLDecode }}   // "hello world"

// Unicode / Unicode
{{ "\\u4f60\\u597d" | UnicodeDecode }}  // "你好"
```

### 类型转换 / Type Conversion

```go
// 数字转换 / Number conversion
{{ "3.14" | ToFloat64 }}    // 3.14
{{ ToFixed 3.14159 2 }}     // "3.14"
{{ "123" | Int }}           // 123

// 安全类型转换 / Safe type conversion
{{ value | Default "N/A" }}     // Return "N/A" if empty
{{ value | IsNil }}             // Check if nil
{{ value | IsEmpty }}           // Check if empty
```

### 条件判断 / Conditional Logic

```go
// 条件输出 / Conditional output
{{ If user.LoggedIn "Welcome" "Please login" }}

// 空值处理 / Empty value handling
{{ user.Name | Default "Anonymous" }}
{{ If (NotEmpty user.Email) "Email provided" "No email" }}
```

### 数据结构操作 / Data Structure Operations

```go
// 创建Map / Create map
{{ MakeMap "name" "John" "age" 30 }}

// 创建Slice / Create slice
{{ MakeSlice "apple" "banana" "orange" }}

// 检查元素 / Check element
{{ InSlice "apple" fruits }}        // true
{{ InStrSlice names "John" }}       // true
{{ SearchStrSlice names "John" }}  // 0
```

## 注意事项 / Notes

1. **线程安全 / Thread Safety**: `TplFuncMap` 是全局变量，在并发环境中使用是安全的。  
   `TplFuncMap` is a global variable and is safe to use in concurrent environments.

2. **错误处理 / Error Handling**: 部分函数在遇到错误时会记录日志并返回默认值。  
   Some functions log errors and return default values when errors occur.

3. **性能考虑 / Performance Considerations**: 频繁调用复杂函数（如 `Math`）可能影响性能，建议在控制器层处理复杂逻辑。  
   Frequent calls to complex functions (like `Math`) may impact performance. Consider handling complex logic at the controller level.

4. **HTML安全 / HTML Security**: 使用 `ToHTML`、`ToJS`、`ToCSS` 等函数时，确保内容是可信的，以避免XSS攻击。  
   When using functions like `ToHTML`, `ToJS`, `ToCSS`, ensure the content is trusted to avoid XSS attacks.

5. **日期格式 / Date Format**: Go使用特定的日期格式参考时间 "2006-01-02 15:04:05"，而不是常用的格式字符串。  
   Go uses a specific reference date "2006-01-02 15:04:05" for date formatting instead of common format strings.

## 常见问题 / FAQ

### Q: 如何自定义函数 / How to add custom functions? / A: 可以通过 `New()` 函数创建新的 `FuncMap` 并添加自定义函数：

```go
customFuncs := tplfunc.New()
customFuncs["customFunc"] = func(s string) string {
    return "custom: " + s
}
template.New("example").Funcs(customFuncs)
```

### Q: 时间格式化字符串是什么意思 / What do the date format strings mean? / A: Go使用特定的参考时间 `Mon Jan 2 15:04:05 MST 2006` 来定义格式：
- `2006` - 年 / Year
- `01` - 月 / Month
- `02` - 日 / Day
- `15` - 时(24小时制) / Hour (24-hour)
- `04` - 分 / Minute
- `05` - 秒 / Second

### Q: 如何处理空值 / How to handle empty values? / A: 使用 `Default` 函数提供默认值：

```go
{{ user.Name | Default "Guest" }}
```

## 许可证 / License

Apache License 2.0

## 更多信息 / More Information

- [Go Template Documentation](https://pkg.go.dev/text/template)
- [Echo Framework](https://github.com/webx-top/echo)
