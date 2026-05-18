package subtask

import (
	"gorm.io/gorm"

	"github.com/snehmatic/mindloop/models"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based subtask repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// FindSubTasks retrieves all subtasks from the database
func (r *sqlRepository) FindSubTasks() ([]models.SubTask, error) {
	var subTasks []models.SubTask
	result := r.DB.Order("CreatedAt DESC").Find(&subTasks)
	return subTasks, result.Error
}

// CreateSubTasks creates multiple subtasks in the database
func (r *sqlRepository) CreateSubTasks(subTasks []models.SubTask) error {
	return r.DB.Create(&subTasks).Error
}
