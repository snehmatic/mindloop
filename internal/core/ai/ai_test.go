package ai_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/ai"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupTestDB(t *testing.T) *gorm.DB {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	err = database.AutoMigrate(&models.AppSetting{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}
	return database
}

func TestSaveAndGetSettingsWithBaseURL(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(db)

	err := svc.SaveSettings("custom", "test-model", "test-token", "http://localhost:11434/v1")
	if err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	provider, model, token, baseURL, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings failed: %v", err)
	}

	if provider != "custom" {
		t.Errorf("Expected custom, got %v", provider)
	}
	if baseURL != "http://localhost:11434/v1" {
		t.Errorf("Expected http://localhost:11434/v1, got %v", baseURL)
	}
	if model != "test-model" {
		t.Errorf("Expected test-model, got %v", model)
	}
	if token != "test-token" {
		t.Errorf("Expected test-token, got %v", token)
	}
}
