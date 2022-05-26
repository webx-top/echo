package logzero

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/webx-top/echo/logger"
)

var (
	// CallerSkipFrameCount caller skip frame count
	CallerSkipFrameCount               = 4
	WriterDefaultLevel                 = zerolog.InfoLevel
	global               *Logger       = New()
	_                    logger.Logger = &Logger{}
)

func SetGlobal(l *Logger) {
	global = l
}

// Default default global logger
func Default() *Logger {
	return global
}

func NewConsoleOutput() zerolog.ConsoleWriter {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: `2006-01-02T15:04:05`,
		NoColor:    !logger.Colorable(os.Stdout),
	}
	// output.FormatLevel = func(i interface{}) string {
	// 	return fmt.Sprintf("| %-6s|", i)
	// }
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	return output
}

func New(writers ...io.Writer) *Logger {
	return NewLogger(CallerSkipFrameCount, writers...)
}

func NewLogger(callerSkipFrameCount int, writers ...io.Writer) *Logger {
	if len(writers) == 0 {
		writers = append(writers, NewConsoleOutput())
	}
	o := zerolog.MultiLevelWriter(writers...)
	var withCaller bool
	var l zerolog.Logger
	if callerSkipFrameCount < 0 {
		l = zerolog.New(o).With().Timestamp().Logger()
	} else {
		withCaller = true
		l = zerolog.New(o).With().CallerWithSkipFrameCount(callerSkipFrameCount).Timestamp().Logger()
	}
	return &Logger{
		Logger:                 &l,
		Base:                   &logger.Base{},
		mutex:                  &sync.RWMutex{},
		WriterLevel:            WriterDefaultLevel,
		writers:                writers,
		writerCallerSkipFrames: 1,
		withCaller:             withCaller,
		subs:                   make(map[string]*Logger),
	}
}

func newSub(root *Logger, category string, writers ...io.Writer) *Logger {
	subLogger := NewLogger(CallerSkipFrameCount-1, writers...)
	l := subLogger.With().Str("category", category).Logger()
	subLogger.Logger = &l
	subLogger.writerCallerSkipFrames = 2
	subLogger.withCaller = root.withCaller
	return subLogger
}

type Logger struct {
	*zerolog.Logger
	*logger.Base
	WriterLevel            zerolog.Level
	writers                []io.Writer
	writerCallerSkipFrames int
	withCaller             bool
	subs                   map[string]*Logger
	mutex                  *sync.RWMutex
}

func (a *Logger) Debug(s ...interface{}) {
	a.Logger.Debug().Msg(fmt.Sprint(s...))
}

func (a *Logger) Debugf(t string, s ...interface{}) {
	a.Logger.Debug().Msgf(t, s...)
}

func (a *Logger) Info(s ...interface{}) {
	a.Logger.Info().Msg(fmt.Sprint(s...))
}

func (a *Logger) Infof(t string, s ...interface{}) {
	a.Logger.Info().Msgf(t, s...)
}

func (a *Logger) Warn(s ...interface{}) {
	a.Logger.Warn().Msg(fmt.Sprint(s...))
}

func (a *Logger) Warnf(t string, s ...interface{}) {
	a.Logger.Warn().Msgf(t, s...)
}

func (a *Logger) Error(s ...interface{}) {
	a.Logger.Error().Msg(fmt.Sprint(s...))
}

func (a *Logger) Errorf(t string, s ...interface{}) {
	a.Logger.Error().Msgf(t, s...)
}

func (a *Logger) Fatal(s ...interface{}) {
	a.Logger.Fatal().Msg(fmt.Sprint(s...))
}

func (a *Logger) Fatalf(t string, s ...interface{}) {
	a.Logger.Fatal().Msgf(t, s...)
}

func (a *Logger) GetLogger(category string, writers ...io.Writer) *Logger {
	a.mutex.RLock()
	subLogger, ok := a.subs[category]
	a.mutex.RUnlock()
	if !ok {
		subLogger = newSub(a, category, writers...)
		a.mutex.Lock()
		a.subs[category] = subLogger
		a.mutex.Unlock()
	}
	return subLogger
}

func (a *Logger) Write(p []byte) (int, error) {
	n := len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	if a.withCaller {
		a.Logger.WithLevel(a.WriterLevel).CallerSkipFrame(a.writerCallerSkipFrames).Msg(string(p))
	} else {
		a.Logger.WithLevel(a.WriterLevel).Msg(string(p))
	}
	return n, nil
}

func (a *Logger) SetLevel(level string) {
	level = strings.ToLower(level)
	lv, err := zerolog.ParseLevel(level)
	if err == nil {
		a.Logger.Level(lv)
	}
}
