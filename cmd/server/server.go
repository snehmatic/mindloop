package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/config"
	"github.com/snehmatic/mindloop/db"
)

const (
	AppName = "MindLoop"
	Port    = "8080"
)

func ServeMindloop() {
	r, err := CreateRouter()
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating router")
	}

	srv := &http.Server{
		Addr:      ":8080",
		Handler:   r,
		TLSConfig: nil,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info().Msg("Starting Mindloop server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("ListenAndServe(): %v", err)
		}
	}()

	<-stop
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Server Shutdown Failed:%+v", err)
	}
	log.Info().Msg("Server exited properly")
}

func main() {

	// Init global config
	config.InitConfig(AppName, fmt.Sprintf(":"+Port))
	appConfig := config.GetConfig()

	dbConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBConfig.Host,
		appConfig.DBConfig.Port,
		appConfig.DBConfig.User,
		appConfig.DBConfig.Password,
		appConfig.DBConfig.Name,
	)

	_, err := db.Conn(dbConnString) // to be used later
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB")
	}

	ServeMindloop()
}
