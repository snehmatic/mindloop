package cli

import (
	"strconv"
	"strings"

	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/quest"
	. "github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var questService *quest.Service

var questCmd = &cobra.Command{
	Use:     "quest",
	Short:   "Manage your side quests",
	Long:    `Side quests are ad-hoc tasks that interrupt your main flow.`,
	Example: `mindloop quest start "Fix prod issue"`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		questService = quest.NewService(gdb)
		intentService = intent.NewService(gdb)
		focusService = focus.NewService(gdb)
	},
}

var questStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new side quest",
	Long:    `Start a new side quest. This will pause any active intent or focus session.`,
	Example: `mindloop quest start "Investigate outage"`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]

		// 1. Check/Pause active intent
		activeIntents, err := intentService.ListActiveIntents()
		if err != nil {
			PrintErrorln("Error checking active intents:", err)
			return
		}
		for _, i := range activeIntents {
			_, err := intentService.PauseIntent(i.ID)
			if err != nil {
				PrintErrorln("Error pausing intent:", err)
				return
			}
			PrintInfof("Paused intent: %s\n", i.Name)
		}

		// 2. Check/Pause active focus session
		activeSession, err := focusService.GetActiveSession()
		if err == nil && activeSession != nil {
			_, err := focusService.PauseSession(activeSession.ID)
			if err != nil {
				PrintErrorln("Error pausing focus session:", err)
				return
			}
			PrintInfof("Paused focus session: %s\n", activeSession.Title)
		}

		// 3. Start Side Quest
		q, err := questService.StartQuest(title)
		if err != nil {
			PrintErrorln("Error starting quest:", err)
			return
		}

		PrintSuccessf("Side Quest '%s' started! Main quest paused.\n", q.Title)
		PrintRocketln("Adventure awaits!")
	},
}

var questStopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop the active side quest",
	Long:    `Stop the active side quest. You can provide a note as arguments, or open an editor.`,
	Example: `mindloop quest stop "Fixed the issue"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get active quest
		q, err := questService.GetActiveQuest()
		if err != nil {
			PrintErrorln("Error checking active side quest:", err)
			return
		}
		if q == nil {
			PrintErrorln("No active side quest found.")
			return
		}

		var note string
		if len(args) > 0 {
			note = strings.Join(args, " ")
		} else {
			// Interactive mode
			header := "# Quest Summary\n# Describe what you accomplished:\n\n"
			note, err = CaptureWithEditor("mindloop_quest_*.md", header, "")
			if err != nil {
				PrintErrorln("Error capturing note:", err)
				return
			}
		}

		if note == "" {
			PrintWarnln("No note provided. Saving with empty note.")
		}

		q, err = questService.StopQuest(q.ID, note)
		if err != nil {
			PrintErrorln("Error stopping quest:", err)
			return
		}

		PrintSuccessf("Side Quest '%s' complete!\n", q.Title)
		PrintInfoln("Don't forget to resume your main intent/focus if needed.")
	},
}

var questDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a side quest",
	Example: `mindloop quest delete <id>`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			PrintErrorln("Invalid ID:", idStr)
			return
		}

		err = questService.DeleteQuest(uint(id))
		if err != nil {
			PrintErrorln("Error deleting quest:", err)
			return
		}

		PrintSuccessf("Quest %d deleted successfully.\n", id)
	},
}

var questListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all side quests",
	Run: func(cmd *cobra.Command, args []string) {
		quests, err := questService.ListQuests()
		if err != nil {
			PrintErrorln("Error listing quests:", err)
			return
		}

		if len(quests) == 0 {
			PrintInfoln("No quests found.")
			return
		}

		var views []models.SideQuestView
		for _, q := range quests {
			views = append(views, models.ToSideQuestView(q))
		}
		PrintTable(views)
	},
}

func init() {
	questCmd.AddCommand(questStartCmd)
	questCmd.AddCommand(questStopCmd)
	questCmd.AddCommand(questDeleteCmd)
	questCmd.AddCommand(questListCmd)
	rootCmd.AddCommand(questCmd)
}