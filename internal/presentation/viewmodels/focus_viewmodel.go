package viewmodels

import (
	"math"

	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/shared/utils"
)

// FocusSessionView represents a focus session for display purposes
type FocusSessionView struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	EndTime   string `json:"end_time"`
	Duration  string `json:"duration"`
	Rating    string `json:"rating"`
	CreatedAt string `json:"created_at"`
}

// ToFocusSessionView converts a domain focus session to view model
func ToFocusSessionView(session *entities.FocusSession) FocusSessionView {
	endTime := "Focus on!"
	if session.IsEnded() {
		endTime = session.EndTime.Format("2006-01-02 15:04:05")
	}

	rating := "Not rated"
	if session.Rating >= 0 {
		rating = string(rune(session.Rating+'0')) + "/10"
	}

	duration := session.GetCurrentDuration()
	if duration < 1 {
		duration = 1 // Show at least 1 minute
	}

	return FocusSessionView{
		ID:        session.ID,
		Title:     session.Title,
		Status:    string(session.Status),
		EndTime:   endTime,
		Duration:  utils.FormatMinutes(math.Floor(duration)),
		Rating:    rating,
		CreatedAt: session.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToFocusSessionViews converts multiple domain focus sessions to view models
func ToFocusSessionViews(sessions []*entities.FocusSession) []FocusSessionView {
	views := make([]FocusSessionView, len(sessions))
	for i, session := range sessions {
		views[i] = ToFocusSessionView(session)
	}
	return views
}
