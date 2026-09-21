package focus

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSession(session *models.FocusSession) error
	ListSessions() ([]models.FocusSession, error)
	GetSession(id int) (*models.FocusSession, error)
	GetSessionByID(id uint) (*models.FocusSession, error)
	UpdateSession(session *models.FocusSession) error
	DeleteSession(id int) error
	DeleteAll() error
	GetActiveSession() (*models.FocusSession, error)

	GetDB() *gorm.DB
}
