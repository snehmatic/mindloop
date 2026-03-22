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
	"github.com/snehmatic/mindloop/internal/core/backup"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/quest"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/web"
)

const (
	AppName = "Mindloop"
	Port    = "8765"
)

func CreateRouter(mlh *v1.MindloopHandler) *mux.Router {
	r := mux.NewRouter()

	// Static files from embedded FS
	staticFS := http.FS(web.WebFS)
	r.PathPrefix("/static/").Handler(http.FileServer(staticFS))

	// Routes
	r.HandleFunc("/", mlh.HandleHome).Methods("GET")
	r.HandleFunc("/healthz", mlh.HandleHealthz).Methods("GET")

	// Journal Routes
	r.HandleFunc("/journal", mlh.HandleJournalList).Methods("GET")
	r.HandleFunc("/journal/new", mlh.HandleJournalCreate).Methods("POST")
	r.HandleFunc("/journal/view/{id}", mlh.HandleJournalView).Methods("GET")
	r.HandleFunc("/journal/update", mlh.HandleJournalUpdate).Methods("POST")
	r.HandleFunc("/journal/update-live", mlh.HandleJournalUpdateLive).Methods("POST", "PUT")
	r.HandleFunc("/journal/delete", mlh.HandleJournalDelete).Methods("POST")

	// Note Routes
	r.HandleFunc("/notes", mlh.HandleNoteList).Methods("GET")
	r.HandleFunc("/notes/new", mlh.HandleNoteCreate).Methods("POST")
	r.HandleFunc("/notes/view/{id}", mlh.HandleNoteView).Methods("GET")
	r.HandleFunc("/notes/update-live", mlh.HandleNoteUpdateLive).Methods("POST", "PUT")
	r.HandleFunc("/notes/delete", mlh.HandleNoteDelete).Methods("POST")

	// Habit Routes
	r.HandleFunc("/habits", mlh.HandleHabitList).Methods("GET")
	r.HandleFunc("/habits/new", mlh.HandleHabitCreate).Methods("POST")
	r.HandleFunc("/habits/view/{id}", mlh.HandleHabitView).Methods("GET")
	r.HandleFunc("/habits/update", mlh.HandleHabitUpdate).Methods("POST")
	r.HandleFunc("/habits/log", mlh.HandleHabitLog).Methods("POST")
	r.HandleFunc("/habits/unlog", mlh.HandleHabitUnlog).Methods("POST")
	r.HandleFunc("/habits/delete", mlh.HandleHabitDelete).Methods("POST")

	// Focus Routes
	r.HandleFunc("/focus", mlh.HandleFocus).Methods("GET")
	r.HandleFunc("/focus/start", mlh.HandleFocusStart).Methods("POST")
	r.HandleFunc("/focus/update", mlh.HandleFocusUpdate).Methods("POST")
	r.HandleFunc("/focus/stop", mlh.HandleFocusStop).Methods("POST")
	r.HandleFunc("/focus/delete", mlh.HandleFocusDelete).Methods("POST")

	// Intent Routes
	r.HandleFunc("/intent", mlh.HandleIntent).Methods("GET")
	r.HandleFunc("/intent/set", mlh.HandleIntentSet).Methods("POST")
	r.HandleFunc("/intent/update", mlh.HandleIntentUpdate).Methods("POST")
	r.HandleFunc("/intent/complete", mlh.HandleIntentComplete).Methods("POST")
	r.HandleFunc("/intent/delete", mlh.HandleIntentDelete).Methods("POST")
	r.HandleFunc("/intent/resume", mlh.HandleIntentResume).Methods("POST")

	// Quest Routes
	r.HandleFunc("/quest/start", mlh.HandleQuestStart).Methods("POST")
	r.HandleFunc("/quest/stop", mlh.HandleQuestStop).Methods("POST")
	r.HandleFunc("/quest/delete", mlh.HandleQuestDelete).Methods("POST")

	// Summary Route
	r.HandleFunc("/summary", mlh.HandleSummary).Methods("GET")

	// Settings Route
	r.HandleFunc("/settings", mlh.HandleSettings).Methods("GET")
	r.HandleFunc("/settings/update", mlh.HandleSettingsUpdate).Methods("POST")

	// Backup Routes
	r.HandleFunc("/backup/export", mlh.HandleBackupExport).Methods("GET")
	r.HandleFunc("/backup/import", mlh.HandleBackupImport).Methods("POST")

	// Quote Route
	r.HandleFunc("/api/quote", mlh.HandleQuote).Methods("GET")

	// Maintenance
	r.HandleFunc("/cleanslate", mlh.HandleCleanSlate).Methods("POST")
	r.HandleFunc("/about", mlh.HandleAbout).Methods("GET")
	r.HandleFunc("/void", mlh.HandleVoid).Methods("GET")

	return r
}

func ServeMindloop(mlh *v1.MindloopHandler) {
	r := CreateRouter(mlh)

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
		log.Error().Msgf("Server Shutdown Failed:%+v", err)
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
	backupService := backup.NewService(database)
	focusService := focus.NewService(database)
	intentService := intent.NewService(database)
	questService := quest.NewService(database)
	summaryService := summary.NewService(database)
	habitService := habit.NewService(database)

	mlh := v1.NewMindloopHandler(
		journalService,
		noteService,
		habitService,
		focusService,
		intentService,
		questService,
		summaryService,
		backupService,
	)

	ServeMindloop(mlh)
}
