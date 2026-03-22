package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication settings",
}

var authSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Enable or change local web UI password",
	Run: func(cmd *cobra.Command, args []string) {
		uc := config.UserConfig{}
		if err := uc.ReadFromYAML(); err != nil {
			utils.PrintErrorln("Failed to read config")
			return
		}

		fmt.Print("Enter new password for web UI: ")
		var password string
		fmt.Scanln(&password)

		if password == "" {
			utils.PrintWarnln("Password cannot be empty")
			return
		}

		hash := sha256.Sum256([]byte(password))
		uc.AuthConfig.PasswordHash = hex.EncodeToString(hash[:])
		uc.AuthConfig.Enabled = true
		uc.WriteToYAML()

		utils.PrintSuccessln("Password set and authentication enabled for Web UI")
	},
}

var authDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable local web UI password protection",
	Run: func(cmd *cobra.Command, args []string) {
		uc := config.UserConfig{}
		if err := uc.ReadFromYAML(); err != nil {
			utils.PrintErrorln("Failed to read config")
			return
		}

		uc.AuthConfig.Enabled = false
		uc.WriteToYAML()

		utils.PrintSuccessln("Authentication disabled for Web UI")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authSetupCmd)
	authCmd.AddCommand(authDisableCmd)
}
