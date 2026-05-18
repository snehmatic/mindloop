package intent

import (
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Repository defines the interface for intent data access
type Repository interface {
	StartIntent(name string) (*models.Intent, error)
	ListIntents() ([]models.Intent, error)
	ListActiveIntents() ([]models.Intent, error)
	GetOngoingIntent() (*models.Intent, error)
	GetIntent(id string) (*models.Intent, error)
	UpdateIntent(intent *models.Intent) error
	EndIntent(idStr string, pointsToAward int) (*models.Intent, bool, error)
	DeleteIntent(id string) error
	DeleteAll() error
	PauseIntent(id uint) (*models.Intent, error)
	ResumeIntent(id uint) (*models.Intent, error)
	GetIntentStats(start, end time.Time) ([]models.IntentStats, error)
	GetIntentsInRange(start, end time.Time) ([]models.Intent, error)
	GetDB() *gorm.DB
}
