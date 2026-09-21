package habit

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

func (r *SQLRepository) CreateHabit(habit *models.Habit) error {
	return r.db.Create(habit).Error
}

func (r *SQLRepository) GetHabit(id string) (*models.Habit, error) {
	var habits []models.Habit
	if err := r.db.Where("id = ?", id).Limit(1).Find(&habits).Error; err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return nil, ErrHabitNotFound
	}
	return &habits[0], nil
}

func (r *SQLRepository) UpdateHabit(habit *models.Habit) error {
	return r.db.Save(habit).Error
}

func (r *SQLRepository) DeleteHabit(id string) error {
	var habits []models.Habit
	if err := r.db.Where("id = ?", id).Limit(1).Find(&habits).Error; err != nil {
		return err
	}
	if len(habits) == 0 {
		return ErrHabitNotFound
	}
	return r.db.Delete(&habits[0]).Error
}

func (r *SQLRepository) ListHabits(interval models.IntervalType, now time.Time) ([]models.Habit, error) {
	var habits []models.Habit
	query := r.db
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}
	query = query.Where("EndDate IS NULL OR EndDate > ?", now)

	result := query.Order("CreatedAt DESC").Find(&habits)
	return habits, result.Error
}

func (r *SQLRepository) ListEndedHabits(now time.Time) ([]models.Habit, error) {
	var habits []models.Habit
	result := r.db.Where("EndDate IS NOT NULL AND EndDate <= ?", now).Order("CreatedAt DESC").Find(&habits)
	return habits, result.Error
}

func (r *SQLRepository) GetHabitLogForInterval(habitID uint, startRange, endRange time.Time) (*models.HabitLog, error) {
	var existingLogs []models.HabitLog
	res := r.db.Where("HabitID = ? AND CreatedAt >= ? AND CreatedAt < ?", habitID, startRange, endRange).Limit(1).Find(&existingLogs)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(existingLogs) == 0 {
		return nil, nil // Not found
	}
	return &existingLogs[0], nil
}

func (r *SQLRepository) CreateHabitLog(log *models.HabitLog) error {
	return r.db.Create(log).Error
}

func (r *SQLRepository) UpdateHabitLog(log *models.HabitLog) error {
	return r.db.Save(log).Error
}

func (r *SQLRepository) ListHabitLogs(interval models.IntervalType) ([]models.HabitLog, error) {
	var habitLogs []models.HabitLog
	query := r.db
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}
	result := query.Order("CreatedAt DESC").Find(&habitLogs)
	return habitLogs, result.Error
}

func (r *SQLRepository) ListLogsForHabit(habitID uint) ([]models.HabitLog, error) {
	var habitLogs []models.HabitLog
	result := r.db.Where("HabitID = ?", habitID).Order("CreatedAt ASC").Find(&habitLogs)
	return habitLogs, result.Error
}

func (r *SQLRepository) ListLogsForHabits(habitIDs []uint) ([]models.HabitLog, error) {
	var allLogs []models.HabitLog
	if err := r.db.Where("HabitID IN ?", habitIDs).Order("CreatedAt asc").Find(&allLogs).Error; err != nil {
		return nil, err
	}
	return allLogs, nil
}

func (r *SQLRepository) DeleteAll() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.HabitLog{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Habit{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *SQLRepository) RecalibrateAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&models.Habit{}).Update("RecalibratedAt", time.Now()).Error
}
