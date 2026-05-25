package focus

import (
	"github.com/snehmatic/mindloop/internal/repository/focus"
	"github.com/snehmatic/mindloop/models"
)

// Service handles the logic for managing focus sessions
type Service struct {
	repository focus.Repository
}

func NewService(repo focus.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) StartSession(title string) (*models.FocusSession, error) {
	return s.repository.StartSession(title)
}

func (s *Service) ListSessions() ([]models.FocusSession, error) {
	return s.repository.ListSessions()
}

func (s *Service) GetSession(id int) (*models.FocusSession, error) {
	return s.repository.GetSession(id)
}

func (s *Service) UpdateSession(session *models.FocusSession) error {
	return s.repository.UpdateSession(session)
}

func (s *Service) EndSession(id int, pointsToAward int) (*models.FocusSession, bool, error) {
	return s.repository.EndSession(id, pointsToAward)
}

func (s *Service) RateSession(id int, rating int) (*models.FocusSession, error) {
	return s.repository.RateSession(id, rating)
}

func (s *Service) DeleteSession(id int) error {
	return s.repository.DeleteSession(id)
}

func (s *Service) DeleteAll() error {
	return s.repository.DeleteAll()
}

func (s *Service) PauseSession(id uint) (*models.FocusSession, error) {
	return s.repository.PauseSession(id)
}

func (s *Service) ResumeSession(id uint) (*models.FocusSession, error) {
	return s.repository.ResumeSession(id)
}

func (s *Service) GetActiveSession() (*models.FocusSession, error) {
	return s.repository.GetActiveSession()
}
