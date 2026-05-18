package quest

import (
	"errors"
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

// NewSQLRepository creates a new SQL-based side quest repository
func NewSQLRepository(db *gorm.DB, logger log.Logger) Repository {
	return &sqlRepository{DB: db, logger: logger}
}

// StartQuest creates a new side quest
func (r *sqlRepository) StartQuest(title string) (*models.SideQuest, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	// Check if there is already an active quest
	var quests []models.SideQuest
	if err := r.DB.Where("status = ?", "active").Limit(1).Find(&quests).Error; err != nil {
		return nil, err
	}
	if len(quests) > 0 {
		return nil, errors.New("a side quest is already active")
	}

	quest := &models.SideQuest{
		Title:  title,
		Status: "active",
	}

	if err := r.DB.Create(quest).Error; err != nil {
		return nil, err
	}
	return quest, nil
}

// StopQuest ends a side quest and awards points if successful
func (r *sqlRepository) StopQuest(id uint, note string, pointsToAward int) (*models.SideQuest, bool, error) {
	var quest models.SideQuest
	if err := r.DB.First(&quest, id).Error; err != nil {
		return nil, false, err
	}

	if quest.Status != "active" {
		return nil, false, errors.New("side quest is not active")
	}

	quest.Status = "done"
	quest.Note = note
	now := time.Now()
	quest.EndedAt = &now

	if err := r.DB.Save(&quest).Error; err != nil {
		return nil, false, err
	}

	var milestoneReached bool
	var totalPoints int
	if err := r.DB.Model(&models.PointTransaction{}).Select("COALESCE(SUM(Points), 0)").Scan(&totalPoints).Error; err == nil {
		tx := models.PointTransaction{
			ActivityType: models.CategoryQuest,
			ActivityID:   quest.ID,
			Points:       pointsToAward,
		}
		if r.DB.Create(&tx).Error == nil {
			newTotal := totalPoints + pointsToAward
			if newTotal/points.MilestoneInterval > totalPoints/points.MilestoneInterval {
				milestoneReached = true
			}
		}
	}

	return &quest, milestoneReached, nil
}

// ListQuests retrieves all side quests from the database
func (r *sqlRepository) ListQuests() ([]models.SideQuest, error) {
	var quests []models.SideQuest
	result := r.DB.Order("CreatedAt DESC").Find(&quests)
	return quests, result.Error
}

// GetActiveQuest retrieves the currently active side quest
func (r *sqlRepository) GetActiveQuest() (*models.SideQuest, error) {
	var quests []models.SideQuest
	err := r.DB.Where("status = ?", "active").Limit(1).Find(&quests).Error
	if err != nil {
		return nil, err
	}
	if len(quests) == 0 {
		return nil, nil
	}
	return &quests[0], nil
}

// DeleteQuest removes a side quest from the database
func (r *sqlRepository) DeleteQuest(id uint) error {
	return r.DB.Delete(&models.SideQuest{}, id).Error
}
