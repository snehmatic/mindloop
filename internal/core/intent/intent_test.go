package intent_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/intent"
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

	err = db.AutoMigrate(&models.Intent{}, &models.PointTransaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestIntentService(t *testing.T) {
	db := setupTestDB(t)
	s := intent.NewService(db)

	// 1. Start Intent
	i, err := s.StartIntent("Test Intent")
	if err != nil {
		t.Fatalf("Failed to start intent: %v", err)
	}
	if i.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", i.Status)
	}

	// 2. List Active Intents
	active, err := s.ListActiveIntents()
	if err != nil {
		t.Fatalf("Failed to list active intents: %v", err)
	}
	if len(active) != 1 {
		t.Errorf("Expected 1 active intent, got %d", len(active))
	}

	// 3. End Intent
	id := "1"
	ended, _, err := s.EndIntent(id, 10)
	if err != nil {
		t.Fatalf("Failed to end intent: %v", err)
	}
	if ended.Status != "done" {
		t.Errorf("Expected status 'done', got '%s'", ended.Status)
	}
	if ended.EndedAt == nil {
		t.Error("Expected EndedAt to be set")
	}

	// 4. List Active Intents (should be 0)
	active, _ = s.ListActiveIntents()
	if len(active) != 0 {
		t.Errorf("Expected 0 active intents, got %d", len(active))
	}

	// 5. List All Intents (should be 1)
	all, _ := s.ListIntents()
	if len(all) != 1 {
		t.Errorf("Expected 1 total intent, got %d", len(all))
	}

	// 6. Delete Intent
	err = s.DeleteIntent(id)
	if err != nil {
		t.Fatalf("Failed to delete intent: %v", err)
	}
	all, _ = s.ListIntents()
	if len(all) != 0 {
		t.Errorf("Expected 0 intents after deletion, got %d", len(all))
	}
}
