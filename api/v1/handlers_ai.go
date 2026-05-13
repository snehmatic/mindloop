package v1

import (
	"encoding/json"
	"net/http"

	"github.com/snehmatic/mindloop/internal/core/ai"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/utils"
)

type AISettingsRequest struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Token    string `json:"token"`
}

func (mlh *MindloopHandler) HandleGetAISettings(w http.ResponseWriter, r *http.Request) {
	// Re-initialize AI service with current DB
	aiService := ai.NewService(mlh.journal.DB) // hack: accessing DB via a service that has it
	provider, model, token, _ := aiService.GetSettings()

	hasToken := token != ""

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"provider": provider,
		"model":    model,
		"hasToken": hasToken,
	})
}

func (mlh *MindloopHandler) HandleSaveAISettings(w http.ResponseWriter, r *http.Request) {
	var req AISettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	aiService := ai.NewService(mlh.journal.DB)
	if err := aiService.SaveSettings(req.Provider, req.Model, req.Token); err != nil {
		http.Error(w, "Failed to save settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Settings saved successfully"})
}

func (mlh *MindloopHandler) HandleListAIModels(w http.ResponseWriter, r *http.Request) {
	aiService := ai.NewService(mlh.journal.DB)
	models, err := aiService.ListModels()
	if err != nil {
		http.Error(w, "Failed to list models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"models": models,
	})
}

func (mlh *MindloopHandler) HandleTestAIConnection(w http.ResponseWriter, r *http.Request) {
	var req AISettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	aiService := ai.NewService(mlh.journal.DB)

	// If token is empty in the request, fallback to the saved token
	if req.Token == "" {
		_, _, savedToken, _ := aiService.GetSettings()
		req.Token = savedToken
	}

	if err := aiService.TestConnection(req.Provider, req.Model, req.Token); err != nil {
		http.Error(w, "Connection test failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Connection successful!"})
}

func (mlh *MindloopHandler) HandleGenerateAIJournal(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	start, end := utils.GetDateRange(period)
	summaryService := summary.NewService(mlh.journal.DB)
	report, err := summaryService.GenerateSummary(start, end)
	if err != nil {
		http.Error(w, "Failed to generate summary data", http.StatusInternalServerError)
		return
	}

	aiService := ai.NewService(mlh.journal.DB)
	generatedText, err := aiService.GenerateJournal(report)
	if err != nil {
		http.Error(w, "AI Generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"content": generatedText,
		"title":   "AI Summary: " + report.DateRange,
	})
}
