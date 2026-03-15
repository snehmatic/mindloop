package cli

import (
	"fmt"
	"strconv"

	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	noteService *note.Service
	noteTitle   string
	noteLabels  string
)

var noteCmd = &cobra.Command{
	Use:     "note",
	Short:   "Capture and manage simple notes",
	Long:    `Notes help you quickly capture ideas, snippets, and important information.`,
	Example: `mindloop note "Meeting at 3 PM"`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		noteService = note.NewService(gdb)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// Create a quick note
			content := args[0]
			n, err := noteService.CreateNote(noteTitle, content, noteLabels)
			if err != nil {
				utils.PrintErrorln("Error creating note:", err)
				return
			}
			utils.PrintSuccessf("Note saved with ID %d!\n", n.ID)
			return
		}
		// If no args, list notes
		runNoteList(cmd, args)
	},
}

var noteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all notes",
	Run:     runNoteList,
}

func runNoteList(cmd *cobra.Command, args []string) {
	notes, err := noteService.ListNotes()
	if err != nil {
		utils.PrintErrorln("Error listing notes:", err)
		return
	}
	if len(notes) == 0 {
		utils.PrintInfoln("No notes found. Try 'mindloop note \"your note text\"'")
		return
	}

	var views []models.NoteView
	for _, n := range notes {
		views = append(views, models.ToNoteView(n))
	}
	utils.PrintTable(views)
}

var noteNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note using your default $EDITOR",
	Run: func(cmd *cobra.Command, args []string) {
		header := "# Mindloop Note\n# Title: " + noteTitle + "\n# Labels: " + noteLabels + "\n# Write your note below. Lines starting with # will be ignored.\n\n"
		content, err := utils.CaptureWithEditor("mindloop_note_*.md", header, "")
		if err != nil {
			utils.PrintErrorln("Error capturing note:", err)
			return
		}
		if content == "" {
			utils.PrintWarnln("Empty note. Nothing saved.")
			return
		}

		n, err := noteService.CreateNote(noteTitle, content, noteLabels)
		if err != nil {
			utils.PrintErrorln("Failed to save note:", err)
			return
		}
		utils.PrintSuccessf("Note '%s' saved with ID %d!\n", n.Title, n.ID)
	},
}

var noteViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a note",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			utils.PrintErrorln("Invalid ID:", args[0])
			return
		}
		n, err := noteService.GetNote(id)
		if err != nil {
			utils.PrintErrorln("Note not found:", err)
			return
		}
		fmt.Println("-------------------------------")
		utils.PrintInfoln("Title:", n.Title)
		utils.PrintInfoln("Labels:", n.Labels)
		utils.PrintLoadingln("Updated:", n.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("-------------------------------")
		fmt.Println(n.Content)
		fmt.Println("-------------------------------")
	},
}

var noteEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a note",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			utils.PrintErrorln("Invalid ID:", args[0])
			return
		}
		n, err := noteService.GetNote(id)
		if err != nil {
			utils.PrintErrorln("Note not found:", err)
			return
		}

		header := "# Mindloop Note Edit\n# Lines starting with # will be ignored.\n\n"
		content, err := utils.CaptureWithEditor("mindloop_note_edit_*.md", header, n.Content)
		if err != nil {
			utils.PrintErrorln("Error capturing note:", err)
			return
		}

		_, err = noteService.UpdateNote(id, n.Title, content, n.Labels)
		if err != nil {
			utils.PrintErrorln("Failed to update note:", err)
			return
		}
		utils.PrintSuccessln("Note updated.")
	},
}

var noteDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a note",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			utils.PrintErrorln("Invalid ID:", args[0])
			return
		}
		err = noteService.DeleteNote(id)
		if err != nil {
			utils.PrintErrorln("Error deleting note:", err)
			return
		}
		utils.PrintSuccessln("Note deleted.")
	},
}

func init() {
	noteCmd.Flags().StringVarP(&noteTitle, "title", "t", "", "Note title")
	noteCmd.Flags().StringVarP(&noteLabels, "labels", "l", "", "Comma separated labels")

	noteNewCmd.Flags().StringVarP(&noteTitle, "title", "t", "", "Note title")
	noteNewCmd.Flags().StringVarP(&noteLabels, "labels", "l", "", "Comma separated labels")

	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteNewCmd)
	noteCmd.AddCommand(noteViewCmd)
	noteCmd.AddCommand(noteEditCmd)
	noteCmd.AddCommand(noteDeleteCmd)

	rootCmd.AddCommand(noteCmd)
}
