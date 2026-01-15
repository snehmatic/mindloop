package v1

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCleanSlateTest(t *testing.T) (*MindloopHandler, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to memory db: %v", err)
	}

	// AutoMigrate all models
	err = db.AutoMigrate(
		&models.JournalEntry{},
		&models.Habit{},
		&models.HabitLog{},
		&models.FocusSession{},
		&models.Intent{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	hService := habit.NewService(db)
	jService := journal.NewService(db)
	fService := focus.NewService(db)
	iService := intent.NewService(db)
	sService := summary.NewService(db)

	mlh := NewMindloopHandler(jService, hService, fService, iService, sService)
	return mlh, db
}

func TestCleanSlate(t *testing.T) {
	mlh, db := setupCleanSlateTest(t)

	// 1. Seed Data
	db.Create(&models.JournalEntry{Title: "ToDelete", Content: "Content", Mood: "happy"})
	db.Create(&models.Habit{Title: "ToDelete", TargetCount: 1, Interval: "daily"})

	// Verify Seeding
	var jCount, hCount int64
	db.Model(&models.JournalEntry{}).Count(&jCount)
	db.Model(&models.Habit{}).Count(&hCount)

	if jCount != 1 || hCount != 1 {
		t.Fatalf("Setup failed: expected 1 journal and 1 habit, got %d and %d", jCount, hCount)
	}

	// 2. Perform Clean Slate Request (Type=all)
	data := url.Values{}
	data.Set("type", "all")
	req, _ := http.NewRequest("POST", "/cleanslate", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	mlh.HandleCleanSlate(rr, req)

	// 3. Verify Result
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect 303, got %d", rr.Code)
	}

	// Check DB counts
	db.Model(&models.JournalEntry{}).Count(&jCount)
	db.Model(&models.Habit{}).Count(&hCount)

	if jCount != 0 {
		t.Errorf("Clean Slate failed: expected 0 journals, got %d", jCount)
	}
	if hCount != 0 {
		t.Errorf("Clean Slate failed: expected 0 habits, got %d", hCount)
	}
}

func TestCleanSlateJournalOnly(t *testing.T) {
	mlh, db := setupCleanSlateTest(t)

	// 1. Seed Data
	db.Create(&models.JournalEntry{Title: "ToDelete", Content: "Content"})
	db.Create(&models.Habit{Title: "KeepMe", TargetCount: 1})

	// 2. Perform Clean Slate (Type=journal)
	data := url.Values{}
	data.Set("type", "journal")
	req, _ := http.NewRequest("POST", "/cleanslate", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	mlh.HandleCleanSlate(rr, req)

	// 3. Verify
	var jCount, hCount int64
	db.Model(&models.JournalEntry{}).Count(&jCount)
	db.Model(&models.Habit{}).Count(&hCount)

	if jCount != 0 {
		t.Errorf("Expected 0 journals, got %d", jCount)
	}
	if hCount != 1 {
		t.Errorf("Expected 1 habit remaining, got %d", hCount)
	}
}
