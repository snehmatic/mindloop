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
	"github.com/snehmatic/mindloop/internal/application"
	"github.com/snehmatic/mindloop/internal/infrastructure/config"
)

const (
	AppName = "Mindloop"
	Port    = "8080"
)

func CreateRouter(mlh *v1.MindloopHandler) (*mux.Router, error) {
	r := mux.NewRouter()

	// Add CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// API v1 routes
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Health and info routes
	apiV1.HandleFunc("/", mlh.HandleHome).Methods("GET")
	apiV1.HandleFunc("/healthz", mlh.HandleHealthz).Methods("GET")

	// Habit routes
	apiV1.HandleFunc("/habits", mlh.HandleCreateHabit).Methods("POST")
	apiV1.HandleFunc("/habits", mlh.HandleListHabits).Methods("GET")
	apiV1.HandleFunc("/habits/{id}", mlh.HandleGetHabit).Methods("GET")
	apiV1.HandleFunc("/habits/{id}", mlh.HandleDeleteHabit).Methods("DELETE")
	apiV1.HandleFunc("/habits/{id}/log", mlh.HandleLogHabit).Methods("POST")

	// Intent routes
	apiV1.HandleFunc("/intents", mlh.HandleCreateIntent).Methods("POST")
	apiV1.HandleFunc("/intents", mlh.HandleListIntents).Methods("GET")
	apiV1.HandleFunc("/intents/{id}/end", mlh.HandleEndIntent).Methods("POST")
	apiV1.HandleFunc("/intents/{id}", mlh.HandleDeleteIntent).Methods("DELETE")

	// Focus routes
	apiV1.HandleFunc("/focus", mlh.HandleCreateFocus).Methods("POST")
	apiV1.HandleFunc("/focus", mlh.HandleListFocus).Methods("GET")
	apiV1.HandleFunc("/focus/{id}/end", mlh.HandleEndFocus).Methods("POST")
	apiV1.HandleFunc("/focus/{id}/pause", mlh.HandlePauseFocus).Methods("POST")
	apiV1.HandleFunc("/focus/{id}/resume", mlh.HandleResumeFocus).Methods("POST")
	apiV1.HandleFunc("/focus/{id}/rate", mlh.HandleRateFocus).Methods("POST")
	apiV1.HandleFunc("/focus/{id}", mlh.HandleDeleteFocus).Methods("DELETE")

	// Journal routes
	apiV1.HandleFunc("/journal", mlh.HandleCreateJournal).Methods("POST")
	apiV1.HandleFunc("/journal", mlh.HandleListJournal).Methods("GET")
	apiV1.HandleFunc("/journal/{id}", mlh.HandleGetJournal).Methods("GET")
	apiV1.HandleFunc("/journal/{id}", mlh.HandleUpdateJournal).Methods("PUT")
	apiV1.HandleFunc("/journal/{id}", mlh.HandleDeleteJournal).Methods("DELETE")

	// Summary routes
	apiV1.HandleFunc("/summary/daily", mlh.HandleDailySummary).Methods("GET")
	apiV1.HandleFunc("/summary/weekly", mlh.HandleWeeklySummary).Methods("GET")
	apiV1.HandleFunc("/summary/monthly", mlh.HandleMonthlySummary).Methods("GET")
	apiV1.HandleFunc("/summary/yearly", mlh.HandleYearlySummary).Methods("GET")
	apiV1.HandleFunc("/summary/custom", mlh.HandleCustomSummary).Methods("POST")

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
	// Initialize configuration with local mode for SQLite
	config.InitConfig(AppName, "local", fmt.Sprintf(":"+Port))
	appConfig := config.GetConfig()

	// Initialize dependency injection container
	container, err := application.NewContainer(appConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application container")
	}
	defer container.Close()

	mlh := v1.NewMindloopHandler(container)

	ServeMindloop(mlh)
}
