package task

import (
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Repository defines the data access methods for Tasks and SubTasks
type Repository interface {
	CreateTask(t *models.Task) error
	GetTask(id uint) (*models.Task, error)
	UpdateTask(t *models.Task) error
	ListTasks() ([]models.Task, error)
	GetTasksByIntent(intentID uint) ([]models.Task, error)
	GetTasksByFocusSession(focusID uint) ([]models.Task, error)
	DeleteTask(id uint) error
	DeleteSubTasksByTaskID(taskID uint) error
	ReorderTasks(ids []uint) error
	RecalibrateTasks(status string, before time.Time) error

	CreateSubTask(st *models.SubTask) error
	GetSubTask(id uint) (*models.SubTask, error)
	UpdateSubTask(st *models.SubTask) error
	DeleteSubTask(id uint) error
	ReorderSubTasks(ids []uint) error

	GetDB() *gorm.DB // Escape hatch for cross-module transactions (like points)
}
