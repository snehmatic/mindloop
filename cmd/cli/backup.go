package cli

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/core/backup"
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var (
	backupService *backup.Service
)

var backupCmd = &cobra.Command{
	Use:     "backup",
	Short:   "Backup your Mindloop data to a JSON file",
	Example: `mindloop backup mindloop_backup.json`,
	Args:    cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		backupService = backup.NewService(gdb, points.NewService(*pRepo), *fRepo, *hRepo, *hlRepo, *iRepo, *jRepo, *nRepo, *pRepo, *qRepo, *rRepo, *sRepo, *tRepo)
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		utils.PrintRocketf("Exporting data to %s...\n", filePath)
		if err := backupService.Export(filePath); err != nil {
			utils.PrintErrorln("Backup failed:", err)
			return
		}
		utils.PrintSuccessln("Backup completed successfully!")
	},
}

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Restore your Mindloop data from a JSON file",
	Long:    "WARNING: This will replace all existing local data with the data from the backup file.",
	Example: `mindloop restore mindloop_backup.json`,
	Args:    cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		backupService = backup.NewService(gdb, points.NewService(*pRepo), *fRepo, *hRepo, *hlRepo, *iRepo, *jRepo, *nRepo, *pRepo, *qRepo, *rRepo, *sRepo, *tRepo)
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		utils.PrintWarnf("Restoring data from %s. All existing local data will be replaced!\n", filePath)
		utils.PrintInfoln("Type 'yes' to confirm:")
		var confirm string
		_, _ = fmt.Scanln(&confirm)
		if confirm != "yes" {
			utils.PrintInfoln("Restore cancelled.")
			return
		}

		utils.PrintRocketln("Importing data...")
		if err := backupService.Import(filePath); err != nil {
			utils.PrintErrorln("Restore failed:", err)
			return
		}
		utils.PrintSuccessln("Restore completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
}
