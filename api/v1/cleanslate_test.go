package v1

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/backup"
	"github.com/snehmatic/mindloop/internal/core/focus"
	coreHabit "github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	coreJournal "github.com/snehmatic/mindloop/internal/core/journal"
	coreNote "github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/points"
	coreQuest "github.com/snehmatic/mindloop/internal/core/quest"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/core/task"
	"github.com/snehmatic/mindloop/internal/log"
	focusRepo "github.com/snehmatic/mindloop/internal/repository/focus"
	habitRepo "github.com/snehmatic/mindloop/internal/repository/habit"
	"github.com/snehmatic/mindloop/internal/repository/habitlog"
	intentRepo "github.com/snehmatic/mindloop/internal/repository/intent"
	journalRepo "github.com/snehmatic/mindloop/internal/repository/journal"
	noteRepo "github.com/snehmatic/mindloop/internal/repository/note"
	pointRepo "github.com/snehmatic/mindloop/internal/repository/point"
	questRepoPkg "github.com/snehmatic/mindloop/internal/repository/quest"
	routineRepo "github.com/snehmatic/mindloop/internal/repository/routine"
	subtaskRepo "github.com/snehmatic/mindloop/internal/repository/subtask"
	taskRepo "github.com/snehmatic/mindloop/internal/repository/task"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

func setupCleanSlateTest(t *testing.T) (*MindloopHandler, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to memory db: %v", err)
	}

	// AutoMigrate all models
	err = db.AutoMigrate(
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
		t.Fatalf("Failed to migrate: %v", err)
	}

	logger := log.Get()
	fRepo := focusRepo.NewSQLRepository(db)
	hRepo := habitRepo.NewSQLRepository(db, logger)
	hlRepo := habitlog.NewSQLRepository(db)
	iRepo := intentRepo.NewSQLRepository(db)
	jRepo := journalRepo.NewSQLRepository(db)
	nRepo := noteRepo.NewSQLRepository(db)
	pRepo := pointRepo.NewSQLRepository(db)
	qRepo := questRepoPkg.NewSQLRepository(db, logger)
	rRepo := routineRepo.NewSQLRepository(db)
	stRepo := subtaskRepo.NewSQLRepository(db)
	tRepo := taskRepo.NewSQLTaskRepository(db)
	pointSvc := points.NewService(pRepo)
	tSvcRepo := taskRepo.NewSQLTaskRepository(db)
	tService := task.NewService(tSvcRepo, config.GetUserConfig(), pointSvc, logger)
	bService := backup.NewService(db, pointSvc, fRepo, hRepo, hlRepo, iRepo, jRepo, nRepo, pRepo, qRepo, rRepo, stRepo, tRepo)
	hService := coreHabit.NewService(hRepo)
	jService := coreJournal.NewService(jRepo)
	nService := coreNote.NewService(nRepo)
	fService := focus.NewService(fRepo)
	iService := intent.NewService(iRepo)
	qService := coreQuest.NewService(qRepo, logger)
	sService := summary.NewService(fRepo, hRepo, iRepo, pRepo, tRepo, logger)

	mlh := NewMindloopHandler(db, jService, nService, hService, fService, iService, qService, sService, bService, tService)
	return mlh, db
}

func TestCleanSlate(t *testing.T) {
	mlh, db := setupCleanSlateTest(t)

	// 1. Seed Data
	db.Create(&models.JournalEntry{Title: "ToDelete", Content: "Content", Mood: "happy"})
	db.Create(&models.Habit{Title: "ToDelete", TargetCount: 1, Interval: "daily"})

	// Verify Seeding
	var jCount, hCount int64
	db.Model(&models.JournalEntry{}).Count(&jCount)
	db.Model(&models.Habit{}).Count(&hCount)

	if jCount != 1 || hCount != 1 {
		t.Fatalf("Setup failed: expected 1 journal and 1 habit, got %d and %d", jCount, hCount)
	}

	// 2. Perform Clean Slate Request (Type=all)
	data := url.Values{}
	data.Set("type", "all")
	req, _ := http.NewRequest("POST", "/cleanslate", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	mlh.HandleCleanSlate(rr, req)

	// 3. Verify Result
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect 303, got %d", rr.Code)
	}

	// Check DB counts
	db.Model(&models.JournalEntry{}).Count(&jCount)
	db.Model(&models.Habit{}).Count(&hCount)

	if jCount != 0 {
		t.Errorf("Clean Slate failed: expected 0 journals, got %d", jCount)
	}
	if hCount != 0 {
		t.Errorf("Clean Slate failed: expected 0 habits, got %d", hCount)
	}
}

func TestCleanSlateJournalOnly(t *testing.T) {
	mlh, db := setupCleanSlateTest(t)

	// 1. Seed Data
	db.Create(&models.JournalEntry{Title: "ToDelete", Content: "Content"})
	db.Create(&models.Habit{Title: "KeepMe", TargetCount: 1})

	// 2. Perform Clean Slate (Type=journal)
	data := url.Values{}
	data.Set("type", "journal")
	req, _ := http.NewRequest("POST", "/cleanslate", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	mlh.HandleCleanSlate(rr, req)

	// 3. Verify
	var jCount, hCount int64
	db.Model(&models.JournalEntry{}).Count(&jCount)
	db.Model(&models.Habit{}).Count(&hCount)

	if jCount != 0 {
		t.Errorf("Expected 0 journals, got %d", jCount)
	}
	if hCount != 1 {
		t.Errorf("Expected 1 habit remaining, got %d", hCount)
	}
}
