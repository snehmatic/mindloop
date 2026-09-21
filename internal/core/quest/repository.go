package quest

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateQuest(quest *models.SideQuest) error
	ListQuests() ([]models.SideQuest, error)
	GetActiveQuest() (*models.SideQuest, error)
	GetQuest(id uint) (*models.SideQuest, error)
	UpdateQuest(quest *models.SideQuest) error
	DeleteQuest(id uint) error

	GetDB() *gorm.DB
}
