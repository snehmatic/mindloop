package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	v1 "github.com/snehmatic/mindloop/api/v1"
	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/summary"
)

const (
	AppName = "Mindloop"
	Port    = "8765"
)

func CreateRouter(mlh *v1.MindloopHandler) (*mux.Router, error) {
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Routes
	r.HandleFunc("/", mlh.HandleHome).Methods("GET")
	r.HandleFunc("/healthz", mlh.HandleHealthz).Methods("GET")

	// Journal Routes
	r.HandleFunc("/journal", mlh.HandleJournalList).Methods("GET")
	r.HandleFunc("/journal/new", mlh.HandleJournalCreate).Methods("POST")
	r.HandleFunc("/journal/delete", mlh.HandleJournalDelete).Methods("POST")

	// Note Routes
	r.HandleFunc("/notes", mlh.HandleNoteList).Methods("GET")
	r.HandleFunc("/notes/new", mlh.HandleNoteCreate).Methods("POST")
	r.HandleFunc("/notes/view/{id}", mlh.HandleNoteView).Methods("GET")
	r.HandleFunc("/notes/delete", mlh.HandleNoteDelete).Methods("POST")

	// Habit Routes
	r.HandleFunc("/habits", mlh.HandleHabitList).Methods("GET")
	r.HandleFunc("/habits/new", mlh.HandleHabitCreate).Methods("POST")
	r.HandleFunc("/habits/log", mlh.HandleHabitLog).Methods("POST")
	r.HandleFunc("/habits/unlog", mlh.HandleHabitUnlog).Methods("POST")
	r.HandleFunc("/habits/delete", mlh.HandleHabitDelete).Methods("POST")

	// Focus Routes
	r.HandleFunc("/focus", mlh.HandleFocus).Methods("GET")
	r.HandleFunc("/focus/start", mlh.HandleFocusStart).Methods("POST")
	r.HandleFunc("/focus/stop", mlh.HandleFocusStop).Methods("POST")
	r.HandleFunc("/focus/delete", mlh.HandleFocusDelete).Methods("POST")

	// Intent Routes
	r.HandleFunc("/intent", mlh.HandleIntent).Methods("GET")
	r.HandleFunc("/intent/set", mlh.HandleIntentSet).Methods("POST")
	r.HandleFunc("/intent/complete", mlh.HandleIntentComplete).Methods("POST")
	r.HandleFunc("/intent/delete", mlh.HandleIntentDelete).Methods("POST")

	// Summary Route
	r.HandleFunc("/summary", mlh.HandleSummary).Methods("GET")

	// Settings Route
	r.HandleFunc("/settings", mlh.HandleSettings).Methods("GET")
	r.HandleFunc("/settings/update", mlh.HandleSettingsUpdate).Methods("POST")

	// Quote Route
	r.HandleFunc("/api/quote", mlh.HandleQuote).Methods("GET")

	// Maintenance
	r.HandleFunc("/cleanslate", mlh.HandleCleanSlate).Methods("POST")
	r.HandleFunc("/about", mlh.HandleAbout).Methods("GET")
	r.HandleFunc("/void", mlh.HandleVoid).Methods("GET")

	return r, nil
}

func ServeMindloop(mlh *v1.MindloopHandler) {
	r, err := CreateRouter(mlh)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating router")
	}

	appConfig := config.GetConfig()
	srv := &http.Server{
		Addr:      appConfig.Port,
		Handler:   r,
		TLSConfig: nil,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info().Msgf("Starting Mindloop server on %s", appConfig.Port)
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

	// Parse flags
	port := flag.String("port", Port, "Port to run the server on")
	mode := flag.String("mode", "local", "Mode to run the server in (local, byodb, api)")
	flag.Parse()

	// Init global config
	config.InitConfig(AppName, *mode, fmt.Sprintf(":%s", *port))
	appConfig := config.GetConfig()

	database, err := db.ConnectToDb(*appConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB")
	}

	// Initialize core services
	journalService := journal.NewService(database)
	noteService := note.NewService(database)
	focusService := focus.NewService(database)
	intentService := intent.NewService(database)
	summaryService := summary.NewService(database)
	habitService := habit.NewService(database)

	mlh := v1.NewMindloopHandler(
		journalService,
		noteService,
		habitService,
		focusService,
		intentService,
		summaryService,
	)

	ServeMindloop(mlh)
}
