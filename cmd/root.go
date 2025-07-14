package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mindloop",
	Short: "mindloop is a CLI tool for productivity tracking",
	Long:  `Mindloop helps track intent, focus sessions, and habits via CLI.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// define persistent flags here
}

func initConfig() {
	// Setup config, env vars, etc
}
