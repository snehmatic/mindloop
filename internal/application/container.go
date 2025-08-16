package application

import (
	"github.com/snehmatic/mindloop/internal/application/usecases"
	"github.com/snehmatic/mindloop/internal/domain/ports"
	"github.com/snehmatic/mindloop/internal/infrastructure/config"
	"github.com/snehmatic/mindloop/internal/infrastructure/persistence"
	"github.com/snehmatic/mindloop/internal/shared/ui"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	Config *config.Config

	// Database
	DB *persistence.Database

	// Repositories
	HabitRepo        ports.HabitRepository
	HabitLogRepo     ports.HabitLogRepository
	IntentRepo       ports.IntentRepository
	FocusSessionRepo ports.FocusSessionRepository
	JournalRepo      ports.JournalRepository
	UnitOfWork       ports.UnitOfWork

	// Use Cases
	HabitUseCase   usecases.HabitUseCase
	IntentUseCase  usecases.IntentUseCase
	FocusUseCase   usecases.FocusUseCase
	JournalUseCase usecases.JournalUseCase
	SummaryUseCase usecases.SummaryUseCase

	// UI
	UI ui.Interface
}

// NewContainer creates and configures the dependency injection container
func NewContainer(appConfig *config.Config) (*Container, error) {
	// Initialize database
	database, err := persistence.CreateDatabaseFromConfig(appConfig)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	db := database.GetDB()
	habitRepo := persistence.NewHabitRepository(db)
	habitLogRepo := persistence.NewHabitLogRepository(db)
	intentRepo := persistence.NewIntentRepository(db)
	focusSessionRepo := persistence.NewFocusSessionRepository(db)
	journalRepo := persistence.NewJournalRepository(db)
	unitOfWork := persistence.NewUnitOfWork(db)

	// Initialize use cases
	habitUseCase := usecases.NewHabitUseCase(habitRepo, habitLogRepo)
	intentUseCase := usecases.NewIntentUseCase(intentRepo)
	focusUseCase := usecases.NewFocusUseCase(focusSessionRepo)
	journalUseCase := usecases.NewJournalUseCase(journalRepo)
	summaryUseCase := usecases.NewSummaryUseCase(focusSessionRepo, habitRepo, habitLogRepo, intentRepo)

	// Initialize UI
	uiInterface := ui.NewCLIInterface(appConfig.Logger)

	return &Container{
		Config:           appConfig,
		DB:               database,
		HabitRepo:        habitRepo,
		HabitLogRepo:     habitLogRepo,
		IntentRepo:       intentRepo,
		FocusSessionRepo: focusSessionRepo,
		JournalRepo:      journalRepo,
		UnitOfWork:       unitOfWork,
		HabitUseCase:     habitUseCase,
		IntentUseCase:    intentUseCase,
		FocusUseCase:     focusUseCase,
		JournalUseCase:   journalUseCase,
		SummaryUseCase:   summaryUseCase,
		UI:               uiInterface,
	}, nil
}

// Close cleans up resources
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
