package cli

import (
	v1 "github.com/snehmatic/mindloop/api/v1"
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
	"github.com/snehmatic/mindloop/internal/repository/appsettings"
	taskRepo "github.com/snehmatic/mindloop/internal/repository/task"
	"github.com/snehmatic/mindloop/internal/server"
	"github.com/spf13/cobra"
)

var (
	serverPort string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Mindloop web server",
	Run: func(cmd *cobra.Command, args []string) {
		uc := config.UserConfig{}
		_ = uc.ReadFromYAML()
		// Initialize core services
		journalService := journal.NewService(*jRepo)
		noteService := note.NewService(*nRepo)
		pointSvc := points.NewService(*pRepo)
		backupService := backup.NewService(gdb, pointSvc, *fRepo, *hRepo, *hlRepo, *iRepo, *jRepo, *nRepo, *pRepo, *qRepo, *rRepo, *sRepo, *tRepo)
		focusService := focus.NewService(*fRepo)
		intentService := intent.NewService(*iRepo)
		questService := quest.NewService(*qRepo, log.Get())
		summaryService := summary.NewService(*fRepo, *hRepo, *iRepo, *pRepo, *tRepo, log.Get())
		habitService := habit.NewService(*hRepo)

		taskRepoInst := taskRepo.NewSQLTaskRepository(gdb) // TODO make factory with all repositories linked to proper persistence layer defined in config file e.g. SQL
		taskService := task.NewService(
			taskRepoInst,
			&uc,
			pointSvc,
			log.Get(),
		)

		appSettingsRepo := appsettings.NewSQLRepository(gdb)

		mlh := v1.NewMindloopHandler(
			gdb,
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

		server.Serve(mlh, serverPort)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "8765", "Port to run the server on")
	rootCmd.AddCommand(serverCmd)
}
