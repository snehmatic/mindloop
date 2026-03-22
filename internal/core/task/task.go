package task

import (
	"errors"

	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

var logger = log.Get()

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateTask(title string, intentID, focusID *uint) (*models.Task, error) {
	t := &models.Task{
		Title:          title,
		IntentID:       intentID,
		FocusSessionID: focusID,
	}
	if err := s.db.Create(t).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to create task")
		return nil, err
	}
	return t, nil
}

func (s *Service) CompleteTask(id uint) error {
	var task models.Task
	if err := s.db.First(&task, id).Error; err != nil {
		return errors.New("task not found")
	}

	task.Status = "completed"
	if err := s.db.Save(&task).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ListTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := s.db.Preload("SubTasks").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

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

func (s *Service) CompleteSubTask(id uint) error {
	var st models.SubTask
	if err := s.db.First(&st, id).Error; err != nil {
		return errors.New("subtask not found")
	}

	st.Status = "completed"
	if err := s.db.Save(&st).Error; err != nil {
		return err
	}
	return nil
}
