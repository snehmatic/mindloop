package cli

import (
	"github.com/rs/zerolog"
	v1 "github.com/snehmatic/mindloop/api/v1"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/backup"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/quest"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/core/task"
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
		journalService := journal.NewService(gdb)
		noteService := note.NewService(gdb)
		backupService := backup.NewService(gdb)
		focusService := focus.NewService(gdb)
		intentService := intent.NewService(gdb)
		questService := quest.NewService(gdb)
		summaryService := summary.NewService(gdb)
		habitService := habit.NewService(gdb)

		taskRepo := taskRepo.NewSQLTaskRepository(gdb) // TODO make factory with all repositories linked to proper persistence layer defined in config file e.g. SQL
		taskService := task.NewService(
			taskRepo,
			&uc,
			zerolog.Nop(),
		)

		mlh := v1.NewMindloopHandler(
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
