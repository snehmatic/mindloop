package cli

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/repository/focus"
	"github.com/snehmatic/mindloop/internal/repository/habit"
	"github.com/snehmatic/mindloop/internal/repository/habitlog"
	"github.com/snehmatic/mindloop/internal/repository/intent"
	"github.com/snehmatic/mindloop/internal/repository/journal"
	"github.com/snehmatic/mindloop/internal/repository/note"
	"github.com/snehmatic/mindloop/internal/repository/point"
	"github.com/snehmatic/mindloop/internal/repository/quest"
	"github.com/snehmatic/mindloop/internal/repository/routine"
	"github.com/snehmatic/mindloop/internal/repository/subtask"
	taskRepo "github.com/snehmatic/mindloop/internal/repository/task"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var gdb *gorm.DB
var ac *config.Config

// repoHolders stores shared repository and service instances created after DB init.
// Each field has a pointer so nil-checking is trivial.
var (
	fRepo     *focus.Repository
	hRepo     *habit.Repository
	hlRepo    *habitlog.Repository
	iRepo     *intent.Repository
	jRepo     *journal.Repository
	nRepo     *note.Repository
	pRepo     *point.Repository
	qRepo     *quest.Repository
	rRepo     *routine.Repository
	sRepo     *subtask.Repository
	tRepo     *taskRepo.TaskRepository
	logger    log.Logger
	reposDone bool
)

func setupRepos(db *gorm.DB) {
	if reposDone {
		return
	}
	logger = log.Get()
	fRepo = ptr(focus.NewSQLRepository(db))
	hRepo = ptr(habit.NewSQLRepository(db, logger))
	hlRepo = ptr(habitlog.NewSQLRepository(db))
	iRepo = ptr(intent.NewSQLRepository(db))
	jRepo = ptr(journal.NewSQLRepository(db))
	nRepo = ptr(note.NewSQLRepository(db))
	pRepo = ptr(point.NewSQLRepository(db))
	qRepo = ptr(quest.NewSQLRepository(db, logger))
	rRepo = ptr(routine.NewSQLRepository(db))
	sRepo = ptr(subtask.NewSQLRepository(db))
	tRepo = ptr(taskRepo.NewSQLTaskRepository(db))
	reposDone = true
}

func ptr[T any](v T) *T { return &v }

var rootCmd = &cobra.Command{
	Use:       "mindloop",
	Version:   config.Version,
	Short:     "mindloop is a CLI tool for productivity tracking",
	Long:      `Mindloop helps track intent, focus sessions, and habits via CLI.`,
	Example:   `mindloop intent start "Get this work done"`,
	ValidArgs: []string{"intent", "focus", "habit", "log", "stats"},
	Args:      cobra.OnlyValidArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.ValidateUserConfig(cmd)

		if db.LocalDBFileExists() {
			ac.Logger.Info("Found local DB file, using it for local mode.")
		} else {
			ac.Logger.Warn("No local DB file found, a new one will be created.")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrettyPrintBanner()
		if ac.UserName != "" {
			utils.PrintRocketln(fmt.Sprintf("Welcome back, %s! Use 'mindloop help' to see available commands.", ac.UserName))
		} else {
			utils.PrintRocketln("Welcome to Mindloop! Use 'mindloop help' to see available commands.")
		}
		utils.PrintInfoln("For starters, try 'mindloop configure' to set up your profile.")
		ac.Logger.Info("User accessed root command, prompting for help.")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧠 Thank you for using Mindloop! This is still a work in progress.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)
	// define persistent flags here
}

func initLogger() {
	logPath := "mindloop.log"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logPath = config.GetDataDir() + "/mindloop.log"
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	opts := &log.InitOptions{
		Out:   logFile,
		Level: zerolog.DebugLevel,
	}
	log.Init(opts)
}

func initConfig() {
	ac = config.GetConfig()
	// Initialize local db
	dbConnection, err := db.ConnectToDb(*ac)
	if err != nil {
		utils.PrintErrorf("Error connecting to DB: %v\n", err)
		ac.Logger.Error("Error connecting to DB", err)
		utils.PrintErrorln("Please check your database connection or configuration.")
		ac.Logger.Warn("Exiting due to DB connection error.")
		os.Exit(1)
	}
	gdb = dbConnection
	setupRepos(dbConnection)
}
