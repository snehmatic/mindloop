package task

import (
	"testing"

	"github.com/glebarez/sqlite"
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

	err = db.AutoMigrate(
		&models.Task{},
		&models.SubTask{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	return db
}

func TestTaskRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLTaskRepository(db)

	// Test CreateTask
	task, err := repo.CreateTask("Test Task", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", task.Title)
	}

	// Test GetTask
	retrieved, err := repo.GetTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if retrieved.ID != task.ID {
		t.Errorf("Expected ID %d, got %d", task.ID, retrieved.ID)
	}

	// Test ListTasks
	tasks, err := repo.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// Test AddSubTask
	subTask, err := repo.AddSubTask(task.ID, "Test SubTask")
	if err != nil {
		t.Fatalf("Failed to add subtask: %v", err)
	}
	if subTask.Title != "Test SubTask" {
		t.Errorf("Expected title 'Test SubTask', got '%s'", subTask.Title)
	}

	// Test GetSubTask
	retrievedSt, err := repo.GetSubTask(subTask.ID)
	if err != nil {
		t.Fatalf("Failed to get subtask: %v", err)
	}
	if retrievedSt.ID != subTask.ID {
		t.Errorf("Expected ID %d, got %d", subTask.ID, retrievedSt.ID)
	}

	// Test CompleteSubTask
	completedSt, err := repo.CompleteSubTask(subTask.ID)
	if err != nil {
		t.Fatalf("Failed to complete subtask: %v", err)
	}
	if completedSt.Status != models.TaskStatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", completedSt.Status)
	}

	// Test CompleteTask
	completed, err := repo.CompleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}
	if completed.Status != models.TaskStatusCompleted {
		t.Errorf("Expected status 'completed', got '%s'", completed.Status)
	}

	// Test DeleteSubTask
	err = repo.DeleteSubTask(subTask.ID)
	if err != nil {
		t.Fatalf("Failed to delete subtask: %v", err)
	}

	// Test DeleteTask
	err = repo.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
}

func TestTaskRepositoryReorder(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLTaskRepository(db)

	// Create multiple tasks
	task1, _ := repo.CreateTask("Task 1", nil, nil)
	task2, _ := repo.CreateTask("Task 2", nil, nil)
	task3, _ := repo.CreateTask("Task 3", nil, nil)

	// Reorder
	err := repo.ReorderTasks([]uint{
		task3.ID,
		task1.ID,
		task2.ID,
	})
	if err != nil {
		t.Fatalf("Failed to reorder tasks: %v", err)
	}

	// Verify order
	tasks, _ := repo.ListTasks()
	if len(tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != task3.ID || tasks[1].ID != task1.ID || tasks[2].ID != task2.ID {
		t.Errorf("Tasks not in expected order")
	}
}
