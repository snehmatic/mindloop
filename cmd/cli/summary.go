package cli

import (
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var summaryCmd = &cobra.Command{
	Use:     "summary",
	Short:   "Get a summary of your productivity",
	Example: `mindloop summary`,
	Run: func(cmd *cobra.Command, args []string) {
		PrettyPrintBanner()
		PrintRocketln("Generating your productivity summary...")

		PrintInfoln("This feature is under development. Stay tuned for updates!")
		ac.Logger.Info().Msg("Requested productivity summary, but feature is not yet implemented.")
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
