// Package log provides a centralized logging utility for Mindloop
package log

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	once     sync.Once
	instance *logger
)

type logger struct {
	log zerolog.Logger
}

type Field struct {
	Key   string
	Value any
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, err error, fields ...Field)
	Fatal(msg string, err error, fields ...Field)
	With(fields ...Field) Logger
}

type InitOptions struct {
	Out   io.Writer
	Level zerolog.Level
}

// Init initializes the global logger with the specified output and level
func Init(opts *InitOptions) {
	once.Do(func() {
		out := io.Writer(os.Stdout)
		level := zerolog.DebugLevel

		if opts != nil {
			if opts.Out != nil {
				out = opts.Out
			}
			if opts.Level != 0 {
				level = opts.Level
			}
		}

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		l := zerolog.New(out).
			With().
			Timestamp().
			Caller().
			Logger().
			Level(level)
		instance = &logger{log: l}
	})
}

// Get returns the global logger instance as a Logger interface
func Get() Logger {
	if instance == nil {
		// Auto-initialize with default values for ease of use in tests
		Init(&InitOptions{
			Level: zerolog.DebugLevel,
			Out:   io.Discard,
		})
	}
	return instance
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.applyFields(l.log.Debug(), fields...).Msg(msg)
}

func (l *logger) Info(msg string, fields ...Field) {
	l.applyFields(l.log.Info(), fields...).Msg(msg)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.applyFields(l.log.Warn(), fields...).Msg(msg)
}

func (l *logger) Error(msg string, err error, fields ...Field) {
	e := l.log.Error().Err(err)
	l.applyFields(e, fields...).Msg(msg)
}

func (l *logger) Fatal(msg string, err error, fields ...Field) {
	l.applyFields(l.log.Fatal(), fields...).Err(err).Msg(msg)
}

func (l *logger) With(fields ...Field) Logger {
	ctx := l.log.With()
	for _, f := range fields {
		ctx = ctx.Interface(f.Key, f.Value)
	}
	return &logger{log: ctx.Logger()}
}

func (l *logger) applyFields(e *zerolog.Event, fields ...Field) *zerolog.Event {
	for _, f := range fields {
		e = e.Interface(f.Key, f.Value)
	}
	return e
}
