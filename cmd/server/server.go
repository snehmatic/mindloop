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
	v1 "github.com/snehmatic/mindloop/api/v1"
	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/backup"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/core/quest"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/core/task"
	"github.com/snehmatic/mindloop/internal/log"
	appsettings "github.com/snehmatic/mindloop/internal/repository/appsettings"
	focusRepo "github.com/snehmatic/mindloop/internal/repository/focus"
	habitRepo "github.com/snehmatic/mindloop/internal/repository/habit"
	"github.com/snehmatic/mindloop/internal/repository/habitlog"
	habitIntRepo "github.com/snehmatic/mindloop/internal/repository/intent"
	journalRepo "github.com/snehmatic/mindloop/internal/repository/journal"
	noteRepo "github.com/snehmatic/mindloop/internal/repository/note"
	pointRepo "github.com/snehmatic/mindloop/internal/repository/point"
	questRepo "github.com/snehmatic/mindloop/internal/repository/quest"
	routineRepo "github.com/snehmatic/mindloop/internal/repository/routine"
	subtaskRepo "github.com/snehmatic/mindloop/internal/repository/subtask"
	taskRepo "github.com/snehmatic/mindloop/internal/repository/task"
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

	// Task Routes
	r.HandleFunc("/tasks", mlh.HandleTaskList).Methods("GET")
	r.HandleFunc("/tasks/new", mlh.HandleTaskCreate).Methods("POST")
	r.HandleFunc("/tasks/complete", mlh.HandleTaskComplete).Methods("POST")
	r.HandleFunc("/tasks/delete", mlh.HandleTaskDelete).Methods("POST")
	r.HandleFunc("/tasks/subtask/new", mlh.HandleSubtaskCreate).Methods("POST")
	r.HandleFunc("/tasks/subtask/complete", mlh.HandleSubtaskComplete).Methods("POST")
	r.HandleFunc("/tasks/subtask/delete", mlh.HandleSubtaskDelete).Methods("POST")

	// Settings Route
	r.HandleFunc("/settings", mlh.HandleSettings).Methods("GET")
	r.HandleFunc("/settings/update", mlh.HandleSettingsUpdate).Methods("POST")

	// AI Routes
	r.HandleFunc("/api/v1/ai/settings", mlh.HandleGetAISettings).Methods("GET")
	r.HandleFunc("/api/v1/ai/settings", mlh.HandleSaveAISettings).Methods("POST")
	r.HandleFunc("/api/v1/ai/models", mlh.HandleListAIModels).Methods("GET")
	r.HandleFunc("/api/v1/ai/test", mlh.HandleTestAIConnection).Methods("POST")
	r.HandleFunc("/api/v1/ai/generate", mlh.HandleGenerateAIJournal).Methods("GET")

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
		log.Get().Info(fmt.Sprintf("Starting Mindloop server on %s", appConfig.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Get().Error("ListenAndServe() error", err)
			os.Exit(1)
		}
	}()

	<-stop
	log.Get().Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Get().Error("Server Shutdown Failed", err)
	}
	log.Get().Info("Server exited properly")
}

func main() {
	// Parse flags
	port := flag.String("port", Port, "Port to run the server on")
	mode := flag.String("mode", "local", "Mode to run the server in (local, byodb, api)")
	flag.Parse()

	// Init global config
	if err := config.Init(AppName, *mode, fmt.Sprintf(":%s", *port)); err != nil {
		log.Get().Error("Failed to initialize config", err)
		os.Exit(1)
	}
	appConfig := config.GetConfig()

	database, err := db.ConnectToDb(*appConfig)
	if err != nil {
		log.Get().Fatal("Error connecting to DB", err)
	}

	logger := log.Get()

	// Initialize repository instances
	fRepo := focusRepo.NewSQLRepository(database)
	hRepo := habitRepo.NewSQLRepository(database, logger)
	hlRepo := habitlog.NewSQLRepository(database)
	iRepo := habitIntRepo.NewSQLRepository(database)
	jRepo := journalRepo.NewSQLRepository(database)
	nRepo := noteRepo.NewSQLRepository(database)
	pRepo := pointRepo.NewSQLRepository(database)
	qRepo := questRepo.NewSQLRepository(database, logger)
	rRepo := routineRepo.NewSQLRepository(database)
	sRepo := subtaskRepo.NewSQLRepository(database)
	tRepo := taskRepo.NewSQLTaskRepository(database)

	// Initialize core services with repository injection
	pointSvc := points.NewService(pRepo)
	journalService := journal.NewService(jRepo)
	noteService := note.NewService(nRepo)
	backupService := backup.NewService(database, pointSvc, fRepo, hRepo, hlRepo, iRepo, jRepo, nRepo, pRepo, qRepo, rRepo, sRepo, tRepo)
	focusService := focus.NewService(fRepo)
	intentService := intent.NewService(iRepo)
	questService := quest.NewService(qRepo, logger)
	summaryService := summary.NewService(fRepo, hRepo, iRepo, pRepo, tRepo, logger)
	habitService := habit.NewService(hRepo)
	taskService := task.NewService(tRepo, config.GetUserConfig(), pointSvc, logger)
	appSettingsRepo := appsettings.NewSQLRepository(database)

	mlh := v1.NewMindloopHandler(
		database,
		appSettingsRepo,
		journalService,
		noteService,
		habitService,
		focusService,
		intentService,
		questService,
		summaryService,
		backupService,
		taskService,
	)

	ServeMindloop(mlh)
}

func applyMilestoneInterval() {
	uc := config.UserConfig{}
	_ = uc.ReadFromYAML()
	points.SetMilestoneInterval(uc.PointsConfig.MilestoneInterval)
}
