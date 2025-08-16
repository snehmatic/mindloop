package usecases

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/domain/ports"
)

// JournalUseCase defines the interface for journal use cases
type JournalUseCase interface {
	CreateEntry(title, content string, mood entities.Mood) (*entities.JournalEntry, error)
	GetEntry(id uint) (*entities.JournalEntry, error)
	GetAllEntries() ([]*entities.JournalEntry, error)
	UpdateEntry(entry *entities.JournalEntry) error
	DeleteEntry(id uint) error
}

type journalUseCase struct {
	journalRepo ports.JournalRepository
}

// NewJournalUseCase creates a new journal use case
func NewJournalUseCase(journalRepo ports.JournalRepository) JournalUseCase {
	return &journalUseCase{
		journalRepo: journalRepo,
	}
}

// CreateEntry creates a new journal entry
func (j *journalUseCase) CreateEntry(title, content string, mood entities.Mood) (*entities.JournalEntry, error) {
	entry := entities.NewJournalEntry(title, content, mood)

	if err := entry.Validate(); err != nil {
		return nil, fmt.Errorf("journal entry validation failed: %w", err)
	}

	if err := j.journalRepo.Create(entry); err != nil {
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}

	return entry, nil
}

// GetEntry retrieves a journal entry by ID
func (j *journalUseCase) GetEntry(id uint) (*entities.JournalEntry, error) {
	entry, err := j.journalRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entry: %w", err)
	}
	return entry, nil
}

// GetAllEntries retrieves all journal entries
func (j *journalUseCase) GetAllEntries() ([]*entities.JournalEntry, error) {
	entries, err := j.journalRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entries: %w", err)
	}
	return entries, nil
}

// UpdateEntry updates an existing journal entry
func (j *journalUseCase) UpdateEntry(entry *entities.JournalEntry) error {
	if err := entry.Validate(); err != nil {
		return fmt.Errorf("journal entry validation failed: %w", err)
	}

	if err := j.journalRepo.Update(entry); err != nil {
		return fmt.Errorf("failed to update journal entry: %w", err)
	}

	return nil
}

// DeleteEntry deletes a journal entry
func (j *journalUseCase) DeleteEntry(id uint) error {
	// Check if entry exists
	_, err := j.journalRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("journal entry not found: %w", err)
	}

	if err := j.journalRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete journal entry: %w", err)
	}

	return nil
}
