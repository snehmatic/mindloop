package focus

import (
	"errors"
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/core/hooks"
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Service handles the logic for managing focus sessions
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

	var activeSessions []models.FocusSession
	if err := s.DB.Where("status = ?", "active").Limit(1).Find(&activeSessions).Error; err != nil {
		return nil, err
	}
	if len(activeSessions) > 0 {
		return nil, errors.New("a focus session is already active")
	}

	session := &models.FocusSession{
		Title:  title,
		Status: "active",
	}

	if err := s.DB.Create(session).Error; err != nil {
		return nil, err
	}
	hooks.ExecuteHook("focus_start", map[string]string{
		"MINDLOOP_FOCUS_TITLE": session.Title,
	})
	return session, nil
}

func (s *Service) ListSessions() ([]models.FocusSession, error) {
	var sessions []models.FocusSession
	result := s.DB.Order("CreatedAt DESC").Find(&sessions)
	return sessions, result.Error
}

func (s *Service) GetSession(id int) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := s.DB.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Service) UpdateSession(session *models.FocusSession) error {
	return s.DB.Save(session).Error
}

func (s *Service) EndSession(id int, pointsToAward int) (*models.FocusSession, bool, error) {
	var session models.FocusSession
	if err := s.DB.First(&session, id).Error; err != nil {
		return nil, false, err
	}

	if session.Status != "active" {
		return nil, false, errors.New("focus session is not active")
	}

	session.Status = "ended"
	session.EndTime = time.Now()
	session.Duration = session.EndTime.Sub(session.CreatedAt).Minutes()

	if err := s.DB.Save(&session).Error; err != nil {
		return nil, false, err
	}

	milestoneReached, _ := points.AwardPoints(s.DB, models.CategoryFocus, session.ID, pointsToAward)

	hooks.ExecuteHook("focus_stop", map[string]string{
		"MINDLOOP_FOCUS_TITLE":    session.Title,
		"MINDLOOP_FOCUS_DURATION": fmt.Sprintf("%f", session.Duration),
	})
	return &session, milestoneReached, nil
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

func (s *Service) DeleteSession(id int) error {
	s.DB.Model(&models.Task{}).Where("FocusSessionID = ?", id).Update("FocusSessionID", nil)
	return s.DB.Delete(&models.FocusSession{}, id).Error
}

func (s *Service) DeleteAll() error {
	return s.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.FocusSession{}).Error
}

func (s *Service) PauseSession(id uint) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := s.DB.First(&session, id).Error; err != nil {
		return nil, err
	}

	if session.Status != "active" {
		return nil, errors.New("focus session is not active")
	}

	session.Status = "paused"
	if err := s.DB.Save(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Service) ResumeSession(id uint) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := s.DB.First(&session, id).Error; err != nil {
		return nil, err
	}

	if session.Status != "paused" {
		return nil, errors.New("focus session is not paused")
	}

	session.Status = "active"
	if err := s.DB.Save(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Service) GetActiveSession() (*models.FocusSession, error) {
	var sessions []models.FocusSession
	err := s.DB.Where("status = ?", "active").Limit(1).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}
	return &sessions[0], nil
}
