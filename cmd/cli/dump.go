package cli

import (
	"github.com/snehmatic/mindloop/internal/core/dump"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:     "dump <text>",
	Short:   "Capture a quick thought (BrainDump)",
	Example: `mindloop dump "Remember to buy milk"`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dumpService := dump.NewService(gdb)
		text := args[0]
		// join args if they passed multiple words without quotes
		if len(args) > 1 {
			for i := 1; i < len(args); i++ {
				text += " " + args[i]
			}
		}
		bd, err := dumpService.CreateDump(text)
		if err != nil {
			utils.PrintErrorln("Error saving dump:", err)
			return
		}
		utils.PrintSuccessf("Captured dump #%d successfully!\n", bd.ID)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
