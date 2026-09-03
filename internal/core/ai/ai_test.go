package ai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
func TestCustomProviderListModels(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(db)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			t.Errorf("Expected path /models, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", r.Header.Get("Authorization"))
		}

		response := map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "gpt-4-custom"},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	err := svc.SaveSettings("custom", "test-model", "test-token", mockServer.URL+"/")
	if err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	models, err := svc.ListModels()
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}

	if len(models) != 1 || models[0] != "gpt-4-custom" {
		t.Errorf("Expected models [gpt-4-custom], got %v", models)
	}
}

func TestCustomProviderTestConnection(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(db)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected path /chat/completions, got %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "Connection successful",
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	err := svc.SaveSettings("custom", "test-model", "test-token", mockServer.URL)
	if err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	err = svc.TestConnection("custom", "test-model", "test-token", mockServer.URL)
	if err != nil {
		t.Fatalf("TestConnection failed: %v", err)
	}
}

func TestCustomProviderGenerateJournal(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(db)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected path /chat/completions, got %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "Journal generated",
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	err := svc.SaveSettings("custom", "test-model", "test-token", mockServer.URL)
	if err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	summary := models.SummaryReport{}

	journal, err := svc.GenerateJournal(summary)
	if err != nil {
		t.Fatalf("GenerateJournal failed: %v", err)
	}

	if journal != "Journal generated" {
		t.Errorf("Expected 'Journal generated', got '%s'", journal)
	}
}

func TestCustomProviderListModels_NonOpenAIModels(t *testing.T) {
	db := setupTestDB(t)
	svc := ai.NewService(db)

	// Simulate an Ollama-style response with non-OpenAI model names
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "llama3.1:8b"},
				{"id": "qwen3-coder:latest"},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	err := svc.SaveSettings("custom", "llama3.1:8b", "ollama", mockServer.URL)
	if err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	models, err := svc.ListModels()
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}

	if len(models) != 2 {
		t.Errorf("Expected 2 models, got %d: %v", len(models), models)
	}
}

func TestGenerateChunker_NoToken(t *testing.T) {
	db := setupTestDB(t)

	svc := ai.NewService(db)

	_, err := svc.GenerateChunker("do my homework")
	if err == nil {
		t.Error("Expected error for unconfigured token, got nil")
	}
}
