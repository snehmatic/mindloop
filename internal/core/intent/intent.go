package intent

import (
	"errors"
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Service handles the logic for managing user intents
type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) StartIntent(name string) (*models.Intent, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	intent := &models.Intent{
		Name:   name,
		Status: "active",
	}

	if err := s.DB.Create(intent).Error; err != nil {
		return nil, err
	}
	return intent, nil
}

func (s *Service) ListIntents() ([]models.Intent, error) {
	var intents []models.Intent
	result := s.DB.Order("CreatedAt DESC").Find(&intents)
	return intents, result.Error
}

func (s *Service) ListActiveIntents() ([]models.Intent, error) {
	var intents []models.Intent
	result := s.DB.Where("status = ?", "active").Order("CreatedAt DESC").Find(&intents)
	return intents, result.Error
}

func (s *Service) GetOngoingIntent() (*models.Intent, error) {
	var intents []models.Intent
	result := s.DB.Where("status IN ?", []string{"active", "paused"}).Limit(1).Find(&intents)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(intents) == 0 {
		return nil, nil
	}
	return &intents[0], nil
}

func (s *Service) GetIntent(id string) (*models.Intent, error) {
	var intent models.Intent
	if err := s.DB.Where("id = ?", id).First(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

func (s *Service) UpdateIntent(intent *models.Intent) error {
	return s.DB.Save(intent).Error
}

func (s *Service) EndIntent(idStr string, pointsToAward int) (*models.Intent, bool, error) {
	var intent models.Intent
	if err := s.DB.Where("id = ?", idStr).First(&intent).Error; err != nil {
		return nil, false, err
	}

	now := time.Now()
	intent.Status = "done"
	intent.EndedAt = &now

	if err := s.DB.Save(&intent).Error; err != nil {
		return nil, false, err
	}

	milestoneReached, _ := points.AwardPoints(s.DB, models.CategoryIntent, intent.ID, pointsToAward)

	return &intent, milestoneReached, nil
}

func (s *Service) DeleteIntent(id string) error {
	s.DB.Model(&models.Task{}).Where("IntentID = ?", id).Update("IntentID", nil)
	return s.DB.Delete(&models.Intent{}, "id = ?", id).Error
}

func (s *Service) DeleteAll() error {
	return s.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Intent{}).Error
}

func (s *Service) PauseIntent(id uint) (*models.Intent, error) {
	var intent models.Intent
	if err := s.DB.First(&intent, id).Error; err != nil {
		return nil, err
	}

	if intent.Status != "active" {
		return nil, errors.New("intent is not active")
	}

	intent.Status = "paused"
	if err := s.DB.Save(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

func (s *Service) ResumeIntent(id uint) (*models.Intent, error) {
	var intent models.Intent
	if err := s.DB.First(&intent, id).Error; err != nil {
		return nil, err
	}

	if intent.Status != "paused" {
		return nil, errors.New("intent is not paused")
	}

	intent.Status = "active"
	if err := s.DB.Save(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}
