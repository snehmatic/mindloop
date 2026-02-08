package note_test

import (
	"testing"

	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/models"
	"github.com/glebarez/sqlite"
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

	err = db.AutoMigrate(&models.Note{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestNoteService(t *testing.T) {
	db := setupTestDB(t)
	s := note.NewService(db)

	// 1. Create Note
	n, err := s.CreateNote("Test Title", "Test Content", "label1,label2")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}
	if n.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", n.Title)
	}

	// 2. List Notes
	notes, err := s.ListNotes()
	if err != nil {
		t.Fatalf("Failed to list notes: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(notes))
	}

	// 3. Get Note
	n2, err := s.GetNote(int(n.ID))
	if err != nil {
		t.Fatalf("Failed to get note: %v", err)
	}
	if n2.Content != "Test Content" {
		t.Errorf("Expected content 'Test Content', got '%s'", n2.Content)
	}

	// 4. Update Note
	n3, err := s.UpdateNote(int(n.ID), "Updated Title", "Updated Content", "newlabel")
	if err != nil {
		t.Fatalf("Failed to update note: %v", err)
	}
	if n3.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", n3.Title)
	}

	// 5. Delete Note
	err = s.DeleteNote(int(n.ID))
	if err != nil {
		t.Fatalf("Failed to delete note: %v", err)
	}

	notes, _ = s.ListNotes()
	if len(notes) != 0 {
		t.Errorf("Expected 0 notes after deletion, got %d", len(notes))
	}
}

func TestCreateNoteEmpty(t *testing.T) {
	db := setupTestDB(t)
	s := note.NewService(db)

	_, err := s.CreateNote("", "", "")
	if err == nil {
		t.Error("Expected error when creating empty note, got nil")
	}
}
