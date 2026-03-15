package habit_test

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/habit"
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

	err = db.AutoMigrate(&models.Habit{}, &models.HabitLog{}, &models.PointTransaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestHabitService(t *testing.T) {
	db := setupTestDB(t)
	s := habit.NewService(db)

	// 1. Create Habit
	h := &models.Habit{
		Title:       "Read",
		TargetCount: 1,
		Interval:    models.Daily,
	}
	err := s.CreateHabit(h)
	if err != nil {
		t.Fatalf("Failed to create habit: %v", err)
	}

	// 2. Log Habit
	_, log, _, err := s.LogHabit("1", 5)
	if err != nil {
		t.Fatalf("Failed to log habit: %v", err)
	}
	if log.ActualCount != 1 {
		t.Errorf("Expected actual count 1, got %d", log.ActualCount)
	}

	// 3. Log Habit again (already completed)
	_, _, _, err = s.LogHabit("1", 5)
	if err == nil {
		t.Error("Expected error when logging already completed habit, got nil")
	}

	// 4. Calculate Streak
	streak, err := s.CalculateStreak(h.ID, models.Daily)
	if err != nil {
		t.Fatalf("Failed to calculate streak: %v", err)
	}
	if streak != 1 {
		t.Errorf("Expected streak 1, got %d", streak)
	}

	// 5. Unlog Habit
	_, err = s.UnlogHabit("1")
	if err != nil {
		t.Fatalf("Failed to unlog habit: %v", err)
	}

	// 6. Delete Habit
	err = s.DeleteHabit("1")
	if err != nil {
		t.Fatalf("Failed to delete habit: %v", err)
	}
}

func TestCalculateStreak(t *testing.T) {
	db := setupTestDB(t)
	s := habit.NewService(db)

	h := &models.Habit{Title: "Run", TargetCount: 1, Interval: models.Daily}
	if err := s.CreateHabit(h); err != nil {
		t.Fatalf("Failed to create habit: %v", err)
	}

	// Create manual logs for 3 consecutive days (including today)
	today := time.Now().Truncate(24 * time.Hour)
	for i := 0; i < 3; i++ {
		date := today.AddDate(0, 0, -i)
		db.Create(&models.HabitLog{
			HabitID:     h.ID,
			Title:       h.Title,
			Interval:    models.Daily,
			TargetCount: 1,
			ActualCount: 1,
			EndedAt:     date,
			Model:       gorm.Model{CreatedAt: date},
		})
	}

	streak, _ := s.CalculateStreak(h.ID, models.Daily)
	if streak != 3 {
		t.Errorf("Expected streak 3, got %d", streak)
	}

	// Add a gap
	gapDate := today.AddDate(0, 0, -4) // Day 3 is missing
	db.Create(&models.HabitLog{
		HabitID:     h.ID,
		Title:       h.Title,
		Interval:    models.Daily,
		TargetCount: 1,
		ActualCount: 1,
		EndedAt:     gapDate,
		Model:       gorm.Model{CreatedAt: gapDate},
	})

	streak, _ = s.CalculateStreak(h.ID, models.Daily)
	if streak != 3 {
		t.Errorf("Expected streak 3 after gap, got %d", streak)
	}
}
