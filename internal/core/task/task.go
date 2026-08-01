package task

import (
	"time"
	"errors"

	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/nlp"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

var logger = log.Get()

// Service handles business logic for tasks and sub-tasks
type Service struct {
	db *gorm.DB
	uc *config.UserConfig
}

// NewService creates a new task Service instance
func NewService(db *gorm.DB) *Service {
	uc := config.UserConfig{}
	_ = uc.ReadFromYAML()

	return &Service{
		db: db,
		uc: &uc,
	}
}

// CreateTask persists a new task to the database
func (s *Service) CreateTask(title string, intentID, focusID *uint) (*models.Task, error) {
	cleanedTitle, dueDate := nlp.ExtractDate(title)

	t := &models.Task{
		Title:          cleanedTitle,
		IntentID:       intentID,
		FocusSessionID: focusID,
		DueDate:        dueDate,
	}
	if err := s.db.Create(t).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to create task")
		return nil, err
	}
	return t, nil
}

// CompleteTask marks a task as completed in the database
func (s *Service) CompleteTask(id uint, pointsVal int) (bool, error) {
	var task models.Task
	if err := s.db.Preload("SubTasks").First(&task, id).Error; err != nil {
		return false, errors.New("task not found")
	}

	task.Status = "completed"
	if err := s.db.Save(&task).Error; err != nil {
		return false, err
	}

	for _, st := range task.SubTasks {
		if st.Status != "completed" {
			if _, err := s.CompleteSubTask(st.ID, s.uc.PointsConfig.SubTask); err != nil {
				logger.Error().Err(err).Uint("subtask_id", st.ID).Msg("Failed to complete subtask while completing task")
			}
		}
	}

	milestoneReached, err := points.AwardPoints(s.db, models.CategoryTask, task.ID, pointsVal)
	if err != nil {
		logger.Error().Err(err).Msg("Error awarding points for task")
	}

	return milestoneReached, nil
}

// ListTasks retrieves all tasks from the database
func (s *Service) ListTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := s.db.Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// AddSubTask persists a new sub-task to the database
func (s *Service) AddSubTask(taskID uint, title string) (*models.SubTask, error) {
	st := &models.SubTask{
		TaskID: taskID,
		Title:  title,
	}
	if err := s.db.Create(st).Error; err != nil {
		return nil, err
	}
	return st, nil
}

// CompleteSubTask marks a sub-task as completed in the database
func (s *Service) CompleteSubTask(id uint, pointsVal int) (bool, error) {
	var st models.SubTask
	if err := s.db.First(&st, id).Error; err != nil {
		return false, errors.New("subtask not found")
	}

	st.Status = "completed"
	if err := s.db.Save(&st).Error; err != nil {
		return false, err
	}

	milestoneReached, err := points.AwardPoints(s.db, models.CategorySubTask, st.ID, pointsVal)
	if err != nil {
		logger.Error().Err(err).Msg("Error awarding points for subtask")
	}

	return milestoneReached, nil
}

// GetTasksByIntent retrieves all tasks linked to a specific intent
func (s *Service) GetTasksByIntent(intentID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := s.db.Where("IntentID = ?", intentID).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTasksByFocusSession retrieves all tasks linked to a specific focus session
func (s *Service) GetTasksByFocusSession(focusID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := s.db.Where("FocusSessionID = ?", focusID).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// DeleteTask removes a task from the database
func (s *Service) DeleteTask(id uint) error {
	if err := s.db.Where("TaskID = ?", id).Delete(&models.SubTask{}).Error; err != nil {
		return err
	}
	return s.db.Delete(&models.Task{}, id).Error
}

// DeleteSubTask removes a subtask from the database
func (s *Service) DeleteSubTask(id uint) error {
	return s.db.Delete(&models.SubTask{}, id).Error
}

// ReorderTasks updates the position of a list of tasks
func (s *Service) ReorderTasks(ids []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.Task{}).Where("id = ?", id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ReorderSubTasks updates the position of a list of subtasks
func (s *Service) ReorderSubTasks(ids []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.SubTask{}).Where("id = ?", id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetTask retrieves a single task by ID
func (s *Service) GetTask(id uint) (*models.Task, error) {
	var t models.Task
	if err := s.db.Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// RecalibrateTasks clears due dates for all pending tasks that were due in the past
func (s *Service) RecalibrateTasks() error {
	today := time.Now().Truncate(24 * time.Hour)
	return s.db.Model(&models.Task{}).Where("status = ? AND due_date < ?", "pending", today).Update("due_date", nil).Error
}
