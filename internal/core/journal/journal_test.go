package journal_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/journal"
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

	err = db.AutoMigrate(&models.JournalEntry{}, &models.PointTransaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestJournalService(t *testing.T) {
	db := setupTestDB(t)
	s := journal.NewService(db)

	// 1. Create Entry
	_, err := s.CreateEntry("Test Journal", "Test Content", "happy", 5)
	if err != nil {
		t.Fatalf("Failed to create journal entry: %v", err)
	}

	// 2. List Entries
	entries, err := s.ListEntries()
	if err != nil {
		t.Fatalf("Failed to list entries: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Title != "Test Journal" {
		t.Errorf("Expected title 'Test Journal', got '%s'", entries[0].Title)
	}

	// 3. Get Entry
	id := "1" // SQLite in-memory starts with 1
	entry, err := s.GetEntry(id)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}
	if entry.Content != "Test Content" {
		t.Errorf("Expected content 'Test Content', got '%s'", entry.Content)
	}

	// 4. Delete Entry
	err = s.DeleteEntry(id)
	if err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	entries, _ = s.ListEntries()
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries after deletion, got %d", len(entries))
	}
}

func TestCreateEntryValidation(t *testing.T) {
	db := setupTestDB(t)
	s := journal.NewService(db)

	tests := []struct {
		name    string
		title   string
		content string
		wantErr bool
	}{
		{"valid", "Title", "Content", false},
		{"empty title", "", "Content", true},
		{"empty content", "Title", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.CreateEntry(tt.title, tt.content, "neutral", 5)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
