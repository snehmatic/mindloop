package routine

import (
	"errors"
	"fmt"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based routine repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// FindRoutines retrieves all routines from the database
func (r *sqlRepository) FindRoutines() ([]models.Routine, error) {
	var routines []models.Routine
	result := r.DB.Preload("Habits").Find(&routines)
	return routines, result.Error
}

// FindRoutineByID retrieves a single routine by its ID
func (r *sqlRepository) FindRoutineByID(id uint) (*models.Routine, error) {
	var routine models.Routine
	result := r.DB.Preload("Habits").First(&routine, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("routine not found")
		}
		return nil, result.Error
	}
	return &routine, nil
}

// CreateRoutine persists a new routine to the database
func (r *sqlRepository) CreateRoutine(routine *models.Routine) error {
	return r.DB.Create(routine).Error
}

// CreateRoutines creates multiple routines in the database
func (r *sqlRepository) CreateRoutines(routines []models.Routine) error {
	return r.DB.Create(&routines).Error
}

// UpdateRoutine updates an existing routine in the database
func (r *sqlRepository) UpdateRoutine(routine *models.Routine) error {
	return r.DB.Save(routine).Error
}

// DeleteRoutine removes a routine by its ID
func (r *sqlRepository) DeleteRoutine(id uint) error {
	return r.DB.Delete(&models.Routine{}, id).Error
}

// AddHabitToRoutine assigns a habit to a routine
func (r *sqlRepository) AddHabitToRoutine(routineID, habitID uint) error {
	var habit models.Habit
	if err := r.DB.First(&habit, habitID).Error; err != nil {
		return fmt.Errorf("habit not found: %w", err)
	}
	habit.RoutineID = &routineID
	return r.DB.Save(&habit).Error
}

// RemoveHabitFromRoutine unassigns a habit from its routine
func (r *sqlRepository) RemoveHabitFromRoutine(habitID uint) error {
	var habit models.Habit
	if err := r.DB.First(&habit, habitID).Error; err != nil {
		return errors.New("habit not found")
	}
	habit.RoutineID = nil
	return r.DB.Save(&habit).Error
}
