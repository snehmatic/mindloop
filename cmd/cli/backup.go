package cli

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/core/backup"
	. "github.com/snehmatic/mindloop/internal/utils"
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
		backupService = backup.NewService(gdb)
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		PrintRocketf("Exporting data to %s...\n", filePath)
		if err := backupService.Export(filePath); err != nil {
			PrintErrorln("Backup failed:", err)
			return
		}
		PrintSuccessln("Backup completed successfully!")
	},
}

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Restore your Mindloop data from a JSON file",
	Long:    "WARNING: This will replace all existing local data with the data from the backup file.",
	Example: `mindloop restore mindloop_backup.json`,
	Args:    cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		backupService = backup.NewService(gdb)
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		PrintWarnf("Restoring data from %s. All existing local data will be replaced!\n", filePath)
		PrintInfoln("Type 'yes' to confirm:")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			PrintInfoln("Restore cancelled.")
			return
		}

		PrintRocketln("Importing data...")
		if err := backupService.Import(filePath); err != nil {
			PrintErrorln("Restore failed:", err)
			return
		}
		PrintSuccessln("Restore completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
}