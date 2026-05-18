package intent

import (
	"github.com/snehmatic/mindloop/internal/repository/intent"
	"github.com/snehmatic/mindloop/models"
)

// Service handles the logic for managing user intents
type Service struct {
	repository intent.Repository
}

func NewService(repo intent.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) StartIntent(name string) (*models.Intent, error) {
	return s.repository.StartIntent(name)
}

func (s *Service) ListIntents() ([]models.Intent, error) {
	return s.repository.ListIntents()
}

func (s *Service) ListActiveIntents() ([]models.Intent, error) {
	return s.repository.ListActiveIntents()
}

func (s *Service) GetOngoingIntent() (*models.Intent, error) {
	return s.repository.GetOngoingIntent()
}

func (s *Service) GetIntent(id string) (*models.Intent, error) {
	return s.repository.GetIntent(id)
}

func (s *Service) UpdateIntent(intent *models.Intent) error {
	return s.repository.UpdateIntent(intent)
}

func (s *Service) EndIntent(idStr string, pointsToAward int) (*models.Intent, bool, error) {
	return s.repository.EndIntent(idStr, pointsToAward)
}

func (s *Service) DeleteIntent(id string) error {
	return s.repository.DeleteIntent(id)
}

func (s *Service) DeleteAll() error {
	return s.repository.DeleteAll()
}

func (s *Service) PauseIntent(id uint) (*models.Intent, error) {
	return s.repository.PauseIntent(id)
}

func (s *Service) ResumeIntent(id uint) (*models.Intent, error) {
	return s.repository.ResumeIntent(id)
}
