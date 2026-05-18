package habit

import (
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB     *gorm.DB
	logger log.Logger
}

// NewSQLRepository creates a new SQL-based habit repository
func NewSQLRepository(db *gorm.DB, logger log.Logger) Repository {
	return &sqlRepository{DB: db, logger: logger}
}

// CreateHabit creates a new habit in the database
func (r *sqlRepository) CreateHabit(habit *models.Habit) error {
	if habit == nil {
		return ErrHabitCannotBeNil
	}
	if err := habit.ValidateHabit(); err != nil {
		return err
	}
	return r.DB.Create(habit).Error
}

// DeleteHabit removes a habit from the database by its ID
func (r *sqlRepository) DeleteHabit(id string) error {
	var habits []models.Habit
	if err := r.DB.Where("id = ?", id).Limit(1).Find(&habits).Error; err != nil {
		return err
	}
	if len(habits) == 0 {
		return ErrHabitNotFound
	}
	return r.DB.Delete(&habits[0]).Error
}

// GetHabit retrieves a single habit by its ID
func (r *sqlRepository) GetHabit(id string) (*models.Habit, error) {
	var habits []models.Habit
	if err := r.DB.Where("id = ?", id).Limit(1).Find(&habits).Error; err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return nil, ErrHabitNotFound
	}
	return &habits[0], nil
}

// UpdateHabit modifies an existing habit in the database
func (r *sqlRepository) UpdateHabit(habit *models.Habit) error {
	if habit == nil {
		return ErrHabitCannotBeNil
	}
	if err := habit.ValidateHabit(); err != nil {
		return err
	}
	return r.DB.Save(habit).Error
}

// ListHabits retrieves habits based on interval type
func (r *sqlRepository) ListHabits(interval models.IntervalType) ([]models.Habit, error) {
	var habits []models.Habit
	query := r.DB
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}
	// Only show habits that haven't ended yet
	now := time.Now()
	query = query.Where("EndDate IS NULL OR EndDate > ?", now)

	result := query.Order("CreatedAt DESC").Find(&habits)
	return habits, result.Error
}

// ListEndedHabits retrieves habits that have ended
func (r *sqlRepository) ListEndedHabits() ([]models.Habit, error) {
	var habits []models.Habit
	now := time.Now()
	result := r.DB.Where("EndDate IS NOT NULL AND EndDate <= ?", now).Order("CreatedAt DESC").Find(&habits)
	return habits, result.Error
}

// LogHabit records an instance of a habit being performed and awards points if completed
func (r *sqlRepository) LogHabit(habitID string, pointsToAward int) (*models.Habit, *models.HabitLog, bool, error) {
	habit, err := r.GetHabit(habitID)
	if err != nil {
		return nil, nil, false, err
	}

	var existingLogs []models.HabitLog
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	var startRange, endRange time.Time

	switch habit.Interval {
	case models.Daily:
		startRange = today
		endRange = today.AddDate(0, 0, 1)
	case models.Weekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday
		}
		startRange = today.AddDate(0, 0, -(weekday - 1))
		endRange = startRange.AddDate(0, 0, 7)
	}

	res := r.DB.Where("HabitID = ? AND CreatedAt >= ? AND CreatedAt < ?", habit.ID, startRange, endRange).Limit(1).Find(&existingLogs)
	if res.Error != nil {
		return nil, nil, false, res.Error
	}

	milestoneReached := false

	if len(existingLogs) > 0 {
		existingLog := existingLogs[0]
		if existingLog.ActualCount >= habit.TargetCount {
			return habit, &existingLog, false, ErrHabitAlreadyCompleted
		}

		existingLog.ActualCount++
		if err := r.DB.Save(&existingLog).Error; err != nil {
			return nil, nil, false, err
		}

		if existingLog.ActualCount == habit.TargetCount {
			var totalPoints int
			if err := r.DB.Model(&models.PointTransaction{}).Select("COALESCE(SUM(Points), 0)").Scan(&totalPoints).Error; err == nil {
				tx := models.PointTransaction{
					ActivityType: models.CategoryHabit,
					ActivityID:   habit.ID,
					Points:       pointsToAward,
				}
				if r.DB.Create(&tx).Error == nil {
					newTotal := totalPoints + pointsToAward
					if newTotal/points.MilestoneInterval > totalPoints/points.MilestoneInterval {
						milestoneReached = true
					}
				}
			}
		}

		return habit, &existingLog, milestoneReached, nil
	}

	// Create new log
	habitLog := &models.HabitLog{
		HabitID:     habit.ID,
		Title:       habit.Title,
		Interval:    habit.Interval,
		TargetCount: habit.TargetCount,
		ActualCount: 1,
		EndedAt:     endRange.AddDate(0, 0, -1), // Represents the last day of the interval
	}
	if err := r.DB.Create(habitLog).Error; err != nil {
		return nil, nil, false, err
	}

	if habitLog.ActualCount == habit.TargetCount {
		var totalPoints int
		if err := r.DB.Model(&models.PointTransaction{}).Select("COALESCE(SUM(Points), 0)").Scan(&totalPoints).Error; err == nil {
			tx := models.PointTransaction{
				ActivityType: models.CategoryHabit,
				ActivityID:   habit.ID,
				Points:       pointsToAward,
			}
			if r.DB.Create(&tx).Error == nil {
				newTotal := totalPoints + pointsToAward
				if newTotal/points.MilestoneInterval > totalPoints/points.MilestoneInterval {
					milestoneReached = true
				}
			}
		}
	}

	return habit, habitLog, milestoneReached, nil
}

// UnlogHabit removes an instance of a habit being performed
func (r *sqlRepository) UnlogHabit(habitID string) (*models.Habit, error) {
	habit, err := r.GetHabit(habitID)
	if err != nil {
		return nil, err
	}

	var existingLogs []models.HabitLog
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	var startRange, endRange time.Time

	switch habit.Interval {
	case models.Daily:
		startRange = today
		endRange = today.AddDate(0, 0, 1)
	case models.Weekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startRange = today.AddDate(0, 0, -(weekday - 1))
		endRange = startRange.AddDate(0, 0, 7)
	}

	res := r.DB.Where("HabitID = ? AND CreatedAt >= ? AND CreatedAt < ?", habit.ID, startRange, endRange).Limit(1).Find(&existingLogs)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(existingLogs) == 0 {
		return nil, ErrNoExistingLog
	}

	existingLog := existingLogs[0]
	if existingLog.ActualCount <= 0 {
		return nil, ErrHabitAlreadyUndone
	}

	existingLog.ActualCount--
	if err := r.DB.Save(&existingLog).Error; err != nil {
		return nil, err
	}

	return habit, nil
}

// ListHabitLogs retrieves habit logs based on interval type
func (r *sqlRepository) ListHabitLogs(interval models.IntervalType) ([]models.HabitLog, error) {
	var habitLogs []models.HabitLog
	query := r.DB
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}
	result := query.Order("CreatedAt DESC").Find(&habitLogs)
	return habitLogs, result.Error
}

// ListLogsForHabit retrieves all logs for a specific habit
func (r *sqlRepository) ListLogsForHabit(habitID uint) ([]models.HabitLog, error) {
	var habitLogs []models.HabitLog
	result := r.DB.Where("HabitID = ?", habitID).Order("CreatedAt ASC").Find(&habitLogs)
	return habitLogs, result.Error
}

// DeleteAll removes all habits and habit logs from the database
func (r *sqlRepository) DeleteAll() error {
	// Transaction to delete both logs and habits
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.HabitLog{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Habit{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// CalculateStreak calculates the streak for a habit
func (r *sqlRepository) CalculateStreak(habitID uint, interval models.IntervalType) (int, error) {
	if interval != models.Daily {
		return 0, nil // fast track for non-daily for now
	}

	var logs []models.HabitLog
	// Fetch all logs for this habit ordered by date descending
	if err := r.DB.Where("HabitID = ?", habitID).Order("CreatedAt desc").Find(&logs).Error; err != nil {
		return 0, err
	}

	if len(logs) == 0 {
		return 0, nil
	}

	streak := 0
	today := time.Now().Truncate(24 * time.Hour)

	lastLogDate := logs[0].CreatedAt.Truncate(24 * time.Hour)
	daysDiff := today.Sub(lastLogDate).Hours() / 24

	if daysDiff > 1 {
		return 0, nil
	}

	expectedDate := lastLogDate
loop:
	for _, log := range logs {
		logDate := log.CreatedAt.Truncate(24 * time.Hour)

		switch {
		case logDate.Equal(expectedDate):
			if log.ActualCount >= log.TargetCount {
				streak++
				expectedDate = expectedDate.AddDate(0, 0, -1)
			}
		case logDate.After(expectedDate):
			continue
		default:
			break loop
		}
	}

	return streak, nil
}

// GetHabitStats returns habit statistics for the given time range
func (r *sqlRepository) GetHabitStats(start, end time.Time) ([]models.HabitStats, error) {
	var habits []models.Habit
	if err := r.DB.Order("CreatedAt DESC").Find(&habits).Error; err != nil {
		return nil, err
	}
	if len(habits) == 0 {
		return []models.HabitStats{}, nil
	}

	var habitLogs []models.HabitLog
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"
	if err := r.DB.Where(rangeQuery, start, end).Order("CreatedAt DESC").Find(&habitLogs).Error; err != nil {
		return nil, err
	}

	var stats []models.HabitStats
	for _, habit := range habits {
		totalCompletedLogsForHabit := 0
		totalLogsForHabit := 0
		for _, log := range habitLogs {
			if log.HabitID == habit.ID {
				totalLogsForHabit++
				if log.ActualCount >= log.TargetCount {
					totalCompletedLogsForHabit++
				}
			}
		}
		if totalLogsForHabit > 0 {
			stats = append(stats, models.HabitStats{
				HabitName:      habit.Title,
				CompletionRate: float64(totalCompletedLogsForHabit) * 100 / float64(totalLogsForHabit),
				LogsTracked:    totalLogsForHabit,
				LogsCompleted:  totalCompletedLogsForHabit,
			})
		}
	}
	return stats, nil
}

// GetDB returns the GORM DB instance
func (r *sqlRepository) GetDB() *gorm.DB {
	return r.DB
}

// GetAllHabits retrieves all habits ordered by creation date descending
func (r *sqlRepository) GetAllHabits() ([]models.Habit, error) {
	var habits []models.Habit
	if err := r.DB.Order("CreatedAt DESC").Find(&habits).Error; err != nil {
		return nil, err
	}
	return habits, nil
}

// GetHabitLogsInRange retrieves all habit logs within the given time range
func (r *sqlRepository) GetHabitLogsInRange(start, end time.Time) ([]models.HabitLog, error) {
	var habitLogs []models.HabitLog
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"
	if err := r.DB.Where(rangeQuery, start, end).Order("CreatedAt DESC").Find(&habitLogs).Error; err != nil {
		return nil, err
	}
	return habitLogs, nil
}
