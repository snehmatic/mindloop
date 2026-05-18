package focus

import (
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// Repository defines the interface for focus session data access
type Repository interface {
	StartSession(title string) (*models.FocusSession, error)
	ListSessions() ([]models.FocusSession, error)
	GetSession(id int) (*models.FocusSession, error)
	UpdateSession(session *models.FocusSession) error
	EndSession(id int, pointsToAward int) (*models.FocusSession, bool, error)
	RateSession(id int, rating int) (*models.FocusSession, error)
	DeleteSession(id int) error
	DeleteAll() error
	PauseSession(id uint) (*models.FocusSession, error)
	ResumeSession(id uint) (*models.FocusSession, error)
	GetActiveSession() (*models.FocusSession, error)
	GetSessionsInRange(start, end time.Time) ([]models.FocusSession, error)
	GetFocusStats(start, end time.Time) (models.FocusStats, error)
	GetDB() *gorm.DB
}
