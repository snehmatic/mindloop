package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	cli "github.com/snehmatic/mindloop/cmd/cli"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/log"
)

const (
	AppName = "Mindloop"
)

func main() {
	logPath := "mindloop.log"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logPath = config.GetDataDir() + "/mindloop.log"
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			panic("Failed to close log file: " + err.Error())
		}
	}()
	log.Init(logFile, zerolog.DebugLevel)
	logger := log.Get()
	logger.Info().Msgf("Logging to %s file...", logPath)

	// Init global config
	if err := config.Init(AppName, "local", ""); err != nil {
		panic(fmt.Sprintf("Failed to initialize config: %v", err))
	}

	cli.Execute()
}
