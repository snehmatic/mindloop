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

func (s *Service) DeleteIntent(id string) error {
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
