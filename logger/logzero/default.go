package logzero

import "io"

func Debug(s ...any) {
	Default().Debug(s...)
}

func Debugf(t string, s ...any) {
	Default().Debugf(t, s...)
}

func Info(s ...any) {
	Default().Info(s...)
}

func Infof(t string, s ...any) {
	Default().Infof(t, s...)
}

func Warn(s ...any) {
	Default().Warn(s...)
}

func Warnf(t string, s ...any) {
	Default().Warnf(t, s...)
}

func Error(s ...any) {
	Default().Error(s...)
}

func Errorf(t string, s ...any) {
	Default().Errorf(t, s...)
}

func Fatal(s ...any) {
	Default().Fatal(s...)
}

func Fatalf(t string, s ...any) {
	Default().Fatalf(t, s...)
}

func GetLogger(category string, writers ...io.Writer) *Logger {
	return Default().GetLogger(category, writers...)
}

func Writer() io.Writer {
	return Default()
}

func SetLevel(level string) {
	Default().SetLevel(level)
}
