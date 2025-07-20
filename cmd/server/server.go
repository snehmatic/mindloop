package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	v1 "github.com/snehmatic/mindloop/api/v1"
	"github.com/snehmatic/mindloop/config"
	"github.com/snehmatic/mindloop/db"
)

const (
	AppName = "Mindloop"
	Port    = "8080"
)

func CreateRouter(mlh *v1.MindloopHandler) (*mux.Router, error) {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/", mlh.HandleHome)
	r.HandleFunc("/healthz", mlh.HandleHealthz)

	return r, nil
}

func ServeMindloop(mlh *v1.MindloopHandler) {
	r, err := CreateRouter(mlh)
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
	config.InitConfig(AppName, "api", fmt.Sprintf(":"+Port))
	appConfig := config.GetConfig()

	_, err := db.ConnectToDb(*appConfig) // to be used later
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB")
	}

	mlh := v1.NewMindloopHandler()

	ServeMindloop(mlh)
}
