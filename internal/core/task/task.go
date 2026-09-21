package task

import (
	"time"

	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/nlp"
	"github.com/snehmatic/mindloop/models"
)

var logger = log.Get()

// Service handles business logic for tasks and sub-tasks
type Service struct {
	repo Repository
}

// NewService creates a new task Service instance
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
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
	if err := s.repo.CreateTask(t); err != nil {
		logger.Error().Err(err).Msg("Failed to create task")
		return nil, err
	}
	return t, nil
}

// CompleteTask marks a task as completed in the database
func (s *Service) CompleteTask(id uint, pointsVal int) (bool, error) {
	task, err := s.repo.GetTask(id)
	if err != nil {
		return false, ErrTaskNotFound
	}

	task.Status = models.StatusCompleted
	if err := s.repo.UpdateTask(task); err != nil {
		return false, err
	}

	for _, st := range task.SubTasks {
		if st.Status != models.StatusCompleted {
			if _, err := s.CompleteSubTask(st.ID, config.GetUserConfig().PointsConfig.SubTask); err != nil {
				logger.Error().Err(err).Uint("subtask_id", st.ID).Msg("Failed to complete subtask while completing task")
			}
		}
	}

	milestoneReached, err := points.AwardPoints(s.repo.GetDB(), models.CategoryTask, task.ID, pointsVal)
	if err != nil {
		logger.Error().Err(err).Msg("Error awarding points for task")
	}

	return milestoneReached, nil
}

// ListTasks retrieves all tasks from the database
func (s *Service) ListTasks() ([]models.Task, error) {
	return s.repo.ListTasks()
}

// AddSubTask persists a new sub-task to the database
func (s *Service) AddSubTask(taskID uint, title string) (*models.SubTask, error) {
	st := &models.SubTask{
		TaskID: taskID,
		Title:  title,
	}
	if err := s.repo.CreateSubTask(st); err != nil {
		return nil, err
	}
	return st, nil
}

// CompleteSubTask marks a sub-task as completed in the database
func (s *Service) CompleteSubTask(id uint, pointsVal int) (bool, error) {
	st, err := s.repo.GetSubTask(id)
	if err != nil {
		return false, ErrSubtaskNotFound
	}

	st.Status = models.StatusCompleted
	if err := s.repo.UpdateSubTask(st); err != nil {
		return false, err
	}

	milestoneReached, err := points.AwardPoints(s.repo.GetDB(), models.CategorySubTask, st.ID, pointsVal)
	if err != nil {
		logger.Error().Err(err).Msg("Error awarding points for subtask")
	}

	return milestoneReached, nil
}

// GetTasksByIntent retrieves all tasks linked to a specific intent
func (s *Service) GetTasksByIntent(intentID uint) ([]models.Task, error) {
	return s.repo.GetTasksByIntent(intentID)
}

// GetTasksByFocusSession retrieves all tasks linked to a specific focus session
func (s *Service) GetTasksByFocusSession(focusID uint) ([]models.Task, error) {
	return s.repo.GetTasksByFocusSession(focusID)
}

// DeleteTask removes a task from the database
func (s *Service) DeleteTask(id uint) error {
	if err := s.repo.DeleteSubTasksByTaskID(id); err != nil {
		return err
	}
	return s.repo.DeleteTask(id)
}

// DeleteSubTask removes a subtask from the database
func (s *Service) DeleteSubTask(id uint) error {
	return s.repo.DeleteSubTask(id)
}

// ReorderTasks updates the position of a list of tasks
func (s *Service) ReorderTasks(ids []uint) error {
	return s.repo.ReorderTasks(ids)
}

// ReorderSubTasks updates the position of a list of subtasks
func (s *Service) ReorderSubTasks(ids []uint) error {
	return s.repo.ReorderSubTasks(ids)
}

// GetTask retrieves a single task by ID
func (s *Service) GetTask(id uint) (*models.Task, error) {
	return s.repo.GetTask(id)
}

// RecalibrateTasks clears due dates for all pending tasks that were due in the past
func (s *Service) RecalibrateTasks() error {
	today := time.Now().Truncate(24 * time.Hour)
	return s.repo.RecalibrateTasks(models.StatusPending, today)
}
