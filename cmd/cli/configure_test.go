package cli

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/core/ai"
	"github.com/snehmatic/mindloop/internal/repository/appsettings"
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

func TestSetupAIConfig(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(appsettings.NewSQLRepository(db))

	input := "openai\ngpt-4\ntest-token\n"
	reader := bufio.NewReader(strings.NewReader(input))
	var writer bytes.Buffer

	err := SetupAIConfig(reader, &writer, svc)
	if err != nil {
		t.Fatalf("SetupAIConfig failed: %v", err)
	}

	provider, model, token, baseURL, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings failed: %v", err)
	}

	if provider != "openai" {
		t.Errorf("Expected openai, got %s", provider)
	}
	if model != "gpt-4" {
		t.Errorf("Expected gpt-4, got %s", model)
	}
	if token != "test-token" {
		t.Errorf("Expected test-token, got %s", token)
	}
	if baseURL != "" {
		t.Errorf("Expected empty baseURL, got %s", baseURL)
	}
}

func TestSetupAIConfigCustom(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(appsettings.NewSQLRepository(db))

	input := "custom\nhttp://localhost:11434/v1\nllama3\nnone\n"
	reader := bufio.NewReader(strings.NewReader(input))
	var writer bytes.Buffer

	err := SetupAIConfig(reader, &writer, svc)
	if err != nil {
		t.Fatalf("SetupAIConfig failed: %v", err)
	}

	provider, model, token, baseURL, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings failed: %v", err)
	}

	if provider != "custom" {
		t.Errorf("Expected custom, got %s", provider)
	}
	if baseURL != "http://localhost:11434/v1" {
		t.Errorf("Expected http://localhost:11434/v1, got %s", baseURL)
	}
	if model != "llama3" {
		t.Errorf("Expected llama3, got %s", model)
	}
	if token != "" {
		t.Errorf("Expected empty token, got %s", token)
	}
}

func TestSetupAIConfigInvalidProvider(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(appsettings.NewSQLRepository(db))

	// first input is invalid, second is valid
	input := "invalid\ngemini\ngemini-1.5-pro\ntoken\n"
	reader := bufio.NewReader(strings.NewReader(input))
	var writer bytes.Buffer

	err := SetupAIConfig(reader, &writer, svc)
	if err != nil {
		t.Fatalf("SetupAIConfig failed: %v", err)
	}

	if !strings.Contains(writer.String(), "Invalid provider") {
		t.Errorf("Expected warning message about invalid provider")
	}

	provider, _, _, _, _ := svc.GetSettings()
	if provider != "gemini" {
		t.Errorf("Expected gemini, got %s", provider)
	}
}

func TestSetupAIConfigInvalidBaseURL(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(appsettings.NewSQLRepository(db))

	// invalid url first, then valid
	input := "custom\nlocalhost:11434\nhttp://localhost:11434/v1\nllama3\nnone\n"
	reader := bufio.NewReader(strings.NewReader(input))
	var writer bytes.Buffer

	err := SetupAIConfig(reader, &writer, svc)
	if err != nil {
		t.Fatalf("SetupAIConfig failed: %v", err)
	}

	if !strings.Contains(writer.String(), "Invalid Base URL") {
		t.Errorf("Expected warning message about invalid Base URL")
	}

	_, _, _, baseURL, _ := svc.GetSettings()
	if baseURL != "http://localhost:11434/v1" {
		t.Errorf("Expected http://localhost:11434/v1, got %s", baseURL)
	}
}
