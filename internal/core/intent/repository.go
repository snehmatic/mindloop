package intent

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateIntent(intent *models.Intent) error
	ListIntents() ([]models.Intent, error)
	ListActiveIntents() ([]models.Intent, error)
	GetOngoingIntent() (*models.Intent, error)
	GetIntent(id string) (*models.Intent, error)
	GetIntentByID(id uint) (*models.Intent, error)
	UpdateIntent(intent *models.Intent) error
	DeleteIntent(id string) error
	DeleteAll() error

	GetDB() *gorm.DB
}
