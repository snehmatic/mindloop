package usecases

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/domain/ports"
)

// IntentUseCase defines the interface for intent use cases
type IntentUseCase interface {
	StartIntent(name string) (*entities.Intent, error)
	EndIntent(id uint) (*entities.Intent, error)
	GetIntent(id uint) (*entities.Intent, error)
	GetAllIntents() ([]*entities.Intent, error)
	GetActiveIntents() ([]*entities.Intent, error)
	DeleteIntent(id uint) error
}

type intentUseCase struct {
	intentRepo ports.IntentRepository
}

// NewIntentUseCase creates a new intent use case
func NewIntentUseCase(intentRepo ports.IntentRepository) IntentUseCase {
	return &intentUseCase{
		intentRepo: intentRepo,
	}
}

// StartIntent creates and starts a new intent
func (i *intentUseCase) StartIntent(name string) (*entities.Intent, error) {
	intent := entities.NewIntent(name)

	if err := intent.Validate(); err != nil {
		return nil, fmt.Errorf("intent validation failed: %w", err)
	}

	if err := i.intentRepo.Create(intent); err != nil {
		return nil, fmt.Errorf("failed to start intent: %w", err)
	}

	return intent, nil
}

// EndIntent ends an existing intent
func (i *intentUseCase) EndIntent(id uint) (*entities.Intent, error) {
	intent, err := i.intentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("intent not found: %w", err)
	}

	if intent.IsDone() {
		return nil, fmt.Errorf("intent is already done")
	}

	intent.End()

	if err := i.intentRepo.Update(intent); err != nil {
		return nil, fmt.Errorf("failed to end intent: %w", err)
	}

	return intent, nil
}

// GetIntent retrieves an intent by ID
func (i *intentUseCase) GetIntent(id uint) (*entities.Intent, error) {
	intent, err := i.intentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get intent: %w", err)
	}
	return intent, nil
}

// GetAllIntents retrieves all intents
func (i *intentUseCase) GetAllIntents() ([]*entities.Intent, error) {
	intents, err := i.intentRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get intents: %w", err)
	}
	return intents, nil
}

// GetActiveIntents retrieves all active intents
func (i *intentUseCase) GetActiveIntents() ([]*entities.Intent, error) {
	intents, err := i.intentRepo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get active intents: %w", err)
	}
	return intents, nil
}

// DeleteIntent deletes an intent
func (i *intentUseCase) DeleteIntent(id uint) error {
	// Check if intent exists
	_, err := i.intentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("intent not found: %w", err)
	}

	if err := i.intentRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete intent: %w", err)
	}

	return nil
}
