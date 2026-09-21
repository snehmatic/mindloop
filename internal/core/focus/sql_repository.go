package focus

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *SQLRepository) CreateSession(session *models.FocusSession) error {
	return r.db.Create(session).Error
}

func (r *SQLRepository) ListSessions() ([]models.FocusSession, error) {
	var sessions []models.FocusSession
	result := r.db.Order("CreatedAt DESC").Find(&sessions)
	return sessions, result.Error
}

func (r *SQLRepository) GetSession(id int) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := r.db.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SQLRepository) GetSessionByID(id uint) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := r.db.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SQLRepository) UpdateSession(session *models.FocusSession) error {
	return r.db.Save(session).Error
}

func (r *SQLRepository) DeleteSession(id int) error {
	r.db.Model(&models.Task{}).Where("FocusSessionID = ?", id).Update("FocusSessionID", nil)
	return r.db.Delete(&models.FocusSession{}, id).Error
}

func (r *SQLRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.FocusSession{}).Error
}

func (r *SQLRepository) GetActiveSession() (*models.FocusSession, error) {
	var sessions []models.FocusSession
	err := r.db.Where("status = ?", models.StatusActive).Limit(1).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}
	return &sessions[0], nil
}
