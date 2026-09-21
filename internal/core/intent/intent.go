package intent

import (
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/nlp"
	"github.com/snehmatic/mindloop/models"
)

// Service handles the logic for managing user intents
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) StartIntent(name string) (*models.Intent, error) {
	if name == "" {
		return nil, ErrNameCannotBeEmpty
	}

	cleanedName, dueDate := nlp.ExtractDate(name)

	intent := &models.Intent{
		Name:    cleanedName,
		Status:  models.StatusActive,
		DueDate: dueDate,
	}

	if err := s.repo.CreateIntent(intent); err != nil {
		return nil, err
	}
	return intent, nil
}

func (s *Service) ListIntents() ([]models.Intent, error) {
	return s.repo.ListIntents()
}

func (s *Service) ListActiveIntents() ([]models.Intent, error) {
	return s.repo.ListActiveIntents()
}

func (s *Service) GetOngoingIntent() (*models.Intent, error) {
	return s.repo.GetOngoingIntent()
}

func (s *Service) GetIntent(id string) (*models.Intent, error) {
	return s.repo.GetIntent(id)
}

func (s *Service) UpdateIntent(intent *models.Intent) error {
	return s.repo.UpdateIntent(intent)
}

func (s *Service) EndIntent(idStr string, pointsToAward int) (*models.Intent, bool, error) {
	intent, err := s.repo.GetIntent(idStr)
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	intent.Status = models.StatusDone
	intent.EndedAt = &now

	if err := s.repo.UpdateIntent(intent); err != nil {
		return nil, false, err
	}

	milestoneReached, _ := points.AwardPoints(s.repo.GetDB(), models.CategoryIntent, intent.ID, pointsToAward)

	return intent, milestoneReached, nil
}

func (s *Service) DeleteIntent(id string) error {
	return s.repo.DeleteIntent(id)
}

func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *Service) PauseIntent(id uint) (*models.Intent, error) {
	intent, err := s.repo.GetIntentByID(id)
	if err != nil {
		return nil, err
	}

	if intent.Status != models.StatusActive {
		return nil, ErrIntentIsNotActive
	}

	intent.Status = models.StatusPaused
	if err := s.repo.UpdateIntent(intent); err != nil {
		return nil, err
	}
	return intent, nil
}

func (s *Service) ResumeIntent(id uint) (*models.Intent, error) {
	intent, err := s.repo.GetIntentByID(id)
	if err != nil {
		return nil, err
	}

	if intent.Status != models.StatusPaused {
		return nil, ErrIntentIsNotPaused
	}

	intent.Status = models.StatusActive
	if err := s.repo.UpdateIntent(intent); err != nil {
		return nil, err
	}
	return intent, nil
}
