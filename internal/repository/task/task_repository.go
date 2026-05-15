package task

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// TaskRepository defines the data access operations for tasks
// Note: This interface only contains data access methods, no business logic
// Business logic (points, logger, etc.) stays in the Service layer
type TaskRepository interface {
	CreateTask(title string, intentID, focusID *uint) (*models.Task, error)
	GetTask(id uint) (*models.Task, error)
	CompleteTask(id uint) (*models.Task, error)
	ListTasks() ([]models.Task, error)
	AddSubTask(taskID uint, title string) (*models.SubTask, error)
	GetSubTask(id uint) (*models.SubTask, error)
	CompleteSubTask(id uint) (*models.SubTask, error)
	GetTasksByIntent(intentID uint) ([]models.Task, error)
	GetTasksByFocusSession(focusID uint) ([]models.Task, error)
	DeleteTask(id uint) error
	DeleteSubTask(id uint) error
	ReorderTasks(ids []uint) error
	ReorderSubTasks(ids []uint) error
	UpdateTask(task *models.Task) error
	UpdateSubTask(subTask *models.SubTask) error
	GetDB() *gorm.DB
}
