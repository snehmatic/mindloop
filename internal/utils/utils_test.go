package utils_test

import (
	"testing"

	"github.com/snehmatic/mindloop/internal/utils"
)

func TestFormatMinutes(t *testing.T) {
	tests := []struct {
		minutes float64
		want    string
	}{
		{0, "0min"},
		{30, "30min"},
		{60, "1hr"},
		{90, "1hr 30min"},
		{120, "2hr"},
		{125, "2hr 5min"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := utils.FormatMinutes(tt.minutes)
			if got != tt.want {
				t.Errorf("FormatMinutes(%f) = %v, want %v", tt.minutes, got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	if utils.FileExists("non_existent_file_12345") {
		t.Error("FileExists returned true for non-existent file")
	}
}
