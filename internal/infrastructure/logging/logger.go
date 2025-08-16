package logging

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	logger zerolog.Logger
	once   sync.Once
)

// InitLogger initializes the global logger
func InitLogger() {
	once.Do(func() {
		// Create or open log file
		logFile, err := os.OpenFile("mindloop.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fallback to stderr if file creation fails
			logger = zerolog.New(os.Stderr).
				Level(zerolog.InfoLevel).
				With().
				Timestamp().
				Logger()
			return
		}

		logger = zerolog.New(logFile).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	})
}

// GetLogger returns the global logger instance
func GetLogger() zerolog.Logger {
	InitLogger()
	return logger
}

// SetLevel sets the logging level
func SetLevel(level zerolog.Level) {
	logger = logger.Level(level)
}

// EnableDebug enables debug level logging
func EnableDebug() {
	SetLevel(zerolog.DebugLevel)
}

// EnableInfo enables info level logging
func EnableInfo() {
	SetLevel(zerolog.InfoLevel)
}
