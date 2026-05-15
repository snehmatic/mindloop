package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAISettingsAndConnection(t *testing.T) {
	mlh := setupTestServer(t)

	// Set up a mock AI provider server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/chat/completions" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"choices": [{"message": {"content": "mocked response"}}]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// 1. Save AI Settings (including mockServer URL as BaseURL)
	savePayload := `{"provider": "custom", "model": "gpt-mock", "token": "test-token", "baseURL": "` + mockServer.URL + `"}`
	req := httptest.NewRequest("POST", "/settings/ai", strings.NewReader(savePayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mlh.HandleSaveAISettings(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to save AI settings, got status %d", w.Code)
	}

	// 2. Get AI Settings and verify
	req = httptest.NewRequest("GET", "/settings/ai", nil)
	w = httptest.NewRecorder()
	mlh.HandleGetAISettings(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get AI settings, got status %d", w.Code)
	}

	var getResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("Failed to parse get settings response: %v", err)
	}

	if getResp["provider"] != "custom" || getResp["model"] != "gpt-mock" || getResp["baseURL"] != mockServer.URL {
		t.Errorf("Get settings returned unexpected values: %v", getResp)
	}

	// 3. Test AI Connection with empty token and baseURL (fallback behavior)
	testPayload := `{"provider": "custom", "model": "gpt-mock", "token": "", "baseURL": ""}`
	req = httptest.NewRequest("POST", "/settings/ai/test", strings.NewReader(testPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mlh.HandleTestAIConnection(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Failed to test AI connection with empty fallback, got status %d, body: %s", w.Code, w.Body.String())
	}

	// 4. Test AI Connection with only token provided (baseURL fallback)
	testPayload = `{"provider": "custom", "model": "gpt-mock", "token": "test-token", "baseURL": ""}`
	req = httptest.NewRequest("POST", "/settings/ai/test", strings.NewReader(testPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mlh.HandleTestAIConnection(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Failed to test AI connection with baseURL fallback, got status %d, body: %s", w.Code, w.Body.String())
	}

	// 5. Test AI Connection with only baseURL provided (token fallback)
	testPayload = `{"provider": "custom", "model": "gpt-mock", "token": "", "baseURL": "` + mockServer.URL + `"}`
	req = httptest.NewRequest("POST", "/settings/ai/test", strings.NewReader(testPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mlh.HandleTestAIConnection(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Failed to test AI connection with token fallback, got status %d, body: %s", w.Code, w.Body.String())
	}
}
