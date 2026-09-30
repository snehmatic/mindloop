package focus

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/core/hooks"
	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
)

// Service handles the logic for managing focus sessions
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) StartSession(title string) (*models.FocusSession, error) {
	if title == "" {
		return nil, ErrTitleCannotBeEmpty
	}

	activeSession, err := s.repo.GetActiveSession()
	if err != nil {
		return nil, err
	}
	if activeSession != nil {
		return nil, ErrAFocusSessionIsAlreadyActive
	}

	session := &models.FocusSession{
		Title:  title,
		Status: models.StatusActive,
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}
	hooks.ExecuteHook("focus_start", map[string]string{
		"MINDLOOP_FOCUS_TITLE": session.Title,
	})
	return session, nil
}

func (s *Service) ListSessions() ([]models.FocusSession, error) {
	return s.repo.ListSessions()
}

func (s *Service) GetSession(id int) (*models.FocusSession, error) {
	return s.repo.GetSession(id)
}

func (s *Service) UpdateSession(session *models.FocusSession) error {
	return s.repo.UpdateSession(session)
}

func (s *Service) EndSession(id int, pointsToAward int) (*models.FocusSession, bool, error) {
	session, err := s.repo.GetSession(id)
	if err != nil {
		return nil, false, err
	}

	if session.Status != models.StatusActive && session.Status != models.StatusPaused {
		return nil, false, ErrFocusSessionIsNotActive
	}

	session.EndTime = time.Now()
	durationSeconds := session.EndTime.Sub(session.CreatedAt).Seconds() - float64(session.PausedDuration)
	if session.Status == models.StatusPaused && session.LastPausedAt != nil {
		durationSeconds -= session.EndTime.Sub(*session.LastPausedAt).Seconds()
	}

	session.Status = models.StatusEnded
	session.Duration = durationSeconds / 60.0

	if err := s.repo.UpdateSession(session); err != nil {
		return nil, false, err
	}

	milestoneReached, _ := points.AwardPoints(s.repo.GetDB(), models.CategoryFocus, session.ID, pointsToAward)

	hooks.ExecuteHook("focus_stop", map[string]string{
		"MINDLOOP_FOCUS_TITLE":    session.Title,
		"MINDLOOP_FOCUS_DURATION": fmt.Sprintf("%f", session.Duration),
	})
	return session, milestoneReached, nil
}

func (s *Service) RateSession(id int, rating int) (*models.FocusSession, error) {
	if rating < 0 || rating > 10 {
		return nil, ErrRatingMustBeBetween0And10
	}

	session, err := s.repo.GetSession(id)
	if err != nil {
		return nil, err
	}

	if session.Status != models.StatusEnded {
		return nil, ErrFocusSessionIsNotEnded
	}

	session.Rating = rating
	if err := s.repo.UpdateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) DeleteSession(id int) error {
	return s.repo.DeleteSession(id)
}

func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *Service) PauseSession(id uint) (*models.FocusSession, error) {
	session, err := s.repo.GetSessionByID(id)
	if err != nil {
		return nil, err
	}

	if session.Status != models.StatusActive {
		return nil, ErrFocusSessionIsNotActive
	}

	session.Status = models.StatusPaused
	now := time.Now()
	session.LastPausedAt = &now
	if err := s.repo.UpdateSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) ResumeSession(id uint) (*models.FocusSession, error) {
	session, err := s.repo.GetSessionByID(id)
	if err != nil {
		return nil, err
	}

	if session.Status != models.StatusPaused {
		return nil, ErrFocusSessionIsNotPaused
	}

	session.Status = models.StatusActive
	if session.LastPausedAt != nil {
		session.PausedDuration += int(time.Since(*session.LastPausedAt).Seconds())
		session.LastPausedAt = nil
	}
	if err := s.repo.UpdateSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) GetActiveSession() (*models.FocusSession, error) {
	return s.repo.GetActiveSession()
}
