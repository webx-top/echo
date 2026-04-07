package logger

import (
	"io"
	"os"

	isatty "github.com/admpub/go-isatty"
)

var _ Logger = &Base{}

type (
	// Logger is the interface that declares Echo's logging system.
	Logger interface {
		Debug(...any)
		Debugf(string, ...any)

		Info(...any)
		Infof(string, ...any)

		Warn(...any)
		Warnf(string, ...any)

		Error(...any)
		Errorf(string, ...any)

		Fatal(...any)
		Fatalf(string, ...any)
	}

	LevelSetter interface {
		SetLevel(string)
	}

	Base struct {
	}
)

func (b *Base) Debug(...any) {
}

func (b *Base) Debugf(string, ...any) {
}

func (b *Base) Info(...any) {
}

func (b *Base) Infof(string, ...any) {
}

func (b *Base) Warn(...any) {
}

func (b *Base) Warnf(string, ...any) {
}

func (b *Base) Error(...any) {
}

func (b *Base) Errorf(string, ...any) {
}

func (b *Base) Fatal(...any) {
}

func (b *Base) Fatalf(string, ...any) {
}

func (b *Base) SetLevel(string) {
}

func Colorable(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	if isatty.IsTerminal(file.Fd()) {
		return true
	}
	if isatty.IsCygwinTerminal(file.Fd()) {
		return true
	}
	return false
}
