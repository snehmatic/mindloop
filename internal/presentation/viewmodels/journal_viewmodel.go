package viewmodels

import (
	"github.com/snehmatic/mindloop/internal/domain/entities"
)

// JournalEntryView represents a journal entry for display purposes
type JournalEntryView struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Mood      string `json:"mood"`
	CreatedAt string `json:"created_at"`
	Preview   string `json:"preview"` // First 100 characters of content
}

// JournalEntryDetailView represents a detailed journal entry for display purposes
type JournalEntryDetailView struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Mood      string `json:"mood"`
	CreatedAt string `json:"created_at"`
}

// ToJournalEntryView converts a domain journal entry to view model
func ToJournalEntryView(entry *entities.JournalEntry) JournalEntryView {
	preview := entry.Content
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}

	return JournalEntryView{
		ID:        entry.ID,
		Title:     entry.Title,
		Mood:      string(entry.Mood),
		CreatedAt: entry.CreatedAt.Format("2006-01-02 15:04:05"),
		Preview:   preview,
	}
}

// ToJournalEntryDetailView converts a domain journal entry to detailed view model
func ToJournalEntryDetailView(entry *entities.JournalEntry) JournalEntryDetailView {
	return JournalEntryDetailView{
		ID:        entry.ID,
		Title:     entry.Title,
		Content:   entry.Content,
		Mood:      string(entry.Mood),
		CreatedAt: entry.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToJournalEntryViews converts multiple domain journal entries to view models
func ToJournalEntryViews(entries []*entities.JournalEntry) []JournalEntryView {
	views := make([]JournalEntryView, len(entries))
	for i, entry := range entries {
		views[i] = ToJournalEntryView(entry)
	}
	return views
}
