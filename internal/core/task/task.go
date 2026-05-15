package task

import (
	"github.com/rs/zerolog"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/repository/task"
	"github.com/snehmatic/mindloop/models"
)

// Service handles business logic for tasks and sub-tasks
type Service struct {
	taskRepository task.TaskRepository
	uc             *config.UserConfig
	logger         zerolog.Logger
}

// NewService creates a new task Service instance
func NewService(repo task.TaskRepository, uc *config.UserConfig, logger zerolog.Logger) *Service {
	return &Service{
		taskRepository: repo,
		uc:             uc,
		logger:         logger,
	}
}

// CreateTask persists a new task to the database
func (s *Service) CreateTask(title string, intentID, focusID *uint) (*models.Task, error) {
	t, err := s.taskRepository.CreateTask(title, intentID, focusID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create task")
		return nil, err
	}
	return t, nil
}

// CompleteTask marks a task as completed in the database
func (s *Service) CompleteTask(id uint, pointsVal int) (bool, error) {
	completedTask, err := s.taskRepository.CompleteTask(id)
	if err != nil {
		return false, ErrorTaskNotFound
	}

	for _, st := range completedTask.SubTasks {
		if st.Status != models.TaskStatusCompleted {
			if _, err := s.CompleteSubTask(st.ID, s.uc.PointsConfig.SubTask); err != nil {
				s.logger.Error().Err(err).Uint("subtask_id", st.ID).Msg("Failed to complete subtask while completing task")
			}
		}
	}

	milestoneReached, err := points.AwardPoints(s.taskRepository.GetDB(), models.CategoryTask, completedTask.ID, pointsVal)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error awarding points for task")
	}

	return milestoneReached, nil
}

// ListTasks retrieves all tasks from the database
func (s *Service) ListTasks() ([]models.Task, error) {
	return s.taskRepository.ListTasks()
}

// AddSubTask persists a new sub-task to the database
func (s *Service) AddSubTask(taskID uint, title string) (*models.SubTask, error) {
	st, err := s.taskRepository.AddSubTask(taskID, title)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// CompleteSubTask marks a sub-task as completed in the database
func (s *Service) CompleteSubTask(id uint, pointsVal int) (bool, error) {
	st, err := s.taskRepository.CompleteSubTask(id)
	if err != nil {
		return false, ErrorSubTaskNotFound
	}

	milestoneReached, err := points.AwardPoints(s.taskRepository.GetDB(), models.CategorySubTask, st.ID, pointsVal)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error awarding points for subtask")
	}

	return milestoneReached, nil
}

// GetTasksByIntent retrieves all tasks linked to a specific intent
func (s *Service) GetTasksByIntent(intentID uint) ([]models.Task, error) {
	return s.taskRepository.GetTasksByIntent(intentID)
}

// GetTasksByFocusSession retrieves all tasks linked to a specific focus session
func (s *Service) GetTasksByFocusSession(focusID uint) ([]models.Task, error) {
	return s.taskRepository.GetTasksByFocusSession(focusID)
}

// DeleteTask removes a task from the database
func (s *Service) DeleteTask(id uint) error {
	return s.taskRepository.DeleteTask(id)
}

// DeleteSubTask removes a subtask from the database
func (s *Service) DeleteSubTask(id uint) error {
	return s.taskRepository.DeleteSubTask(id)
}

// ReorderTasks updates the position of a list of tasks
func (s *Service) ReorderTasks(ids []uint) error {
	return s.taskRepository.ReorderTasks(ids)
}

// ReorderSubTasks updates the position of a list of subtasks
func (s *Service) ReorderSubTasks(ids []uint) error {
	return s.taskRepository.ReorderSubTasks(ids)
}
