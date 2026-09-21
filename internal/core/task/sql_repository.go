package task

import (
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *SQLRepository) CreateTask(t *models.Task) error {
	return r.db.Create(t).Error
}

func (r *SQLRepository) GetTask(id uint) (*models.Task, error) {
	var t models.Task
	if err := r.db.Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *SQLRepository) UpdateTask(t *models.Task) error {
	return r.db.Save(t).Error
}

func (r *SQLRepository) ListTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *SQLRepository) GetTasksByIntent(intentID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("IntentID = ?", intentID).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *SQLRepository) GetTasksByFocusSession(focusID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("FocusSessionID = ?", focusID).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *SQLRepository) DeleteTask(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *SQLRepository) DeleteSubTasksByTaskID(taskID uint) error {
	return r.db.Where("TaskID = ?", taskID).Delete(&models.SubTask{}).Error
}

func (r *SQLRepository) ReorderTasks(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.Task{}).Where("id = ?", id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *SQLRepository) RecalibrateTasks(status string, before time.Time) error {
	return r.db.Model(&models.Task{}).Where("Status = ? AND DueDate < ?", status, before).Update("DueDate", nil).Error
}

func (r *SQLRepository) CreateSubTask(st *models.SubTask) error {
	return r.db.Create(st).Error
}

func (r *SQLRepository) GetSubTask(id uint) (*models.SubTask, error) {
	var st models.SubTask
	if err := r.db.First(&st, id).Error; err != nil {
		return nil, err
	}
	return &st, nil
}

func (r *SQLRepository) UpdateSubTask(st *models.SubTask) error {
	return r.db.Save(st).Error
}

func (r *SQLRepository) DeleteSubTask(id uint) error {
	return r.db.Delete(&models.SubTask{}, id).Error
}

func (r *SQLRepository) ReorderSubTasks(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.SubTask{}).Where("id = ?", id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
