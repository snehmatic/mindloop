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

func Get() zerolog.Logger {
	return instance
}
