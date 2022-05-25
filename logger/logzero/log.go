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
	CallerSkipFrameCount = 4
	WriterDefaultLevel   = zerolog.InfoLevel
	global               *Logger
	_                    logger.Logger = &Logger{}
	once                 sync.Once
)

func initGlobal() {
	global = New()
}

// Default default global logger
func Default() *Logger {
	once.Do(initGlobal)
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
	if len(writers) == 0 {
		writers = append(writers, NewConsoleOutput())
	}
	o := zerolog.MultiLevelWriter(writers...)
	var l zerolog.Logger
	if CallerSkipFrameCount < 0 {
		l = zerolog.New(o).With().Timestamp().Logger()
	} else {
		l = zerolog.New(o).With().CallerWithSkipFrameCount(CallerSkipFrameCount).Timestamp().Logger()
	}
	return &Logger{
		Logger:      &l,
		Base:        &logger.Base{},
		mutex:       &sync.RWMutex{},
		WriterLevel: WriterDefaultLevel,
		writers:     writers,
		subs:        make(map[string]*Logger),
	}
}

func NewSub(root *Logger, category string, writers ...io.Writer) *Logger {
	subLogger := root.Logger.With().Str("category", category).CallerWithSkipFrameCount(CallerSkipFrameCount - 1).Logger()
	if len(writers) > 0 {
		o := zerolog.MultiLevelWriter(writers...)
		subLogger = subLogger.Output(o)
	}
	return &Logger{
		Logger:      &subLogger,
		Base:        root.Base,
		mutex:       root.mutex,
		WriterLevel: WriterDefaultLevel,
		writers:     writers,
		subs:        make(map[string]*Logger),
	}
}

type Logger struct {
	*zerolog.Logger
	*logger.Base
	WriterLevel zerolog.Level
	writers     []io.Writer
	subs        map[string]*Logger
	mutex       *sync.RWMutex
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
		subLogger = NewSub(a, category, writers...)
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
	a.Logger.WithLevel(a.WriterLevel).CallerSkipFrame(-1).Msg(string(p))
	return n, nil
}

func (a *Logger) SetLevel(level string) {
	level = strings.ToLower(level)
	lv, err := zerolog.ParseLevel(level)
	if err == nil {
		a.Logger.Level(lv)
	}
}
