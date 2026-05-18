// Package routine manages groupings of habits for specific times of day
package routine

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/repository/routine"
	"github.com/snehmatic/mindloop/models"
)

var logger = log.Get()

type Service struct {
	repo   routine.Repository
	logger log.Logger
}

func NewService(repo routine.Repository, logger log.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateRoutine persists a new routine to the database
func (s *Service) CreateRoutine(title, timeOfDay string) (*models.Routine, error) {
	r := &models.Routine{
		Title:     title,
		TimeOfDay: timeOfDay,
	}
	if err := s.repo.CreateRoutine(r); err != nil {
		logger.Error("Failed to create routine", err)
		return nil, fmt.Errorf("failed to create routine: %w", err)
	}
	return r, nil
}

// ListRoutines retrieves all routines from the database
func (s *Service) ListRoutines() ([]models.Routine, error) {
	routines, err := s.repo.FindRoutines()
	return routines, err
}

func (s *Service) GetRoutine(id uint) (*models.Routine, error) {
	routineEntry, err := s.repo.FindRoutineByID(id)
	if err != nil {
		return nil, err
	}
	return routineEntry, nil
}

func (s *Service) UpdateRoutine(r *models.Routine) error {
	return s.repo.UpdateRoutine(r)
}

func (s *Service) DeleteRoutine(id uint) error {
	return s.repo.DeleteRoutine(id)
}

func (s *Service) AddHabitToRoutine(routineID, habitID uint) error {
	return s.repo.AddHabitToRoutine(routineID, habitID)
}

func (s *Service) RemoveHabitFromRoutine(habitID uint) error {
	return s.repo.RemoveHabitFromRoutine(habitID)
}
