package utils

import (
	"fmt"
	"math"
)

// FormatMinutes converts float64 minutes into a human-readable string like "1hr 2min"
func FormatMinutes(minutes float64) string {
	totalMinutes := int(math.Round(minutes))
	hours := totalMinutes / 60
	mins := totalMinutes % 60

	switch {
	case hours > 0 && mins > 0:
		return fmt.Sprintf("%dhr %dmin", hours, mins)
	case hours > 0:
		return fmt.Sprintf("%dhr", hours)
	default:
		return fmt.Sprintf("%dmin", mins)
	}
}
