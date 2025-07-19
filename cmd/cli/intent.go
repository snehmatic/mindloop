package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// parent intent command
var intentCmd = &cobra.Command{
	Use:     "intent",
	Short:   "Manage your intents",
	Example: `mindloop intent start "Get this work done"`,
}

// start intent subcommand
var intentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new intent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting new intent...")
	},
	Example: `mindloop intent start "Get this work done"`,
}

// pause intent subcommand
var intentPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause a running intent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pausing intent...")
	},
	Example: `mindloop intent pause "Get this work done"`,
}

// finish intent subcommand
var intentFinishCmd = &cobra.Command{
	Use:     "finish",
	Short:   "Finish intent",
	Example: `mindloop intent finish "Get this work done"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finishing intent...")
	},
}

func init() {
	rootCmd.AddCommand(intentCmd)
	intentCmd.AddCommand(intentStartCmd)
	intentCmd.AddCommand(intentPauseCmd)
	intentCmd.AddCommand(intentFinishCmd)
}
