package cli

import (
	"fmt"
	"strconv"
	"time"

	cfg "github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	focusService      *focus.Service
	focusStatusFormat string
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
		utils.PrintRocketln("That's the spirit! Starting a new focus session...")
		session, err := focusService.StartSession(args[0])
		if err != nil {
			utils.PrintErrorln("Error starting focus session:", err)
			ac.Logger.Error().Msgf("Error starting focus session: %v", err)
			return
		}
		utils.PrintSuccessf("Focus session '%s' started successfully with id %d!\n", session.Title, session.ID)
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
			utils.PrintErrorln("Error listing focus sessions:", err)
			ac.Logger.Error().Msgf("Error listing focus sessions: %v", err)
			return
		}
		if len(sessions) == 0 {
			utils.PrintInfoln("No focus sessions found... Try starting one with 'mindloop focus start <title>'")
			ac.Logger.Info().Msg("No focus sessions found. Prompting user to start a new focus session.")
			return
		}

		var views []models.FocusSessionView
		for _, session := range sessions {
			views = append(views, models.ToFocusSessionView(session))
		}

		ac.Logger.Info().Msg("Listing all focus sessions.")
		utils.PrintInfoln("Focus sessions listed below. Note: Duration is in minutes")
		utils.PrintTable(views)
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
			utils.PrintErrorln("Error parsing session ID:", err)
			return
		}

		uc := cfg.UserConfig{}
		_ = uc.ReadFromYAML()
		session, milestoneReached, err := focusService.EndSession(sessionIDInt, uc.PointsConfig.Focus)
		if err != nil {
			utils.PrintErrorln("Error ending focus session:", err)
			ac.Logger.Error().Msgf("Error ending focus session: %v", err)
			return
		}

		utils.PrintSuccessf("Focus session '%s' ended successfully! (+%d pts) 🎉\n", session.Title, uc.PointsConfig.Focus)
		if milestoneReached {
			utils.PrintRocketln("🏆 MILESTONE REACHED! You're on fire! 🏆")
		}
		utils.PrintRocketln("Great work chief!")
		ac.Logger.Info().Msgf("Focus session '%s' ended successfully!", session.Title)
	},
}

var focusRateCmd = &cobra.Command{Use: "rate",
	Short:   "Rate a focus session",
	Long:    `Rate a completed focus session to provide feedback on your productivity.`,
	Example: `mindloop focus rate <session_id> <rating 0-10>`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := args[0]
		sessionIDInt, err := strconv.Atoi(sessionID)
		if err != nil {
			utils.PrintErrorln("Error parsing session ID:", err)
			return
		}

		rating, err := strconv.Atoi(args[1])
		if err != nil {
			utils.PrintWarnln("Rating must be an integer.")
			return
		}

		session, err := focusService.RateSession(sessionIDInt, rating)
		if err != nil {
			utils.PrintErrorln("Error saving rating:", err)
			ac.Logger.Error().Msgf("Error saving rating for focus session: %v", err)
			return
		}

		utils.PrintSuccessf("'%s' session rated successfully with a score of %d!\n", session.Title, session.Rating)
		ac.Logger.Info().Msgf("Focus session '%s' rated successfully with a score of %d!", session.Title, session.Rating)
	},
}

var focusStatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Get the active focus session status",
	Long:    `Show the currently active focus session, if any.`,
	Example: `mindloop focus status --format=compact`,
	Run: func(cmd *cobra.Command, args []string) {
		session, err := focusService.GetActiveSession()
		if err != nil {
			if focusStatusFormat == "compact" {
				return
			}
			utils.PrintErrorln("Error getting focus status:", err)
			return
		}
		if session == nil {
			if focusStatusFormat == "compact" {
				return
			}
			utils.PrintInfoln("No active focus session.")
			return
		}

		duration := time.Since(session.CreatedAt)
		mins := int(duration.Minutes())

		if focusStatusFormat == "compact" {
			fmt.Printf("⚡ %dm - %s\n", mins, session.Title)
		} else {
			utils.PrintSuccessf("Active Focus: '%s' (%d minutes)\n", session.Title, mins)
		}
	},
}

func init() {
	focusCmd.AddCommand(focusStartCmd)
	focusCmd.AddCommand(focusListCmd)
	focusCmd.AddCommand(focusEndCmd)
	focusCmd.AddCommand(focusRateCmd)

	focusStatusCmd.Flags().StringVar(&focusStatusFormat, "format", "", "Output format (e.g., compact)")
	focusCmd.AddCommand(focusStatusCmd)

	rootCmd.AddCommand(focusCmd)
}
