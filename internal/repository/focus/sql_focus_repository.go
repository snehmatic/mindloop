package focus

import (
	"errors"
	"time"

	"github.com/snehmatic/mindloop/internal/core/points"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based focus session repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// StartSession creates a new focus session
func (r *sqlRepository) StartSession(title string) (*models.FocusSession, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	var activeSessions []models.FocusSession
	if err := r.DB.Where("status = ?", "active").Limit(1).Find(&activeSessions).Error; err != nil {
		return nil, err
	}
	if len(activeSessions) > 0 {
		return nil, errors.New("a focus session is already active")
	}

	session := &models.FocusSession{
		Title:  title,
		Status: "active",
	}

	if err := r.DB.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

// ListSessions retrieves all focus sessions from the database
func (r *sqlRepository) ListSessions() ([]models.FocusSession, error) {
	var sessions []models.FocusSession
	result := r.DB.Order("CreatedAt DESC").Find(&sessions)
	return sessions, result.Error
}

// GetSession retrieves a single focus session by its ID
func (r *sqlRepository) GetSession(id int) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := r.DB.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateSession modifies an existing focus session in the database
func (r *sqlRepository) UpdateSession(session *models.FocusSession) error {
	return r.DB.Save(session).Error
}

// EndSession ends a focus session and awards points if successful
func (r *sqlRepository) EndSession(id int, pointsToAward int) (*models.FocusSession, bool, error) {
	var session models.FocusSession
	if err := r.DB.First(&session, id).Error; err != nil {
		return nil, false, err
	}

	if session.Status != "active" {
		return nil, false, errors.New("focus session is not active")
	}

	session.Status = "ended"
	session.EndTime = time.Now()
	session.Duration = session.EndTime.Sub(session.CreatedAt).Minutes()

	if err := r.DB.Save(&session).Error; err != nil {
		return nil, false, err
	}

	var milestoneReached bool
	var totalPoints int
	if err := r.DB.Model(&models.PointTransaction{}).Select("COALESCE(SUM(Points), 0)").Scan(&totalPoints).Error; err == nil {
		tx := models.PointTransaction{
			ActivityType: models.CategoryFocus,
			ActivityID:   session.ID,
			Points:       pointsToAward,
		}
		if r.DB.Create(&tx).Error == nil {
			newTotal := totalPoints + pointsToAward
			if newTotal/points.MilestoneInterval > totalPoints/points.MilestoneInterval {
				milestoneReached = true
			}
		}
	}

	return &session, milestoneReached, nil
}

// RateSession sets a rating for a focus session
func (r *sqlRepository) RateSession(id int, rating int) (*models.FocusSession, error) {
	if rating < 0 || rating > 10 {
		return nil, errors.New("rating must be between 0 and 10")
	}

	var session models.FocusSession
	if err := r.DB.First(&session, id).Error; err != nil {
		return nil, err
	}

	if session.Status != "ended" {
		return nil, errors.New("focus session is not ended")
	}

	session.Rating = rating
	if err := r.DB.Save(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession removes a focus session from the database
func (r *sqlRepository) DeleteSession(id int) error {
	r.DB.Model(&models.Task{}).Where("FocusSessionID = ?", id).Update("FocusSessionID", nil)
	return r.DB.Delete(&models.FocusSession{}, id).Error
}

// DeleteAll removes all focus sessions from the database
func (r *sqlRepository) DeleteAll() error {
	return r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.FocusSession{}).Error
}

// PauseSession pauses a focus session
func (r *sqlRepository) PauseSession(id uint) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := r.DB.First(&session, id).Error; err != nil {
		return nil, err
	}

	if session.Status != "active" {
		return nil, errors.New("focus session is not active")
	}

	session.Status = "paused"
	if err := r.DB.Save(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// ResumeSession resumes a focus session
func (r *sqlRepository) ResumeSession(id uint) (*models.FocusSession, error) {
	var session models.FocusSession
	if err := r.DB.First(&session, id).Error; err != nil {
		return nil, err
	}

	if session.Status != "paused" {
		return nil, errors.New("focus session is not paused")
	}

	session.Status = "active"
	if err := r.DB.Save(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetActiveSession retrieves the currently active focus session
func (r *sqlRepository) GetActiveSession() (*models.FocusSession, error) {
	var sessions []models.FocusSession
	err := r.DB.Where("status = ?", "active").Limit(1).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}
	return &sessions[0], nil
}

// GetDB returns the GORM DB instance
func (r *sqlRepository) GetDB() *gorm.DB {
	return r.DB
}

// GetSessionsInRange retrieves all focus sessions within the given time range
func (r *sqlRepository) GetSessionsInRange(start, end time.Time) ([]models.FocusSession, error) {
	var sessions []models.FocusSession
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"
	if err := r.DB.Where(rangeQuery, start, end).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetFocusStats returns focus statistics for the given time range
func (r *sqlRepository) GetFocusStats(start, end time.Time) (models.FocusStats, error) {
	var sessions []models.FocusSession
	rangeQuery := "CreatedAt >= ? AND CreatedAt <= ?"

	if err := r.DB.Where(rangeQuery, start, end).Find(&sessions).Error; err != nil {
		return models.FocusStats{}, err
	}
	if len(sessions) == 0 {
		return models.FocusStats{
			TotalSessions:  0,
			TotalDuration:  "0 mins",
			LongestSession: "0 mins",
		}, nil
	}
	totalDuration := 0.0
	longestSession := 0.0
	for _, session := range sessions {
		totalDuration += session.Duration
		if session.Duration > longestSession {
			longestSession = session.Duration
		}
	}
	return models.FocusStats{
		TotalSessions:  len(sessions),
		TotalDuration:  utils.FormatMinutes(totalDuration),
		RawDuration:    totalDuration,
		LongestSession: utils.FormatMinutes(longestSession),
	}, nil
}
