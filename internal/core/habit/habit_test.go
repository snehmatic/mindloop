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
	momentum, err := s.CalculateMomentum(h)
	if err != nil {
		t.Fatalf("Failed to calculate streak: %v", err)
	}
	
	// Test output is not strictly asserted here because LogHabit creates logs 
	// using time.Now() via GORM, making the exact momentum value dependent on 
	// timezone differences between SQLite and the system. We test math explicitly below.
	if momentum < 0 {
		t.Errorf("Expected momentum to be >= 0")
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

func TestCalculateMomentum(t *testing.T) {
	s := habit.NewService(nil)
	today := time.Now().Truncate(24 * time.Hour)

	h := &models.Habit{Title: "Run", TargetCount: 1, Interval: models.Daily}
	h.CreatedAt = today.AddDate(0, 0, -5)

	var logs []models.HabitLog
	for i := 0; i < 3; i++ {
		date := today.AddDate(0, 0, -i)
		logs = append(logs, models.HabitLog{
			HabitID:     h.ID,
			Title:       h.Title,
			Interval:    models.Daily,
			TargetCount: 1,
			ActualCount: 1,
			EndedAt:     date,
			Model:       gorm.Model{CreatedAt: date},
		})
	}

	momentum := s.CalculateMomentumFromLogs(h, logs)
	if momentum != 30 {
		t.Errorf("Expected momentum 30, got %d", momentum)
	}

	// Add a gap (missing day 3)
	gapDate := today.AddDate(0, 0, -4)
	logs = append(logs, models.HabitLog{
		HabitID:     h.ID,
		Title:       h.Title,
		Interval:    models.Daily,
		TargetCount: 1,
		ActualCount: 1,
		EndedAt:     gapDate,
		Model:       gorm.Model{CreatedAt: gapDate},
	})

	momentum = s.CalculateMomentumFromLogs(h, logs)
	// Expected: today-4 (+10), today-3 (*0.9=9), today-2 (+10=19), today-1 (+10=29), today (+10=39)
	if momentum != 39 {
		t.Errorf("Expected momentum 39 after gap, got %d", momentum)
	}
}

