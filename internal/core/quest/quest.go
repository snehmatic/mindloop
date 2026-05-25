package quest

import (
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/repository/quest"
	"github.com/snehmatic/mindloop/models"
)

// Service handles the logic for managing side quests
type Service struct {
	repository quest.Repository
	logger     log.Logger
}

// NewService creates a new quest Service instance
func NewService(repo quest.Repository, logger log.Logger) *Service {
	return &Service{
		repository: repo,
		logger:     logger,
	}
}

func (s *Service) StartQuest(title string) (*models.SideQuest, error) {
	return s.repository.StartQuest(title)
}

func (s *Service) StopQuest(id uint, note string, pointsToAward int) (*models.SideQuest, bool, error) {
	return s.repository.StopQuest(id, note, pointsToAward)
}

func (s *Service) ListQuests() ([]models.SideQuest, error) {
	return s.repository.ListQuests()
}

func (s *Service) GetActiveQuest() (*models.SideQuest, error) {
	return s.repository.GetActiveQuest()
}

func (s *Service) DeleteQuest(id uint) error {
	return s.repository.DeleteQuest(id)
}
