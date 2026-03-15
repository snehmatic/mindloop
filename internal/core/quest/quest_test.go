package quest_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/quest"
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

	err = db.AutoMigrate(&models.SideQuest{}, &models.PointTransaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestQuestService(t *testing.T) {
	db := setupTestDB(t)
	s := quest.NewService(db)

	// 1. Start Quest
	q, err := s.StartQuest("Emergency Fix")
	if err != nil {
		t.Fatalf("Failed to start quest: %v", err)
	}
	if q.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", q.Status)
	}

	// 2. Start Another Quest (Should fail)
	_, err = s.StartQuest("Another One")
	if err == nil {
		t.Error("Expected error when starting second active quest, got nil")
	}

	// 3. Get Active Quest
	active, err := s.GetActiveQuest()
	if err != nil {
		t.Fatalf("Failed to get active quest: %v", err)
	}
	if active.ID != q.ID {
		t.Errorf("Expected active quest ID %d, got %d", q.ID, active.ID)
	}

	// 4. Stop Quest
	completed, _, err := s.StopQuest(q.ID, "Fixed it", 5)
	if err != nil {
		t.Fatalf("Failed to stop quest: %v", err)
	}
	if completed.Status != "done" {
		t.Errorf("Expected status 'done', got '%s'", completed.Status)
	}
	if completed.Note != "Fixed it" {
		t.Errorf("Expected note 'Fixed it', got '%s'", completed.Note)
	}

	// 5. Get Active Quest (Should return nil, nil)
	active, err = s.GetActiveQuest()
	if err != nil {
		t.Fatalf("Expected no error when getting non-existent active quest, got %v", err)
	}
	if active != nil {
		t.Error("Expected nil quest when none is active")
	}

	// 6. List Quests
	list, err := s.ListQuests()
	if err != nil {
		t.Fatalf("Failed to list quests: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 quest in list, got %d", len(list))
	}
}
