package persistence

import (
	"fmt"
	"time"

	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/domain/ports"
	"gorm.io/gorm"
)

// habitRepository implements ports.HabitRepository
type habitRepository struct {
	db *gorm.DB
}

// NewHabitRepository creates a new habit repository
func NewHabitRepository(db *gorm.DB) ports.HabitRepository {
	return &habitRepository{db: db}
}

func (r *habitRepository) Create(habit *entities.Habit) error {
	return r.db.Create(habit).Error
}

func (r *habitRepository) GetByID(id uint) (*entities.Habit, error) {
	var habit entities.Habit
	err := r.db.First(&habit, id).Error
	if err != nil {
		return nil, err
	}
	return &habit, nil
}

func (r *habitRepository) GetAll() ([]*entities.Habit, error) {
	var habits []*entities.Habit
	err := r.db.Find(&habits).Error
	return habits, err
}

func (r *habitRepository) Update(habit *entities.Habit) error {
	return r.db.Save(habit).Error
}

func (r *habitRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Habit{}, id).Error
}

// habitLogRepository implements ports.HabitLogRepository
type habitLogRepository struct {
	db *gorm.DB
}

// NewHabitLogRepository creates a new habit log repository
func NewHabitLogRepository(db *gorm.DB) ports.HabitLogRepository {
	return &habitLogRepository{db: db}
}

func (r *habitLogRepository) Create(habitLog *entities.HabitLog) error {
	return r.db.Create(habitLog).Error
}

func (r *habitLogRepository) GetByID(id uint) (*entities.HabitLog, error) {
	var habitLog entities.HabitLog
	err := r.db.First(&habitLog, id).Error
	if err != nil {
		return nil, err
	}
	return &habitLog, nil
}

func (r *habitLogRepository) GetByHabitID(habitID uint) ([]*entities.HabitLog, error) {
	var habitLogs []*entities.HabitLog
	err := r.db.Where("HabitID = ?", habitID).Order("CreatedAt DESC").Find(&habitLogs).Error
	return habitLogs, err
}

func (r *habitLogRepository) GetByDateRange(start, end time.Time) ([]*entities.HabitLog, error) {
	var habitLogs []*entities.HabitLog
	err := r.db.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Order("CreatedAt DESC").Find(&habitLogs).Error
	return habitLogs, err
}

func (r *habitLogRepository) Update(habitLog *entities.HabitLog) error {
	return r.db.Save(habitLog).Error
}

func (r *habitLogRepository) Delete(id uint) error {
	return r.db.Delete(&entities.HabitLog{}, id).Error
}

// intentRepository implements ports.IntentRepository
type intentRepository struct {
	db *gorm.DB
}

// NewIntentRepository creates a new intent repository
func NewIntentRepository(db *gorm.DB) ports.IntentRepository {
	return &intentRepository{db: db}
}

func (r *intentRepository) Create(intent *entities.Intent) error {
	return r.db.Create(intent).Error
}

func (r *intentRepository) GetByID(id uint) (*entities.Intent, error) {
	var intent entities.Intent
	err := r.db.First(&intent, id).Error
	if err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *intentRepository) GetAll() ([]*entities.Intent, error) {
	var intents []*entities.Intent
	err := r.db.Find(&intents).Error
	return intents, err
}

func (r *intentRepository) GetActive() ([]*entities.Intent, error) {
	var intents []*entities.Intent
	err := r.db.Where("status = ?", entities.IntentActive).Find(&intents).Error
	return intents, err
}

func (r *intentRepository) GetByDateRange(start, end time.Time) ([]*entities.Intent, error) {
	var intents []*entities.Intent
	err := r.db.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Find(&intents).Error
	return intents, err
}

func (r *intentRepository) Update(intent *entities.Intent) error {
	return r.db.Save(intent).Error
}

func (r *intentRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Intent{}, id).Error
}

// focusSessionRepository implements ports.FocusSessionRepository
type focusSessionRepository struct {
	db *gorm.DB
}

// NewFocusSessionRepository creates a new focus session repository
func NewFocusSessionRepository(db *gorm.DB) ports.FocusSessionRepository {
	return &focusSessionRepository{db: db}
}

func (r *focusSessionRepository) Create(session *entities.FocusSession) error {
	return r.db.Create(session).Error
}

func (r *focusSessionRepository) GetByID(id uint) (*entities.FocusSession, error) {
	var session entities.FocusSession
	err := r.db.First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *focusSessionRepository) GetAll() ([]*entities.FocusSession, error) {
	var sessions []*entities.FocusSession
	err := r.db.Find(&sessions).Error
	return sessions, err
}

func (r *focusSessionRepository) GetActive() ([]*entities.FocusSession, error) {
	var sessions []*entities.FocusSession
	err := r.db.Where("status = ?", entities.FocusActive).Find(&sessions).Error
	return sessions, err
}

func (r *focusSessionRepository) GetByDateRange(start, end time.Time) ([]*entities.FocusSession, error) {
	var sessions []*entities.FocusSession
	err := r.db.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Find(&sessions).Error
	return sessions, err
}

func (r *focusSessionRepository) Update(session *entities.FocusSession) error {
	return r.db.Save(session).Error
}

func (r *focusSessionRepository) Delete(id uint) error {
	return r.db.Delete(&entities.FocusSession{}, id).Error
}

// journalRepository implements ports.JournalRepository
type journalRepository struct {
	db *gorm.DB
}

// NewJournalRepository creates a new journal repository
func NewJournalRepository(db *gorm.DB) ports.JournalRepository {
	return &journalRepository{db: db}
}

func (r *journalRepository) Create(entry *entities.JournalEntry) error {
	return r.db.Create(entry).Error
}

func (r *journalRepository) GetByID(id uint) (*entities.JournalEntry, error) {
	var entry entities.JournalEntry
	err := r.db.First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *journalRepository) GetAll() ([]*entities.JournalEntry, error) {
	var entries []*entities.JournalEntry
	err := r.db.Order("CreatedAt DESC").Find(&entries).Error
	return entries, err
}

func (r *journalRepository) GetByDateRange(start, end time.Time) ([]*entities.JournalEntry, error) {
	var entries []*entities.JournalEntry
	err := r.db.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Order("CreatedAt DESC").Find(&entries).Error
	return entries, err
}

func (r *journalRepository) Update(entry *entities.JournalEntry) error {
	return r.db.Save(entry).Error
}

func (r *journalRepository) Delete(id uint) error {
	return r.db.Delete(&entities.JournalEntry{}, id).Error
}

// unitOfWork implements ports.UnitOfWork
type unitOfWork struct {
	db               *gorm.DB
	tx               *gorm.DB
	habitRepo        ports.HabitRepository
	habitLogRepo     ports.HabitLogRepository
	intentRepo       ports.IntentRepository
	focusSessionRepo ports.FocusSessionRepository
	journalRepo      ports.JournalRepository
}

// NewUnitOfWork creates a new unit of work
func NewUnitOfWork(db *gorm.DB) ports.UnitOfWork {
	return &unitOfWork{
		db:               db,
		habitRepo:        NewHabitRepository(db),
		habitLogRepo:     NewHabitLogRepository(db),
		intentRepo:       NewIntentRepository(db),
		focusSessionRepo: NewFocusSessionRepository(db),
		journalRepo:      NewJournalRepository(db),
	}
}

func (uow *unitOfWork) Begin() error {
	if uow.tx != nil {
		return fmt.Errorf("transaction already started")
	}
	uow.tx = uow.db.Begin()
	if uow.tx.Error != nil {
		return uow.tx.Error
	}

	// Update repositories to use transaction
	uow.habitRepo = NewHabitRepository(uow.tx)
	uow.habitLogRepo = NewHabitLogRepository(uow.tx)
	uow.intentRepo = NewIntentRepository(uow.tx)
	uow.focusSessionRepo = NewFocusSessionRepository(uow.tx)
	uow.journalRepo = NewJournalRepository(uow.tx)

	return nil
}

func (uow *unitOfWork) Commit() error {
	if uow.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}
	err := uow.tx.Commit().Error
	uow.tx = nil

	// Reset repositories to use main db
	uow.habitRepo = NewHabitRepository(uow.db)
	uow.habitLogRepo = NewHabitLogRepository(uow.db)
	uow.intentRepo = NewIntentRepository(uow.db)
	uow.focusSessionRepo = NewFocusSessionRepository(uow.db)
	uow.journalRepo = NewJournalRepository(uow.db)

	return err
}

func (uow *unitOfWork) Rollback() error {
	if uow.tx == nil {
		return fmt.Errorf("no transaction to rollback")
	}
	err := uow.tx.Rollback().Error
	uow.tx = nil

	// Reset repositories to use main db
	uow.habitRepo = NewHabitRepository(uow.db)
	uow.habitLogRepo = NewHabitLogRepository(uow.db)
	uow.intentRepo = NewIntentRepository(uow.db)
	uow.focusSessionRepo = NewFocusSessionRepository(uow.db)
	uow.journalRepo = NewJournalRepository(uow.db)

	return err
}

func (uow *unitOfWork) HabitRepository() ports.HabitRepository {
	return uow.habitRepo
}

func (uow *unitOfWork) HabitLogRepository() ports.HabitLogRepository {
	return uow.habitLogRepo
}

func (uow *unitOfWork) IntentRepository() ports.IntentRepository {
	return uow.intentRepo
}

func (uow *unitOfWork) FocusSessionRepository() ports.FocusSessionRepository {
	return uow.focusSessionRepo
}

func (uow *unitOfWork) JournalRepository() ports.JournalRepository {
	return uow.journalRepo
}
