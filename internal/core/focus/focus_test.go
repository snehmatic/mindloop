package focus_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/focus"
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

	err = db.AutoMigrate(&models.FocusSession{}, &models.PointTransaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestStartSession(t *testing.T) {
	db := setupTestDB(t)
	s := focus.NewService(db)

	tests := []struct {
		name    string
		title   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "start first session",
			title:   "Session 1",
			wantErr: false,
		},
		{
			name:    "start second session while first is active",
			title:   "Session 2",
			wantErr: true,
			errMsg:  "a focus session is already active",
		},
		{
			name:    "start session with empty title",
			title:   "",
			wantErr: true,
			errMsg:  "title cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.StartSession(tt.title)
			if (err != nil) != tt.wantErr {
				t.Fatalf("StartSession() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("StartSession() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestEndAndRestartSession(t *testing.T) {
	db := setupTestDB(t)
	s := focus.NewService(db)

	// 1. Start session
	sess, err := s.StartSession("Session 1")
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}

	// 2. End session
	_, _, err = s.EndSession(int(sess.ID), 10)
	if err != nil {
		t.Fatalf("Failed to end session: %v", err)
	}

	// 3. Start another session - should succeed
	_, err = s.StartSession("Session 2")
	if err != nil {
		t.Errorf("Failed to start second session after ending first: %v", err)
	}
}
