package intent

import (
	"errors"
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

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
	result := s.DB.Find(&intents)
	return intents, result.Error
}

func (s *Service) ListActiveIntents() ([]models.Intent, error) {
	var intents []models.Intent
	result := s.DB.Where("status = ?", "active").Find(&intents)
	return intents, result.Error
}

func (s *Service) EndIntent(idStr string) (*models.Intent, error) {
	var intent models.Intent
	if err := s.DB.Where("id = ?", idStr).First(&intent).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	intent.Status = "done"
	intent.EndedAt = &now

	if err := s.DB.Save(&intent).Error; err != nil {
		return nil, err
	}

	return &intent, nil
}

func (s *Service) DeleteAll() error {
	return s.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Intent{}).Error
}
