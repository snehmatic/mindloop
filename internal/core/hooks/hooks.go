package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// getHooksDir returns the path to the mindloop hooks directory.
func getHooksDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".mindloop", "hooks"), nil
}

// ExecuteHook asynchronously executes the hook for the given event name if it exists and is executable.
func ExecuteHook(eventName string, envContext map[string]string) {
	hooksDir, err := getHooksDir()
	if err != nil {
		log.Error().Err(err).Msg("failed to get hooks directory")
		return
	}

	hookPath := filepath.Join(hooksDir, eventName)

	// Check if file exists
	info, err := os.Stat(hookPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Error().Err(err).Str("hook", hookPath).Msg("failed to stat hook file")
		}
		// File does not exist, nothing to do
		return
	}

	// Check if executable (owner, group, or other)
	if info.Mode()&0111 == 0 {
		log.Warn().Str("hook", hookPath).Msg("hook file exists but is not executable")
		return
	}

	// Execute asynchronously
	go func() {
		log.Debug().Str("hook", hookPath).Msg("executing hook")
		cmd := exec.Command(hookPath)

		env := os.Environ()
		env = append(env, fmt.Sprintf("MINDLOOP_EVENT=%s", eventName))
		for k, v := range envContext {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Error().Err(err).Str("hook", hookPath).Bytes("output", output).Msg("hook execution failed")
		} else {
			log.Debug().Str("hook", hookPath).Bytes("output", output).Msg("hook executed successfully")
		}
	}()
}
