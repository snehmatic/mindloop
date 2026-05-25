package quest

import (
	"github.com/snehmatic/mindloop/models"
)

// Repository defines the interface for side quest data access
type Repository interface {
	StartQuest(title string) (*models.SideQuest, error)
	StopQuest(id uint, note string, pointsToAward int) (*models.SideQuest, bool, error)
	ListQuests() ([]models.SideQuest, error)
	GetActiveQuest() (*models.SideQuest, error)
	DeleteQuest(id uint) error
}
