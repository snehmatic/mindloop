package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// parent habit command
var habitCmd = &cobra.Command{
	Use:     "habit",
	Short:   "Manage your habits",
	Example: `mindloop habit add "Excercise" --daily --time "12:00 PM" --tags health,fitness`,
}

// add habit subcommand
var habitAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new habit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Adding new habit...")
		fmt.Println(args)
	},
	Example:    `mindloop habit add "Excercise" --daily --time "12:00 PM" --labels health,fitness`,
	Aliases:    []string{"create", "new"},
	ValidArgs:  []string{"daily", "time", "labels"},
	ArgAliases: []string{"d", "t", "l"},
}

// delete habit subcommand
var habitDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a habit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting habit...")
	},
	Example: `mindloop habit delete "Excercise"`,
}

// update habit subcommand
var habitUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a habit",
	Example: `mindloop habit update "Excercise" --time "1:00 PM"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updated habit...")
	},
	Aliases: []string{"edit", "modify"},
}

// mark habit as done subcommand
var habitMarkAsDoneCmd = &cobra.Command{
	Use:     "done",
	Short:   "Mark a habit as done",
	Example: `mindloop habit done "Excercise"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Marked habit as done...")
	},
}

// mark habit as done subcommand
var habitListCmd = &cobra.Command{
	Use:                   "list",
	Short:                 "List all habits",
	Example:               `mindloop habit list`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"all", "active", "completed"},
	Args:                  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all habits...")
		fmt.Println("Args:", args)
		fmt.Printf("%+v\n", cmd.Flags())
	},
}

func init() {
	// cmds
	rootCmd.AddCommand(habitCmd)
	habitCmd.AddCommand(habitAddCmd)
	habitCmd.AddCommand(habitDeleteCmd)
	habitCmd.AddCommand(habitUpdateCmd)
	habitCmd.AddCommand(habitMarkAsDoneCmd)
	habitCmd.AddCommand(habitListCmd)

	// flags
	habitCmd.PersistentFlags().BoolP("all", "A", false, "Select all habits")
	habitAddCmd.Flags().BoolP("daily", "d", false, "Set habit as daily")
	habitAddCmd.Flags().BoolP("weekly", "w", false, "Set habit as weekly")
	habitAddCmd.Flags().StringP("time", "t", "", "Set reminder time for the habit")
	habitAddCmd.Flags().StringSliceP("labels", "l", []string{}, "habit labels (comma separated)")
}
