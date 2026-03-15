package summary_test

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	err = db.AutoMigrate(
		&models.Intent{},
		&models.FocusSession{},
		&models.Habit{},
		&models.HabitLog{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestSummaryService(t *testing.T) {
	db := setupTestDB(t)
	s := summary.NewService(db)

	now := time.Now()
	start := now.AddDate(0, 0, -1)
	end := now.AddDate(0, 0, 1)

	// 1. Seed Focus Session
	db.Create(&models.FocusSession{
		Title:    "Session 1",
		Duration: 30,
		Status:   "ended",
		Model:    gorm.Model{CreatedAt: now},
	})

	// 2. Seed Habit and Log
	h := models.Habit{Title: "Habit 1", TargetCount: 1, Interval: models.Daily}
	db.Create(&h)
	db.Create(&models.HabitLog{
		HabitID:     h.ID,
		Title:       h.Title,
		TargetCount: 1,
		ActualCount: 1,
		Model:       gorm.Model{CreatedAt: now},
	})

	// 3. Seed Intent
	db.Create(&models.Intent{
		Name:   "Intent 1",
		Status: "done",
		Model:  gorm.Model{CreatedAt: now},
	})

	// 4. Generate Summary
	report, err := s.GenerateSummary(start, end)
	if err != nil {
		t.Fatalf("GenerateSummary failed: %v", err)
	}

	if report.Focus.TotalSessions != 1 {
		t.Errorf("Expected 1 focus session, got %d", report.Focus.TotalSessions)
	}
	if report.Focus.RawDuration != 30 {
		t.Errorf("Expected 30 mins focus, got %f", report.Focus.RawDuration)
	}
	if len(report.Habits) != 1 {
		t.Errorf("Expected 1 habit in stats, got %d", len(report.Habits))
	}
	if report.Habits[0].CompletionRate != 100 {
		t.Errorf("Expected 100%% completion, got %f", report.Habits[0].CompletionRate)
	}
	if len(report.Intents) != 1 {
		t.Errorf("Expected 1 intent in stats, got %d", len(report.Intents))
	}
}

func TestGetFocusSeries(t *testing.T) {
	db := setupTestDB(t)
	s := summary.NewService(db)

	today := time.Now().Truncate(24 * time.Hour)
	start := today.AddDate(0, 0, -2) // 3 days total
	end := today.Add(24*time.Hour - 1*time.Second)

	// Day 0: 30 mins
	db.Create(&models.FocusSession{Duration: 30, Model: gorm.Model{CreatedAt: start.Add(1 * time.Hour)}})
	// Day 1: 0 mins
	// Day 2 (today): 60 mins
	db.Create(&models.FocusSession{Duration: 60, Model: gorm.Model{CreatedAt: today.Add(1 * time.Hour)}})

	stats, labels, err := s.GetFocusSeries(start, end)
	if err != nil {
		t.Fatalf("GetFocusSeries failed: %v", err)
	}

	if len(stats) != 3 {
		t.Errorf("Expected 3 days of stats, got %d", len(stats))
	}
	if stats[0] != 30 {
		t.Errorf("Expected 30 mins on day 0, got %f", stats[0])
	}
	if stats[1] != 0 {
		t.Errorf("Expected 0 mins on day 1, got %f", stats[1])
	}
	if stats[2] != 60 {
		t.Errorf("Expected 60 mins on day 2, got %f", stats[2])
	}
	if len(labels) != 3 {
		t.Errorf("Expected 3 labels, got %d", len(labels))
	}
}
