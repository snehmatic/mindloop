package cli

import (
	"strconv"

	"github.com/snehmatic/mindloop/internal/core/focus"
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	focusService *focus.Service
)

var focusCmd = &cobra.Command{
	Use:     "focus",
	Short:   "Manage your focus sessions",
	Long:    `Focus sessions help you track your work and productivity.`,
	Example: `mindloop focus start "Work on project"`,
	Args:    cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		focusService = focus.NewService(gdb)
	},
}

var focusStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new focus session",
	Long:    `Start a new focus session to track your work.`,
	Example: `mindloop focus start "Work on project"`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		PrintRocketln("That's the spirit! Starting a new focus session...")
		session, err := focusService.StartSession(args[0])
		if err != nil {
			PrintErrorln("Error starting focus session:", err)
			ac.Logger.Error().Msgf("Error starting focus session: %v", err)
			return
		}
		PrintSuccessf("Focus session '%s' started successfully with id %d!\n", session.Title, session.ID)
		ac.Logger.Info().Msgf("Focus session '%s' started successfully with id %d!", session.Title, session.ID)
	},
}

var focusListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all focus sessions",
	Long:    `List all your focus sessions to review your productivity.`,
	Example: `mindloop focus list`,
	Run: func(cmd *cobra.Command, args []string) {
		sessions, err := focusService.ListSessions()
		if err != nil {
			PrintErrorln("Error listing focus sessions:", err)
			ac.Logger.Error().Msgf("Error listing focus sessions: %v", err)
			return
		}
		if len(sessions) == 0 {
			PrintInfoln("No focus sessions found... Try starting one with 'mindloop focus start <title>'")
			ac.Logger.Info().Msg("No focus sessions found. Prompting user to start a new focus session.")
			return
		}

		var views []models.FocusSessionView
		for _, session := range sessions {
			views = append(views, models.ToFocusSessionView(session))
		}

		ac.Logger.Info().Msg("Listing all focus sessions.")
		PrintInfoln("Focus sessions listed below. Note: Duration is in minutes")
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
			PrintErrorln("Error parsing session ID:", err)
			return
		}

		session, err := focusService.EndSession(sessionIDInt)
		if err != nil {
			PrintErrorln("Error ending focus session:", err)
			ac.Logger.Error().Msgf("Error ending focus session: %v", err)
			return
		}

		PrintSuccessf("Focus session '%s' ended successfully!\n", session.Title)
		PrintRocketln("Great work chief!")
		ac.Logger.Info().Msgf("Focus session '%s' ended successfully!", session.Title)
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
		sessionIDInt, err := strconv.Atoi(sessionID)
		if err != nil {
			PrintErrorln("Error parsing session ID:", err)
			return
		}

		rating, err := strconv.Atoi(args[1])
		if err != nil {
			PrintWarnln("Rating must be an integer.")
			return
		}

		session, err := focusService.RateSession(sessionIDInt, rating)
		if err != nil {
			PrintErrorln("Error saving rating:", err)
			ac.Logger.Error().Msgf("Error saving rating for focus session: %v", err)
			return
		}

		PrintSuccessf("'%s' session rated successfully with a score of %d!\n", session.Title, session.Rating)
		ac.Logger.Info().Msgf("Focus session '%s' rated successfully with a score of %d!", session.Title, session.Rating)
	},
}

func init() {
	focusCmd.AddCommand(focusStartCmd)
	focusCmd.AddCommand(focusListCmd)
	focusCmd.AddCommand(focusEndCmd)
	focusCmd.AddCommand(focusRateCmd)

	rootCmd.AddCommand(focusCmd)
}
