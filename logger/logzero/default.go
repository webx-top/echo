package logzero

import "io"

func Debug(s ...interface{}) {
	Default().Debug(s...)
}

func Debugf(t string, s ...interface{}) {
	Default().Debugf(t, s...)
}

func Info(s ...interface{}) {
	Default().Info(s...)
}

func Infof(t string, s ...interface{}) {
	Default().Infof(t, s...)
}

func Warn(s ...interface{}) {
	Default().Warn(s...)
}

func Warnf(t string, s ...interface{}) {
	Default().Warnf(t, s...)
}

func Error(s ...interface{}) {
	Default().Error(s...)
}

func Errorf(t string, s ...interface{}) {
	Default().Errorf(t, s...)
}

func Fatal(s ...interface{}) {
	Default().Fatal(s...)
}

func Fatalf(t string, s ...interface{}) {
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
