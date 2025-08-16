package usecases

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/domain/ports"
)

// FocusUseCase defines the interface for focus session use cases
type FocusUseCase interface {
	StartFocusSession(title string) (*entities.FocusSession, error)
	EndFocusSession(id uint) (*entities.FocusSession, error)
	PauseFocusSession(id uint) (*entities.FocusSession, error)
	ResumeFocusSession(id uint) (*entities.FocusSession, error)
	RateFocusSession(id uint, rating int) (*entities.FocusSession, error)
	GetFocusSession(id uint) (*entities.FocusSession, error)
	GetAllFocusSessions() ([]*entities.FocusSession, error)
	GetActiveFocusSessions() ([]*entities.FocusSession, error)
	DeleteFocusSession(id uint) error
}

type focusUseCase struct {
	focusRepo ports.FocusSessionRepository
}

// NewFocusUseCase creates a new focus use case
func NewFocusUseCase(focusRepo ports.FocusSessionRepository) FocusUseCase {
	return &focusUseCase{
		focusRepo: focusRepo,
	}
}

// StartFocusSession creates and starts a new focus session
func (f *focusUseCase) StartFocusSession(title string) (*entities.FocusSession, error) {
	session := entities.NewFocusSession(title)

	if err := session.Validate(); err != nil {
		return nil, fmt.Errorf("focus session validation failed: %w", err)
	}

	if err := f.focusRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to start focus session: %w", err)
	}

	return session, nil
}

// EndFocusSession ends an active focus session
func (f *focusUseCase) EndFocusSession(id uint) (*entities.FocusSession, error) {
	session, err := f.focusRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("focus session not found: %w", err)
	}

	if !session.IsActive() {
		return nil, fmt.Errorf("focus session is not active")
	}

	session.End()

	if err := f.focusRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to end focus session: %w", err)
	}

	return session, nil
}

// PauseFocusSession pauses an active focus session
func (f *focusUseCase) PauseFocusSession(id uint) (*entities.FocusSession, error) {
	session, err := f.focusRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("focus session not found: %w", err)
	}

	if !session.IsActive() {
		return nil, fmt.Errorf("focus session is not active")
	}

	session.Pause()

	if err := f.focusRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to pause focus session: %w", err)
	}

	return session, nil
}

// ResumeFocusSession resumes a paused focus session
func (f *focusUseCase) ResumeFocusSession(id uint) (*entities.FocusSession, error) {
	session, err := f.focusRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("focus session not found: %w", err)
	}

	session.Resume()

	if err := f.focusRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to resume focus session: %w", err)
	}

	return session, nil
}

// RateFocusSession rates a completed focus session
func (f *focusUseCase) RateFocusSession(id uint, rating int) (*entities.FocusSession, error) {
	session, err := f.focusRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("focus session not found: %w", err)
	}

	if err := session.Rate(rating); err != nil {
		return nil, fmt.Errorf("failed to rate focus session: %w", err)
	}

	if err := f.focusRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to save focus session rating: %w", err)
	}

	return session, nil
}

// GetFocusSession retrieves a focus session by ID
func (f *focusUseCase) GetFocusSession(id uint) (*entities.FocusSession, error) {
	session, err := f.focusRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get focus session: %w", err)
	}
	return session, nil
}

// GetAllFocusSessions retrieves all focus sessions
func (f *focusUseCase) GetAllFocusSessions() ([]*entities.FocusSession, error) {
	sessions, err := f.focusRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get focus sessions: %w", err)
	}
	return sessions, nil
}

// GetActiveFocusSessions retrieves all active focus sessions
func (f *focusUseCase) GetActiveFocusSessions() ([]*entities.FocusSession, error) {
	sessions, err := f.focusRepo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get active focus sessions: %w", err)
	}
	return sessions, nil
}

// DeleteFocusSession deletes a focus session
func (f *focusUseCase) DeleteFocusSession(id uint) error {
	// Check if session exists
	_, err := f.focusRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("focus session not found: %w", err)
	}

	if err := f.focusRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete focus session: %w", err)
	}

	return nil
}
