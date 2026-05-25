package gamification

// DefaultMilestoneInterval defines the default point increment for milestone celebrations.
const DefaultMilestoneInterval = 200

// NormalizeMilestoneInterval returns a safe milestone interval value.
func NormalizeMilestoneInterval(interval int) int {
	if interval <= 0 {
		return DefaultMilestoneInterval
	}
	return interval
}
