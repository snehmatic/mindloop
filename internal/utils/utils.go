package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

func PrettyPrint(x any) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func WriteResponse(data interface{}, respWriter http.ResponseWriter, status int) {
	respWriter.Header().Set("content-type", "application/json; charset=utf-8")
	respWriter.WriteHeader(status)
	json.NewEncoder(respWriter).Encode(data)
}

func GetEnvOrDie(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatal().Str("key", key).Msg("failed to get environment variable")
	return ""
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
func FileWrite(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Error().Err(err).Str("file", filename).Msg("failed to write file")
		return err
	}
	log.Info().Str("file", filename).Msg("file written successfully")
	return nil
}
func FileRead(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Error().Err(err).Str("file", filename).Msg("failed to read file")
		return nil, err
	}
	log.Info().Str("file", filename).Msg("file read successfully")
	return data, nil
}
func FileDelete(filename string) error {
	if err := os.Remove(filename); err != nil {
		log.Error().Err(err).Str("file", filename).Msg("failed to delete file")
		return err
	}
	log.Info().Str("file", filename).Msg("file deleted successfully")
	return nil
}

func ValidateUserConfig(cmd *cobra.Command) {
	// check if user_config.yaml exists
	if FileExists(models.UserConfigPath) {
		fmt.Println("User config exists at", models.UserConfigPath)
	} else {
		if cmd.Use != "configure" {
			fmt.Println("Warn: user config does not exist, create a new one or run `mindloop configure`.")
			os.Exit(0)
		}
	}
}
