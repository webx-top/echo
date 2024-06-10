package middleware

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/webx-top/echo"
)

type (
	RewriteRegExp struct {
		Old   *regexp.Regexp
		New   *regexp.Regexp
		Named map[string]string
	}
	// RewriteConfig defines the config for Rewrite middleware.
	RewriteConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		// Rules defines the URL path rewrite rules. The values captured in asterisk can be
		// retrieved by index e.g. $1, $2 and so on.
		// Example:
		// "/old":              "/new",
		// "/api/*":            "/$1",
		// "/js/*":             "/public/javascripts/$1",
		// "/users/*/orders/*": "/user/$1/order/$2",
		// "/users/:id": "/user/$1",
		// "/match/<name:[0-9]+>": "/user/$1",
		// Required.
		Rules             map[string]string `json:"rules"`
		DisableColonParam bool

		addresses  map[string]*RewriteRegExp //"/old": old<*regexp.Regexp>
		rulesRegex map[*regexp.Regexp]string //old<*regexp.Regexp>: "/new"
		rvsesRegex map[*regexp.Regexp]string //new<*regexp.Regexp>: "/old/$1"
	}
)

// Init Initialize
func (c *RewriteConfig) Init() *RewriteConfig {
	c.rulesRegex = map[*regexp.Regexp]string{}
	c.rvsesRegex = map[*regexp.Regexp]string{}
	c.addresses = map[string]*RewriteRegExp{}

	for k, v := range c.Rules {
		c.Set(k, v)
	}
	return c
}

var regexpQuotedBytes = map[rune]struct{}{
	'\\': {},
	'.':  {},
	'+':  {},
	'*':  {},
	'?':  {},
	'(':  {},
	')':  {},
	'|':  {},
	'[':  {},
	']':  {},
	'{':  {},
	'}':  {},
	'^':  {},
	'$':  {},
}

func NeedQuoteMeta(r rune) bool {
	_, ok := regexpQuotedBytes[r]
	return ok
}

func QueryParamToRegexpRule(query string, disableColonParam bool) (pathRegexp string, newPath string, regexpRules []string, namedRules map[string]string) {
	fullRule := strings.Builder{}
	rv := strings.Builder{}
	var regExp bool
	var regExpParam bool
	var param bool
	rule := strings.Builder{}
	name := strings.Builder{}
	namedRules = map[string]string{}
	setRuleByParam := func() {
		rv.WriteString("$" + strconv.Itoa(len(regexpRules)+1))
		nameStr := name.String()
		if len(nameStr) > 0 {
			fullRule.WriteString(`(?P<` + nameStr + `>[^/]+)`)
			rule.WriteString(`(?P<` + nameStr + `>[^/]+)`)
			namedRules[nameStr] = rule.String()
		}
		regexpRules = append(regexpRules, rule.String())
		name.Reset()
		rule.Reset()
	}
	for _, r := range query {
		if regExp {
			if r == '>' {
				regExp = false
				regExpParam = false
				fullRule.WriteRune(')')
				rule.WriteRune(')')
				regexpRules = append(regexpRules, rule.String())
				namedRules[name.String()] = rule.String()
				name.Reset()
				rule.Reset()
				continue
			}
			if !regExpParam {
				if r == ':' {
					nameStr := name.String()
					if len(nameStr) > 0 {
						fullRule.WriteString(`?P<` + nameStr + `>`)
						rule.WriteString(`?P<` + nameStr + `>`)
					}
					regExpParam = true
				} else {
					name.WriteRune(r)
				}
				continue
			}
			rule.WriteRune(r)
		} else {
			if r == '*' {
				rv.WriteString("$" + strconv.Itoa(len(regexpRules)+1))
				fullRule.WriteString(`(\S*)`)
				rule.WriteString(`(\S*)`)
				regexpRules = append(regexpRules, rule.String())
				rule.Reset()
				continue
			}
			if !disableColonParam && r == ':' {
				param = true
				continue
			}
			if r == '<' {
				regExp = true
				rv.WriteString("$" + strconv.Itoa(len(regexpRules)+1))
				fullRule.WriteRune('(')
				rule.WriteRune('(')
				continue
			}
			if param {
				if r != '/' {
					name.WriteRune(r)
					continue
				}
				setRuleByParam()
				param = false
			}
			rv.WriteRune(r)
		}
		if !regExp && NeedQuoteMeta(r) {
			fullRule.WriteRune('\\')
		}
		fullRule.WriteRune(r)
	}
	if name.Len() > 0 {
		setRuleByParam()
	}
	pathRegexp = `^` + fullRule.String() + `$`
	newPath = rv.String()
	return
}

// Set rule
func (c *RewriteConfig) Set(urlPath, newPath string) *RewriteConfig {
	re, ok := c.addresses[urlPath]
	if ok {
		delete(c.rulesRegex, re.Old)
		delete(c.rvsesRegex, re.New)
	}
	c.Add(urlPath, newPath)
	return c
}

// Delete rule
func (c *RewriteConfig) Delete(urlPath string) *RewriteConfig {
	re, ok := c.addresses[urlPath]
	if ok {
		delete(c.addresses, urlPath)
		delete(c.rulesRegex, re.Old)
		delete(c.rvsesRegex, re.New)
	}
	return c
}

func ValidateRewriteRule(urlPath, newPath string, disableParam bool) error {
	r, _, ps, kv := QueryParamToRegexpRule(urlPath, disableParam)
	_, err := regexp.Compile(r)
	if err != nil {
		return fmt.Errorf(`%w: %s (routeURL: %s)`, err, r, urlPath)
	}
	newR := regexp.QuoteMeta(newPath)
	if len(ps) > 0 {
		newR = echo.CaptureTokensByValues(ps, kv, true).Replace(newR)
	}
	_, err = regexp.Compile(`^` + newR + `$`)
	if err != nil {
		if len(ps) > 0 {
			err = fmt.Errorf(`%w: %s (rewriteTo: %s)`, err, newR, newPath)
		} else {
			err = fmt.Errorf(`%w: %s`, err, newR)
		}
	}
	return err
}

// Add rule
func (c *RewriteConfig) Add(urlPath, newPath string) *RewriteConfig {
	r, rv, ps, kv := QueryParamToRegexpRule(urlPath, c.DisableColonParam)
	re := regexp.MustCompile(r)
	c.rulesRegex[re] = newPath
	newR := regexp.QuoteMeta(newPath)
	if len(ps) > 0 {
		newR = echo.CaptureTokensByValues(ps, kv, true).Replace(newR)
	}
	rve := regexp.MustCompile(`^` + newR + `$`)
	c.rvsesRegex[rve] = rv
	c.addresses[urlPath] = &RewriteRegExp{
		Old:   re,
		New:   rve,
		Named: kv,
	}
	return c
}

// Rewrite url
func (c *RewriteConfig) Rewrite(urlPath string) string {
	for k, v := range c.rulesRegex {
		replacer := echo.CaptureTokens(k, urlPath)
		if replacer != nil {
			urlPath = replacer.Replace(v)
			break
		}
	}
	return urlPath
}

// Reverse url
func (c *RewriteConfig) Reverse(urlPath string) string {
	for k, v := range c.rvsesRegex {
		replacer := echo.CaptureTokens(k, urlPath)
		if replacer != nil {
			urlPath = replacer.Replace(v)
			break
		}
	}
	return urlPath
}

var (
	// DefaultRewriteConfig is the default Rewrite middleware config.
	DefaultRewriteConfig = RewriteConfig{
		Skipper: echo.DefaultSkipper,
	}
)

// Rewrite returns a Rewrite middleware.
//
// Rewrite middleware rewrites the URL path based on the provided rules.
func Rewrite(rules map[string]string) echo.MiddlewareFuncd {
	c := DefaultRewriteConfig
	c.Rules = rules
	return RewriteWithConfig(c)
}

// RewriteWithConfig returns a Rewrite middleware with config.
// See: `Rewrite()`.
func RewriteWithConfig(config RewriteConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Rules == nil {
		panic("echo: rewrite middleware requires url path rewrite rules")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultRewriteConfig.Skipper
	}
	config.Init()
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			req.URL().SetPath(config.Rewrite(req.URL().Path()))
			return next.Handle(c)
		}
	}
}

func UnrewriteWithConfig(config RewriteConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Rules == nil {
		panic("echo: rewrite middleware requires url path rewrite rules")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultRewriteConfig.Skipper
	}
	config.Init()
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			req.URL().SetPath(config.Reverse(req.URL().Path()))
			return next.Handle(c)
		}
	}
}
