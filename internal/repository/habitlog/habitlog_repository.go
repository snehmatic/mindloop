package habitlog

import (
	"github.com/snehmatic/mindloop/models"
)

// Repository defines the interface for habit log data access
type Repository interface {
	FindHabitLogs() ([]models.HabitLog, error)
	CreateHabitLogs(logs []models.HabitLog) error
}
