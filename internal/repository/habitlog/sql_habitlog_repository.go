package habitlog

import (
	"gorm.io/gorm"

	"github.com/snehmatic/mindloop/models"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based habit log repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// FindHabitLogs retrieves all habit logs from the database
func (r *sqlRepository) FindHabitLogs() ([]models.HabitLog, error) {
	var logs []models.HabitLog
	result := r.DB.Order("CreatedAt DESC").Find(&logs)
	return logs, result.Error
}

// CreateHabitLogs creates multiple habit logs in the database
func (r *sqlRepository) CreateHabitLogs(logs []models.HabitLog) error {
	return r.DB.Create(&logs).Error
}
