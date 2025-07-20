package main

import (
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
	logFile, err := os.OpenFile("mindloop.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	defer logFile.Close()
	log.Init(logFile, zerolog.DebugLevel)
	logger := log.Get()
	logger.Info().Msg("Logging to mindloop.log file...")

	// Init global config
	config.InitConfig(AppName, "local", "")

	cli.Execute()
}
