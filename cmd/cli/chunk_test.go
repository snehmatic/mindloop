package cli

import (
	"testing"
)

func TestChunkCmd_InvalidArgs(t *testing.T) {
	// Directly invoke the Run function to avoid rootCmd init panic
	if chunkCmd.Run != nil {
		chunkCmd.Run(chunkCmd, []string{"invalid", "not-a-number"})
	}
}
