package viewmodels

import (
	"github.com/snehmatic/mindloop/internal/domain/entities"
)

// IntentView represents an intent for display purposes
type IntentView struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	EndedAt   string `json:"ended_at"`
}

// ToIntentView converts a domain intent to view model
func ToIntentView(intent *entities.Intent) IntentView {
	var endedAt string
	if intent.EndedAt != nil {
		endedAt = intent.EndedAt.Format("2006-01-02 15:04")
	} else {
		endedAt = "-"
	}

	return IntentView{
		ID:        intent.ID,
		Name:      intent.Name,
		Status:    string(intent.Status),
		CreatedAt: intent.CreatedAt.Format("2006-01-02 15:04"),
		EndedAt:   endedAt,
	}
}

// ToIntentViews converts multiple domain intents to view models
func ToIntentViews(intents []*entities.Intent) []IntentView {
	views := make([]IntentView, len(intents))
	for i, intent := range intents {
		views[i] = ToIntentView(intent)
	}
	return views
}
