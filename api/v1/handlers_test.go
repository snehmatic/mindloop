//go:build test

//^ test tag is required here to resolve ResetForTest()

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
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	coreJournal "github.com/snehmatic/mindloop/internal/core/journal"
	coreNote "github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/points"
	coreQuest "github.com/snehmatic/mindloop/internal/core/quest"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/core/task"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/repository/appsettings"
	fRepo "github.com/snehmatic/mindloop/internal/repository/focus"
	hRepo "github.com/snehmatic/mindloop/internal/repository/habit"
	iRepo "github.com/snehmatic/mindloop/internal/repository/intent"
	journalRepo "github.com/snehmatic/mindloop/internal/repository/journal"
	noteRepo "github.com/snehmatic/mindloop/internal/repository/note"
	pRepo "github.com/snehmatic/mindloop/internal/repository/point"
	questRepo "github.com/snehmatic/mindloop/internal/repository/quest"
	"github.com/snehmatic/mindloop/internal/repository/routine"
	"github.com/snehmatic/mindloop/internal/repository/subtask"
	taskRepo "github.com/snehmatic/mindloop/internal/repository/task"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupTestServer(t *testing.T) *v1.MindloopHandler {
	// Reset config for test isolation
	config.ResetForTest()
	// Initialize global config for tests
	if err := config.Init("MindloopTest", "local", ":8765"); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

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
		&models.Routine{},
		&models.Task{},
		&models.SubTask{},
		&models.AppSetting{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	jRepo := journalRepo.NewSQLRepository(database)
	journalService := coreJournal.NewService(jRepo)
	nRepo := noteRepo.NewSQLRepository(database)
	noteService := coreNote.NewService(nRepo)
	appSettingsRepo := appsettings.NewSQLRepository(database)
	fRepo := fRepo.NewSQLRepository(database)
	focusService := focus.NewService(fRepo)
	iRepo := iRepo.NewSQLRepository(database)
	intentService := intent.NewService(iRepo)
	qRepo := questRepo.NewSQLRepository(database, log.Get())
	questService := coreQuest.NewService(qRepo, log.Get())
	hRepo := hRepo.NewSQLRepository(database, nil)
	habitService := habit.NewService(hRepo)
	summaryService := summary.NewService(fRepo, hRepo, iRepo, pRepo.NewSQLRepository(database), taskRepo.NewSQLTaskRepository(database), log.Get())
	backupService := backup.NewService(database, points.NewService(pRepo.NewSQLRepository(database)), fRepo, hRepo, nil, iRepo, jRepo, nRepo, pRepo.NewSQLRepository(database), qRepo, routine.NewSQLRepository(database), subtask.NewSQLRepository(database), taskRepo.NewSQLTaskRepository(database))
	pointSvc := points.NewService(pRepo.NewSQLRepository(database))
	taskService := task.NewService(taskRepo.NewSQLTaskRepository(database), config.GetUserConfig(), pointSvc, log.Get())

	return v1.NewMindloopHandler(
		database,
		appSettingsRepo,
		journalService,
		noteService,
		habitService,
		focusService,
		intentService,
		questService,
		summaryService,
		backupService,
		taskService,
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
