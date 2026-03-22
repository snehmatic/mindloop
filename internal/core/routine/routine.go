// Package routine manages groupings of habits for specific times of day
package routine

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

// CreateRoutine persists a new routine to the database
func (s *Service) CreateRoutine(title, timeOfDay string) (*models.Routine, error) {
	r := &models.Routine{
		Title:     title,
		TimeOfDay: timeOfDay,
	}
	result := s.db.Create(r)
	if result.Error != nil {
		logger.Error().Err(result.Error).Msg("Failed to create routine")
		return nil, result.Error
	}
	return r, nil
}

// ListRoutines retrieves all routines from the database
func (s *Service) ListRoutines() ([]models.Routine, error) {
	var routines []models.Routine
	result := s.db.Preload("Habits").Find(&routines)
	if result.Error != nil {
		return nil, result.Error
	}
	return routines, nil
}

func (s *Service) GetRoutine(id uint) (*models.Routine, error) {
	var routine models.Routine
	result := s.db.Preload("Habits").First(&routine, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("routine not found")
		}
		return nil, result.Error
	}
	return &routine, nil
}

func (s *Service) UpdateRoutine(r *models.Routine) error {
	result := s.db.Save(r)
	return result.Error
}

func (s *Service) DeleteRoutine(id uint) error {
	result := s.db.Delete(&models.Routine{}, id)
	return result.Error
}

func (s *Service) AddHabitToRoutine(routineID, habitID uint) error {
	var habit models.Habit
	if err := s.db.First(&habit, habitID).Error; err != nil {
		return errors.New("habit not found")
	}

	habit.RoutineID = &routineID
	if err := s.db.Save(&habit).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveHabitFromRoutine(habitID uint) error {
	var habit models.Habit
	if err := s.db.First(&habit, habitID).Error; err != nil {
		return errors.New("habit not found")
	}

	habit.RoutineID = nil
	if err := s.db.Save(&habit).Error; err != nil {
		return err
	}
	return nil
}
