package quest

import (
	"errors"
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) StartQuest(title string) (*models.SideQuest, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	// Check if there is already an active quest
	var quests []models.SideQuest
	if err := s.DB.Where("status = ?", "active").Limit(1).Find(&quests).Error; err != nil {
		return nil, err
	}
	if len(quests) > 0 {
		return nil, errors.New("a side quest is already active")
	}

	quest := &models.SideQuest{
		Title:  title,
		Status: "active",
	}

	if err := s.DB.Create(quest).Error; err != nil {
		return nil, err
	}
	return quest, nil
}

func (s *Service) StopQuest(id uint, note string) (*models.SideQuest, error) {
	var quest models.SideQuest
	if err := s.DB.First(&quest, id).Error; err != nil {
		return nil, err
	}

	if quest.Status != "active" {
		return nil, errors.New("side quest is not active")
	}

	quest.Status = "done"
	quest.Note = note
	now := time.Now()
	quest.EndedAt = &now

	if err := s.DB.Save(&quest).Error; err != nil {
		return nil, err
	}

	_ = points.AwardPoints(s.DB, models.CategoryQuest, quest.ID, points.PointsQuest)

	return &quest, nil
}

func (s *Service) ListQuests() ([]models.SideQuest, error) {
	var quests []models.SideQuest
	result := s.DB.Order("CreatedAt DESC").Find(&quests)
	return quests, result.Error
}

func (s *Service) GetActiveQuest() (*models.SideQuest, error) {
	var quests []models.SideQuest
	err := s.DB.Where("status = ?", "active").Limit(1).Find(&quests).Error
	if err != nil {
		return nil, err
	}
	if len(quests) == 0 {
		return nil, nil
	}
	return &quests[0], nil
}

func (s *Service) DeleteQuest(id uint) error {
	return s.DB.Delete(&models.SideQuest{}, id).Error
}
