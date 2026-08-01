package v1_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	v1 "github.com/snehmatic/mindloop/api/v1"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/backup"
	"github.com/snehmatic/mindloop/internal/core/dump"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/quest"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/core/task"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupTestServer(t *testing.T) *v1.MindloopHandler {
	// Initialize global config for tests
	config.InitConfig("MindloopTest", "local", ":8765")

	// Use in-memory DB for testing
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	// AutoMigrate manually to ensure tables exist
	err = database.AutoMigrate(
		&models.Intent{},
		&models.FocusSession{},
		&models.Habit{},
		&models.HabitLog{},
		&models.JournalEntry{},
		&models.Note{},
		&models.SideQuest{},
		&models.PointTransaction{},
		&models.Task{},
		&models.SubTask{},
		&models.AppSetting{},
		&models.BrainDump{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	journalService := journal.NewService(database)
	noteService := note.NewService(database)
	focusService := focus.NewService(database)
	intentService := intent.NewService(database)
	questService := quest.NewService(database)
	summaryService := summary.NewService(database)
	habitService := habit.NewService(database)
	backupService := backup.NewService(database)
	taskService := task.NewService(database)
	dumpService := dump.NewService(database)

	return v1.NewMindloopHandler(
		journalService,
		noteService,
		habitService,
		focusService,
		intentService,
		questService,
		summaryService,
		backupService,
		taskService,
		dumpService,
	)
}

func TestHabitFlow(t *testing.T) {
	mlh := setupTestServer(t)

	// 1. Create Habit
	val := url.Values{}
	val.Add("title", "Test Habit")
	val.Add("target_count", "1")
	val.Add("interval", "daily")

	req := httptest.NewRequest("POST", "/habits/new", strings.NewReader(val.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	mlh.HandleHabitCreate(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Create Habit failed, status: %d", resp.StatusCode)
	}

	// Verify Redirect to success
	loc, _ := resp.Location()
	if !strings.Contains(loc.String(), "success=true") {
		t.Errorf("Expected redirect to success=true, got %v", loc)
	}

	// Need to get the ID of the created habit for next steps.
	// We can list habits to find it.
	req = httptest.NewRequest("GET", "/habits", nil)
	w = httptest.NewRecorder()
	mlh.HandleHabitList(w, req)

	// We expect "Test Habit" in the body.
	// Also since ID autoincrements, it should be 1.
	habitID := "1"
	if !strings.Contains(w.Body.String(), "Test Habit") {
		t.Errorf("Habit list did not contain 'Test Habit'")
	}

	// 2. Log Habit
	val = url.Values{}
	val.Add("habit_id", habitID)
	req = httptest.NewRequest("POST", "/habits/log", strings.NewReader(val.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mlh.HandleHabitLog(w, req)

	resp = w.Result()
	loc, _ = resp.Location()
	if !strings.Contains(loc.String(), "success=done") && !strings.Contains(loc.String(), "success=true") {
		t.Errorf("Log Habit failed/redirected wrong: %v", loc)
	}

	// 3. Unlog Habit
	val = url.Values{}
	val.Add("habit_id", habitID)
	req = httptest.NewRequest("POST", "/habits/unlog", strings.NewReader(val.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mlh.HandleHabitUnlog(w, req)

	resp = w.Result()
	loc, _ = resp.Location()
	if !strings.Contains(loc.String(), "success=true") {
		t.Errorf("Unlog Habit failed/redirected wrong: %v", loc)
	}

	// 4. Delete Habit
	val = url.Values{}
	val.Add("habit_id", habitID)
	req = httptest.NewRequest("POST", "/habits/delete", strings.NewReader(val.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mlh.HandleHabitDelete(w, req)

	resp = w.Result()
	loc, _ = resp.Location()
	if !strings.Contains(loc.String(), "success=true") {
		t.Errorf("Delete Habit failed/redirected wrong: %v", loc)
	}
}

func TestSummaryGeneration(t *testing.T) {
	mlh := setupTestServer(t)

	req := httptest.NewRequest("GET", "/summary", nil)
	w := httptest.NewRecorder()
	mlh.HandleSummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Summary failed with status %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Summary Report") {
		t.Errorf("Summary page content missing expected title")
	}
}

func TestQuickDump(t *testing.T) {
	mlh := setupTestServer(t)

	t.Run("successful request", func(t *testing.T) {
		val := url.Values{}
		val.Add("content", "test dump")
		req := httptest.NewRequest("POST", "/api/dump", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		mlh.HandleQuickDump(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("empty content", func(t *testing.T) {
		val := url.Values{}
		val.Add("content", "")
		req := httptest.NewRequest("POST", "/api/dump", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		mlh.HandleQuickDump(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for empty content, got %d", w.Code)
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/dump", nil)
		w := httptest.NewRecorder()

		mlh.HandleQuickDump(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for GET method, got %d", w.Code)
		}
	})
}
