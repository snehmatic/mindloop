package cli

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var focusCmd = &cobra.Command{
	Use:     "focus",
	Short:   "Manage your focus sessions",
	Long:    `Focus sessions help you track your work and productivity.`,
	Example: `mindloop focus start "Work on project"`,
	Args:    cobra.NoArgs,
}

var focusStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new focus session",
	Long:    `Start a new focus session to track your work.`,
	Example: `mindloop focus start "Work on project"`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		focusSession := &models.FocusSession{
			Title:  args[0],
			Status: "active",
		}
		if err := gdb.Create(focusSession).Error; err != nil {
			cmd.Println("Error starting focus session:", err)
			ac.Logger.Error().Msgf("Error starting focus session: %v", err)
			return
		}
		cmd.Printf("Focus session '%s' started successfully with id %d!\n", focusSession.Title, focusSession.ID)
		ac.Logger.Info().Msgf("Focus session '%s' started successfully with id %d!", focusSession.Title, focusSession.ID)
	},
}

var focusListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all focus sessions",
	Long:    `List all your focus sessions to review your productivity.`,
	Example: `mindloop focus list`,
	Run: func(cmd *cobra.Command, args []string) {
		var focusSessions []models.FocusSession
		if err := gdb.Find(&focusSessions).Error; err != nil {
			cmd.Println("Error listing focus sessions:", err)
			ac.Logger.Error().Msgf("Error listing focus sessions: %v", err)
			return
		}
		if len(focusSessions) == 0 {
			cmd.Println("No focus sessions found... Try starting one with 'mindloop focus start <title>'")
			ac.Logger.Info().Msg("No focus sessions found. Prompting user to start a new focus session.")
			return
		}

		var views []models.FocusSessionView
		for _, session := range focusSessions {
			views = append(views, models.ToFocusSessionView(session))
		}

		ac.Logger.Info().Msg("Listing all focus sessions.")
		fmt.Println("Focus sessions listed below. Note: Duration is in minutes")
		PrintTable(views)
	},
}

var focusEndCmd = &cobra.Command{
	Use:     "end",
	Short:   "End a focus session",
	Long:    `End an active focus session to mark it as completed.`,
	Example: `mindloop focus end <session_id>`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := args[0]
		sessionIDInt, err := strconv.Atoi(sessionID)
		if err != nil {
			cmd.Println("Error parsing session ID:", err)
			ac.Logger.Error().Msgf("Error parsing session ID: %v", err)
			return
		}

		var focusSession models.FocusSession
		if err := gdb.First(&focusSession, sessionIDInt).Error; err != nil {
			cmd.Println("Focus session not found:", err)
			ac.Logger.Error().Msgf("Focus session not found: %v", err)
			return
		}

		if focusSession.Status != "active" {
			cmd.Println("Focus session is not active.")
			ac.Logger.Warn().Msg("Attempted to end a non-active focus session.")
			return
		}

		focusSession.Status = "ended"
		focusSession.EndTime = time.Now()
		focusSession.Duration = focusSession.EndTime.Sub(focusSession.CreatedAt).Minutes()
		if err := gdb.Save(&focusSession).Error; err != nil {
			cmd.Println("Error ending focus session:", err)
			ac.Logger.Error().Msgf("Error ending focus session: %v", err)
			return
		}
		cmd.Printf("Focus session '%s' ended successfully!\n", focusSession.Title)
		ac.Logger.Info().Msgf("Focus session '%s' ended successfully!", focusSession.Title)
	},
}

var focusRateCmd = &cobra.Command{
	Use:     "rate",
	Short:   "Rate a focus session",
	Long:    `Rate a completed focus session to provide feedback on your productivity.`,
	Example: `mindloop focus rate <session_id> <rating 0-10>`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := args[0]
		rating, err := strconv.Atoi(args[1])
		if err != nil || rating < 0 || rating > 10 {
			cmd.Println("Rating must be an integer between 0 and 10.")
			ac.Logger.Warn().Msgf("Invalid rating provided: %s", args[1])
			return
		}

		var focusSession models.FocusSession
		if err := gdb.First(&focusSession, sessionID).Error; err != nil {
			cmd.Println("Focus session not found:", err)
			ac.Logger.Error().Msgf("Focus session not found: %v", err)
			return
		}

		if focusSession.Status != "ended" {
			PrintSuccessln("Focus session is not ended. Please end it before rating.")
			ac.Logger.Warn().Msg("Attempted to rate a non-ended focus session.")
			return
		}

		focusSession.Rating = rating
		if err := gdb.Save(&focusSession).Error; err != nil {
			cmd.Println("Error saving rating:", err)
			ac.Logger.Error().Msgf("Error saving rating for focus session ID %d: %v", focusSession.ID, err)
			return
		}
		PrintSuccessf("'%s' session rated successfully with a score of %d!\n", focusSession.Title, rating)
		ac.Logger.Info().Msgf("Focus session '%s' rated successfully with a score of %d!", focusSession.Title, rating)
	},
}

func init() {
	focusCmd.AddCommand(focusStartCmd)
	focusCmd.AddCommand(focusListCmd)
	focusCmd.AddCommand(focusEndCmd)
	focusCmd.AddCommand(focusRateCmd)

	rootCmd.AddCommand(focusCmd)
}
