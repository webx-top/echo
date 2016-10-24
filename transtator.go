package echo

import "fmt"

type Translator interface {
	T(format string, args ...interface{}) string
}

var DefaultNopTranslate Translator = &NopTranslate{}

type NopTranslate struct {
}

func (n *NopTranslate) T(format string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}
