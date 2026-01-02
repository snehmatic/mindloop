package focus

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

func (s *Service) StartSession(title string) (*models.FocusSession, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	session := &models.FocusSession{
		Title:  title,
		Status: "active",
	}

	if err := s.DB.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) ListSessions() ([]models.FocusSession, error) {
	var sessions []models.FocusSession
	result := s.DB.Find(&sessions)
	return sessions, result.Error
}

func (s *Service) EndSession(id int) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := s.DB.First(&session, id).Error; err != nil {
		return nil, err
	}

	if session.Status != "active" {
		return nil, errors.New("focus session is not active")
	}

	session.Status = "ended"
	session.EndTime = time.Now()
	session.Duration = session.EndTime.Sub(session.CreatedAt).Minutes()

	if err := s.DB.Save(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Service) RateSession(id int, rating int) (*models.FocusSession, error) {
	if rating < 0 || rating > 10 {
		return nil, errors.New("rating must be between 0 and 10")
	}

	var session models.FocusSession
	if err := s.DB.First(&session, id).Error; err != nil {
		return nil, err
	}

	if session.Status != "ended" {
		return nil, errors.New("focus session is not ended")
	}

	session.Rating = rating
	if err := s.DB.Save(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}
