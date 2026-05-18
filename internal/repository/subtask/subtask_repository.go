package subtask

import (
	"github.com/snehmatic/mindloop/models"
)

// Repository defines the interface for subtask data access
type Repository interface {
	FindSubTasks() ([]models.SubTask, error)
	CreateSubTasks(subTasks []models.SubTask) error
}
