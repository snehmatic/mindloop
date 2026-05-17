// Package log provides a centralized logging utility for Mindloop
package log

import (
	"io"
	"sync"

	"github.com/rs/zerolog"
)

var (
	once     sync.Once
	instance zerolog.Logger
)

// Init initializes the global logger with the specified output and level
func Init(out io.Writer, level zerolog.Level) {
	once.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		instance = zerolog.New(out).
			With().
			Timestamp().
			Caller().
			Logger().
			Level(level)
	})
}

// Get returns the global logger instance
func Get() *zerolog.Logger {
	//guard against returning an uninitialized logger.
	if instance.GetLevel() == zerolog.NoLevel {
		panic("log: logger not initialized")
	}
	return &instance
}
