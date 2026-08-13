package task_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/task"
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

	err = db.AutoMigrate(&models.Task{}, &models.SubTask{}, &models.PointTransaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestTaskService(t *testing.T) {
	db := setupTestDB(t)
	s := task.NewService(db)

	// 1. Create Task
	var intentID uint = 1
	var focusID uint = 1
	createdTask, err := s.CreateTask("Test Task", &intentID, &focusID)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	if createdTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", createdTask.Title)
	}

	// 2. Add SubTask
	subTask, err := s.AddSubTask(createdTask.ID, "Test SubTask")
	if err != nil {
		t.Fatalf("Failed to add subtask: %v", err)
	}
	if subTask.Title != "Test SubTask" {
		t.Errorf("Expected title 'Test SubTask', got '%s'", subTask.Title)
	}

	// 3. List Tasks
	tasks, err := s.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}

	// 4. Get Task
	retrievedTask, err := s.GetTask(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if len(retrievedTask.SubTasks) != 1 {
		t.Errorf("Expected 1 subtask, got %d", len(retrievedTask.SubTasks))
	}

	// 5. Complete SubTask
	_, err = s.CompleteSubTask(subTask.ID, 5)
	if err != nil {
		t.Fatalf("Failed to complete subtask: %v", err)
	}

	// 6. Complete Task
	_, err = s.CompleteTask(createdTask.ID, 10)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	// 7. Delete Task
	err = s.DeleteTask(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	tasks, _ = s.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(tasks))
	}
}
