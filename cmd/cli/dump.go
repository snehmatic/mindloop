package cli

import (
	"strings"

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
		text := strings.Join(args, " ")
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
