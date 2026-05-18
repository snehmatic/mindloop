package points

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	pointRepo "github.com/snehmatic/mindloop/internal/repository/point"
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
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.PointTransaction{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestAwardPoints(t *testing.T) {
	db := setupTestDB(t)
	repo := pointRepo.NewSQLRepository(db)
	svc := NewService(repo)

	milestoneReached, err := svc.AwardPoints(models.CategoryHabit, 1, 5)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if milestoneReached {
		t.Errorf("Expected milestone false for initial points")
	}

	var transaction models.PointTransaction
	if err := db.First(&transaction).Error; err != nil {
		t.Fatalf("Failed to retrieve transaction: %v", err)
	}

	if transaction.ActivityType != models.CategoryHabit {
		t.Errorf("Expected ActivityType %s, got %s", models.CategoryHabit, transaction.ActivityType)
	}
	if transaction.ActivityID != 1 {
		t.Errorf("Expected ActivityID 1, got %d", transaction.ActivityID)
	}
	if transaction.Points != 5 {
		t.Errorf("Expected Points %d, got %d", 5, transaction.Points)
	}

	// Test Milestone
	milestoneReached, _ = svc.AwardPoints(models.CategoryHabit, 1, MilestoneInterval)
	if !milestoneReached {
		t.Errorf("Expected milestone to be true after exceeding interval")
	}
}

func TestGetTotalPoints(t *testing.T) {
	db := setupTestDB(t)
	repo := pointRepo.NewSQLRepository(db)
	svc := NewService(repo)

	// Test empty
	total, err := svc.GetTotalPoints()
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if total != 0 {
		t.Errorf("Expected 0 points, got %d", total)
	}

	// Add some points
	_, _ = svc.AwardPoints(models.CategoryHabit, 1, 5)
	_, _ = svc.AwardPoints(models.CategoryFocus, 1, 10)

	total, err = svc.GetTotalPoints()
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if total != 15 {
		t.Errorf("Expected 15 points, got %d", total)
	}
}

func TestGetPointsInRange(t *testing.T) {
	db := setupTestDB(t)
	repo := pointRepo.NewSQLRepository(db)
	svc := NewService(repo)

	now := time.Now()

	// Transaction 1 (Yesterday)
	tx1 := models.PointTransaction{
		ActivityType: models.CategoryHabit,
		ActivityID:   1,
		Points:       5,
	}
	db.Create(&tx1)
	db.Model(&tx1).Update("CreatedAt", now.AddDate(0, 0, -1))

	// Transaction 2 (Today)
	tx2 := models.PointTransaction{
		ActivityType: models.CategoryFocus,
		ActivityID:   1,
		Points:       10,
	}
	db.Create(&tx2)

	// Transaction 3 (Tomorrow - out of range typically)
	tx3 := models.PointTransaction{
		ActivityType: models.CategoryIntent,
		ActivityID:   1,
		Points:       10,
	}
	db.Create(&tx3)
	db.Model(&tx3).Update("CreatedAt", now.AddDate(0, 0, 1))

	startStr := now.AddDate(0, 0, -2).Format("2006-01-02 15:04:05")
	endStr := now.Add(time.Hour).Format("2006-01-02 15:04:05")

	transactions, err := svc.GetPointsInRange(startStr, endStr)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions in range, got %d", len(transactions))
	}

	// Check order
	if transactions[0].Points != 5 || transactions[1].Points != 10 {
		t.Errorf("Expected ordered transactions")
	}
}
