package zerolog

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
	// Default default global logger
	Default = New()

	_ logger.Logger = Default
)

func init() {
	zerolog.CallerSkipFrameCount = 3
}

func New(writers ...io.Writer) *Logger {
	if len(writers) == 0 {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	o := zerolog.MultiLevelWriter(writers...)
	l := zerolog.New(o).With().Caller().Timestamp().Logger()
	return &Logger{
		Logger: &l,
		Base:   &logger.Base{},
		mutex:  &sync.RWMutex{},
		subs:   make(map[string]*Logger),
	}
}

func NewSub(root *Logger, category string, writers ...io.Writer) *Logger {
	subLogger := root.Logger.With().Str("category", category).Logger()
	if len(writers) > 0 {
		o := zerolog.MultiLevelWriter(writers...)
		subLogger = subLogger.Output(o)
	}
	return &Logger{
		Logger: &subLogger,
		Base:   root.Base,
		mutex:  root.mutex,
		subs:   make(map[string]*Logger),
	}
}

type Logger struct {
	*zerolog.Logger
	*logger.Base
	subs  map[string]*Logger
	mutex *sync.RWMutex
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
	a.mutex.Lock()
	subLogger, ok := a.subs[category]
	if !ok {
		subLogger = NewSub(a, category, writers...)
		a.subs[category] = subLogger
	}
	a.mutex.Unlock()
	return subLogger
}

func (a *Logger) SetLevel(level string) {
	level = strings.ToLower(level)
	lv, err := zerolog.ParseLevel(level)
	if err == nil {
		a.Logger.Level(lv)
	}
}
