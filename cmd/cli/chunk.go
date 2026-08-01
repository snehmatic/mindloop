package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/snehmatic/mindloop/db"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/ai"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/task"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var chunkCmd = &cobra.Command{
	Use:   "chunk [intent|task] [id]",
	Short: "Decompose an intent or task into actionable micro-steps",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		itemType := args[0]
		idStr := args[1]
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			utils.PrintErrorln("Invalid ID")
			return
		}

		appConfig := config.GetConfig()
		dbConn, err := db.ConnectToDb(*appConfig)
		if err != nil {
			utils.PrintErrorln("Failed to connect to database")
			return
		}

		aiSvc := ai.NewService(dbConn)
		taskSvc := task.NewService(dbConn)

		var itemName string
		switch itemType {
		case "intent":
			intentSvc := intent.NewService(dbConn)
			i, err := intentSvc.GetIntent(idStr)
			if err != nil || i == nil {
				utils.PrintErrorln("Intent not found")
				return
			}
			itemName = i.Name
		case "task":
			t, err := taskSvc.GetTask(uint(id))
			if err != nil || t == nil {
				utils.PrintErrorln("Task not found")
				return
			}
			itemName = t.Title
		default:
			utils.PrintErrorln("Type must be intent or task")
			return
		}

		utils.PrintInfoln(fmt.Sprintf("Chunking %s: %s...", itemType, itemName))

		res, err := aiSvc.GenerateChunker(itemName)
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("AI error: %v", err))
			return
		}

		var steps []string

		// Clean the response if it has markdown tags
		cleanRes := strings.TrimSpace(res)
		if strings.HasPrefix(cleanRes, "```json") {
			cleanRes = strings.TrimPrefix(cleanRes, "```json")
			cleanRes = strings.TrimSuffix(cleanRes, "```")
		} else if strings.HasPrefix(cleanRes, "```") {
			cleanRes = strings.TrimPrefix(cleanRes, "```")
			cleanRes = strings.TrimSuffix(cleanRes, "```")
		}

		err = json.Unmarshal([]byte(cleanRes), &steps)
		if err != nil {
			utils.PrintErrorln(fmt.Sprintf("Failed to parse AI response: %v\nResponse:\n%s", err, res))
			return
		}

		utils.PrintSuccessln(fmt.Sprintf("Generated %d micro-steps:", len(steps)))

		for _, step := range steps {
			switch itemType {
			case "intent":
				uid := uint(id)
				_, err := taskSvc.CreateTask(step, &uid, nil)
				if err != nil {
					utils.PrintErrorln(fmt.Sprintf("Failed to create task: %v", err))
				} else {
					fmt.Printf("  - [Task] %s\n", step)
				}
			case "task":
				_, err := taskSvc.AddSubTask(uint(id), step)
				if err != nil {
					utils.PrintErrorln(fmt.Sprintf("Failed to create subtask: %v", err))
				} else {
					fmt.Printf("  - [SubTask] %s\n", step)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chunkCmd)
}
