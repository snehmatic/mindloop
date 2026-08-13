package dump_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/dump"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, *dump.Service) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.BrainDump{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	service := dump.NewService(db)
	return db, service
}

func TestCreateDump(t *testing.T) {
	_, service := setupTestDB(t)

	// Test successful creation
	t.Run("success", func(t *testing.T) {
		bd, err := service.CreateDump("Buy milk")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if bd.Content != "Buy milk" {
			t.Errorf("expected content to be 'Buy milk', got %s", bd.Content)
		}
		if bd.ID == 0 {
			t.Error("expected ID to be set")
		}
	})

	// Test empty content
	t.Run("empty content", func(t *testing.T) {
		_, err := service.CreateDump("")
		if err == nil {
			t.Error("expected error for empty content, got nil")
		}
	})
}
