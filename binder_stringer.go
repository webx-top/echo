package echo

import "github.com/webx-top/echo/param"

var (
	TimestampStringer  = param.TimestampStringer
	DateTimeStringer   = param.DateTimeStringer
	WhitespaceStringer = param.WhitespaceStringer
	Ignored            = param.Ignored
)

func TranslateStringer(t Translator, args ...any) param.Stringer {
	return param.StringerFunc(func(v any) string {
		return t.T(param.AsString(v), args...)
	})
}

func ignoreValues(v any) []string {
	return nil
}

func FormStringer(s param.Stringer) BinderValueCustomEncoder {
	if ig, ok := s.(param.Ignorer); ok {
		if ig.Ignore() {
			return ignoreValues
		}
	}
	return func(v any) []string {
		return []string{s.String(v)}
	}
}
