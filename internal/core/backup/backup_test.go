package backup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newBackupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	if err := db.AutoMigrate(
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
	); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	return db
}

func TestExportImportIncludesTasksAndSubTasks(t *testing.T) {
	db := newBackupTestDB(t)
	service := NewService(db)

	task := models.Task{Title: "Back up tasks", Status: models.TaskStatusPending, Position: 7}
	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}
	subTask := models.SubTask{TaskID: task.ID, Title: "Back up subtasks", Status: models.TaskStatusPending, Position: 3}
	if err := db.Create(&subTask).Error; err != nil {
		t.Fatalf("create subtask: %v", err)
	}

	filePath := filepath.Join(t.TempDir(), "mindloop-backup.json")
	if err := service.Export(filePath); err != nil {
		t.Fatalf("export backup: %v", err)
	}

	backupBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	var exported Data
	if err := json.Unmarshal(backupBytes, &exported); err != nil {
		t.Fatalf("unmarshal backup: %v", err)
	}
	if len(exported.Tasks) != 1 || exported.Tasks[0].Title != task.Title {
		t.Fatalf("expected exported task, got %#v", exported.Tasks)
	}
	if len(exported.SubTasks) != 1 || exported.SubTasks[0].Title != subTask.Title {
		t.Fatalf("expected exported subtask, got %#v", exported.SubTasks)
	}

	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.SubTask{}).Error; err != nil {
		t.Fatalf("clear subtasks: %v", err)
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Task{}).Error; err != nil {
		t.Fatalf("clear tasks: %v", err)
	}

	if err := service.Import(filePath); err != nil {
		t.Fatalf("import backup: %v", err)
	}

	var restoredTask models.Task
	if err := db.Preload("SubTasks").First(&restoredTask, task.ID).Error; err != nil {
		t.Fatalf("load restored task: %v", err)
	}
	if restoredTask.Title != task.Title || restoredTask.Position != task.Position {
		t.Fatalf("restored task mismatch: %#v", restoredTask)
	}
	if len(restoredTask.SubTasks) != 1 {
		t.Fatalf("expected one restored subtask, got %#v", restoredTask.SubTasks)
	}
	if restoredTask.SubTasks[0].Title != subTask.Title || restoredTask.SubTasks[0].Position != subTask.Position {
		t.Fatalf("restored subtask mismatch: %#v", restoredTask.SubTasks[0])
	}
}
