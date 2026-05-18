package routine

import (
	"github.com/snehmatic/mindloop/models"
)

// Repository defines the interface for routine data access
type Repository interface {
	FindRoutines() ([]models.Routine, error)
	FindRoutineByID(id uint) (*models.Routine, error)
	CreateRoutine(r *models.Routine) error
	CreateRoutines(routines []models.Routine) error
	UpdateRoutine(r *models.Routine) error
	DeleteRoutine(id uint) error
	AddHabitToRoutine(routineID, habitID uint) error
	RemoveHabitFromRoutine(habitID uint) error
}
