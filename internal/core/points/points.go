package points

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Point constants
const (
	PointsFocus   = 10
	PointsHabit   = 5
	PointsIntent  = 10
	PointsJournal = 5
	PointsQuest   = 5
)

var MilestoneInterval = 100

// AwardPoints creates a new PointTransaction and returns true if a milestone was reached
func AwardPoints(db *gorm.DB, activityType models.PointCategory, activityID uint, points int) (bool, error) {
	// Get points before transaction
	currentTotal, err := GetTotalPoints(db)
	if err != nil {
		return false, err
	}

	transaction := models.PointTransaction{
		ActivityType: activityType,
		ActivityID:   activityID,
		Points:       points,
	}
	err = db.Create(&transaction).Error
	if err != nil {
		return false, err
	}

	newTotal := currentTotal + points
	
	// Check if a milestone boundary was crossed
	currentMilestone := currentTotal / MilestoneInterval
	newMilestone := newTotal / MilestoneInterval

	return newMilestone > currentMilestone, nil
}

// GetTotalPoints calculates the lifetime total points for the user
func GetTotalPoints(db *gorm.DB) (int, error) {
	var total int
	err := db.Model(&models.PointTransaction{}).Select("COALESCE(SUM(points), 0)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetPointsInRange returns the point transactions within a date range
func GetPointsInRange(db *gorm.DB, start, end string) ([]models.PointTransaction, error) {
	var history []models.PointTransaction
	query := db.Model(&models.PointTransaction{})
	if start != "" && end != "" {
		query = query.Where("created_at >= ? AND created_at <= ?", start, end)
	}
	err := query.Order("created_at asc").Find(&history).Error
	return history, err
}
