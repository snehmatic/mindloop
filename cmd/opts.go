package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var optsCmd = &cobra.Command{
	Use:   "opts",
	Short: "Show available product options",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Show product opts here...")
	},
}

func init() {
	rootCmd.AddCommand(optsCmd)
}
