package quest

import (
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) StartQuest(title string) (*models.SideQuest, error) {
	if title == "" {
		return nil, ErrTitleCannotBeEmpty
	}

	activeQuest, err := s.repo.GetActiveQuest()
	if err != nil {
		return nil, err
	}
	if activeQuest != nil {
		return nil, ErrASideQuestIsAlreadyActive
	}

	quest := &models.SideQuest{
		Title:  title,
		Status: models.StatusActive,
	}

	if err := s.repo.CreateQuest(quest); err != nil {
		return nil, err
	}
	return quest, nil
}

func (s *Service) StopQuest(id uint, note string, pointsToAward int) (*models.SideQuest, bool, error) {
	quest, err := s.repo.GetQuest(id)
	if err != nil {
		return nil, false, err
	}

	if quest.Status != models.StatusActive {
		return nil, false, ErrSideQuestIsNotActive
	}

	quest.Status = models.StatusDone
	quest.Note = note
	now := time.Now()
	quest.EndedAt = &now

	if err := s.repo.UpdateQuest(quest); err != nil {
		return nil, false, err
	}

	milestoneReached, _ := points.AwardPoints(s.repo.GetDB(), models.CategoryQuest, quest.ID, pointsToAward)

	return quest, milestoneReached, nil
}

func (s *Service) ListQuests() ([]models.SideQuest, error) {
	return s.repo.ListQuests()
}

func (s *Service) GetActiveQuest() (*models.SideQuest, error) {
	return s.repo.GetActiveQuest()
}

func (s *Service) DeleteQuest(id uint) error {
	return s.repo.DeleteQuest(id)
}
