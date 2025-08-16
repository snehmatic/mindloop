package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/snehmatic/mindloop/internal/application"
	"github.com/snehmatic/mindloop/internal/domain/entities"
)

type MindloopHandler struct {
	container *application.Container
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HabitRequest represents a habit creation/update request
type HabitRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetCount int    `json:"target_count"`
	Interval    string `json:"interval"` // "daily" or "weekly"
}

// IntentRequest represents an intent creation request
type IntentRequest struct {
	Name string `json:"name"`
}

// FocusRequest represents a focus session creation request
type FocusRequest struct {
	Title string `json:"title"`
}

// FocusRateRequest represents a focus session rating request
type FocusRateRequest struct {
	Rating int `json:"rating"` // 0-10
}

// JournalRequest represents a journal entry creation/update request
type JournalRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Mood    string `json:"mood"` // "happy", "sad", "neutral", "angry", "excited"
}

// SummaryRequest represents a summary generation request
type SummaryRequest struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD format
	EndDate   string `json:"end_date"`   // YYYY-MM-DD format
}

func NewMindloopHandler(container *application.Container) *MindloopHandler {
	return &MindloopHandler{
		container: container,
	}
}

// writeJSON writes a JSON response
func (mlh *MindloopHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func (mlh *MindloopHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	response := Response{
		Success: false,
		Error:   message,
	}
	mlh.writeJSON(w, statusCode, response)
}

// writeSuccess writes a success response
func (mlh *MindloopHandler) writeSuccess(w http.ResponseWriter, data interface{}, message string) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	mlh.writeJSON(w, http.StatusOK, response)
}

// HandleHome handles the home endpoint
func (mlh *MindloopHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: "Welcome to Mindloop API!",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"features": []string{
				"habits", "intents", "focus", "journal", "summary",
			},
		},
	}
	mlh.writeJSON(w, http.StatusOK, response)
}

// HandleHealthz handles the health check endpoint
func (mlh *MindloopHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	// Check database health
	if err := mlh.container.DB.Health(); err != nil {
		mlh.writeError(w, http.StatusServiceUnavailable, "Database connection issue")
		return
	}
	mlh.writeSuccess(w, map[string]string{"status": "healthy"}, "Service is healthy")
}

// ===== HABIT HANDLERS =====

// HandleCreateHabit handles habit creation
func (mlh *MindloopHandler) HandleCreateHabit(w http.ResponseWriter, r *http.Request) {
	var req HabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Title == "" {
		mlh.writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.TargetCount <= 0 {
		mlh.writeError(w, http.StatusBadRequest, "Target count must be greater than 0")
		return
	}

	// Parse interval
	var interval entities.IntervalType
	switch req.Interval {
	case "daily":
		interval = entities.Daily
	case "weekly":
		interval = entities.Weekly
	default:
		interval = entities.Daily // default to daily
	}

	habit, err := mlh.container.HabitUseCase.CreateHabit(req.Title, req.Description, req.TargetCount, interval)
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create habit: %v", err))
		return
	}

	mlh.writeSuccess(w, habit, "Habit created successfully")
}

// HandleListHabits handles listing all habits
func (mlh *MindloopHandler) HandleListHabits(w http.ResponseWriter, r *http.Request) {
	habits, err := mlh.container.HabitUseCase.GetAllHabits()
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch habits: %v", err))
		return
	}

	mlh.writeSuccess(w, habits, "Habits retrieved successfully")
}

// HandleGetHabit handles getting a specific habit
func (mlh *MindloopHandler) HandleGetHabit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid habit ID")
		return
	}

	habit, err := mlh.container.HabitUseCase.GetHabit(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusNotFound, fmt.Sprintf("Habit not found: %v", err))
		return
	}

	mlh.writeSuccess(w, habit, "Habit retrieved successfully")
}

// HandleDeleteHabit handles habit deletion
func (mlh *MindloopHandler) HandleDeleteHabit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid habit ID")
		return
	}

	err = mlh.container.HabitUseCase.DeleteHabit(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete habit: %v", err))
		return
	}

	mlh.writeSuccess(w, nil, "Habit deleted successfully")
}

// HabitLogRequest represents a habit logging request
type HabitLogRequest struct {
	ActualCount int `json:"actual_count"`
}

// HandleLogHabit handles habit logging
func (mlh *MindloopHandler) HandleLogHabit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid habit ID")
		return
	}

	var req HabitLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ActualCount < 0 {
		mlh.writeError(w, http.StatusBadRequest, "Actual count must be non-negative")
		return
	}

	err = mlh.container.HabitUseCase.LogHabit(uint(id), req.ActualCount)
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to log habit: %v", err))
		return
	}

	mlh.writeSuccess(w, map[string]interface{}{
		"habit_id":     id,
		"actual_count": req.ActualCount,
		"logged_at":    time.Now(),
	}, "Habit logged successfully")
}

// ===== INTENT HANDLERS =====

// HandleCreateIntent handles intent creation
func (mlh *MindloopHandler) HandleCreateIntent(w http.ResponseWriter, r *http.Request) {
	var req IntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		mlh.writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	intent, err := mlh.container.IntentUseCase.StartIntent(req.Name)
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create intent: %v", err))
		return
	}

	mlh.writeSuccess(w, intent, "Intent created successfully")
}

// HandleListIntents handles listing all intents
func (mlh *MindloopHandler) HandleListIntents(w http.ResponseWriter, r *http.Request) {
	active := r.URL.Query().Get("active")

	var intents []*entities.Intent
	var err error

	if active == "true" {
		intents, err = mlh.container.IntentUseCase.GetActiveIntents()
	} else {
		intents, err = mlh.container.IntentUseCase.GetAllIntents()
	}

	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch intents: %v", err))
		return
	}

	mlh.writeSuccess(w, intents, "Intents retrieved successfully")
}

// HandleEndIntent handles ending an intent
func (mlh *MindloopHandler) HandleEndIntent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid intent ID")
		return
	}

	intent, err := mlh.container.IntentUseCase.EndIntent(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to end intent: %v", err))
		return
	}

	mlh.writeSuccess(w, intent, "Intent ended successfully")
}

// HandleDeleteIntent handles intent deletion
func (mlh *MindloopHandler) HandleDeleteIntent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid intent ID")
		return
	}

	err = mlh.container.IntentUseCase.DeleteIntent(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete intent: %v", err))
		return
	}

	mlh.writeSuccess(w, nil, "Intent deleted successfully")
}

// ===== FOCUS HANDLERS =====

// HandleCreateFocus handles focus session creation
func (mlh *MindloopHandler) HandleCreateFocus(w http.ResponseWriter, r *http.Request) {
	var req FocusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		mlh.writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	session, err := mlh.container.FocusUseCase.StartFocusSession(req.Title)
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create focus session: %v", err))
		return
	}

	mlh.writeSuccess(w, session, "Focus session created successfully")
}

// HandleListFocus handles listing focus sessions
func (mlh *MindloopHandler) HandleListFocus(w http.ResponseWriter, r *http.Request) {
	active := r.URL.Query().Get("active")

	var sessions []*entities.FocusSession
	var err error

	if active == "true" {
		sessions, err = mlh.container.FocusUseCase.GetActiveFocusSessions()
	} else {
		sessions, err = mlh.container.FocusUseCase.GetAllFocusSessions()
	}

	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch focus sessions: %v", err))
		return
	}

	mlh.writeSuccess(w, sessions, "Focus sessions retrieved successfully")
}

// HandleEndFocus handles ending a focus session
func (mlh *MindloopHandler) HandleEndFocus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid focus session ID")
		return
	}

	session, err := mlh.container.FocusUseCase.EndFocusSession(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to end focus session: %v", err))
		return
	}

	mlh.writeSuccess(w, session, "Focus session ended successfully")
}

// HandlePauseFocus handles pausing a focus session
func (mlh *MindloopHandler) HandlePauseFocus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid focus session ID")
		return
	}

	session, err := mlh.container.FocusUseCase.PauseFocusSession(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to pause focus session: %v", err))
		return
	}

	mlh.writeSuccess(w, session, "Focus session paused successfully")
}

// HandleResumeFocus handles resuming a focus session
func (mlh *MindloopHandler) HandleResumeFocus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid focus session ID")
		return
	}

	session, err := mlh.container.FocusUseCase.ResumeFocusSession(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to resume focus session: %v", err))
		return
	}

	mlh.writeSuccess(w, session, "Focus session resumed successfully")
}

// HandleRateFocus handles rating a focus session
func (mlh *MindloopHandler) HandleRateFocus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid focus session ID")
		return
	}

	var req FocusRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Rating < 0 || req.Rating > 10 {
		mlh.writeError(w, http.StatusBadRequest, "Rating must be between 0 and 10")
		return
	}

	session, err := mlh.container.FocusUseCase.RateFocusSession(uint(id), req.Rating)
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to rate focus session: %v", err))
		return
	}

	mlh.writeSuccess(w, session, "Focus session rated successfully")
}

// HandleDeleteFocus handles focus session deletion
func (mlh *MindloopHandler) HandleDeleteFocus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid focus session ID")
		return
	}

	err = mlh.container.FocusUseCase.DeleteFocusSession(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete focus session: %v", err))
		return
	}

	mlh.writeSuccess(w, nil, "Focus session deleted successfully")
}

// ===== JOURNAL HANDLERS =====

// HandleCreateJournal handles journal entry creation
func (mlh *MindloopHandler) HandleCreateJournal(w http.ResponseWriter, r *http.Request) {
	var req JournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		mlh.writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.Content == "" {
		mlh.writeError(w, http.StatusBadRequest, "Content is required")
		return
	}

	// Parse mood
	mood := entities.Mood(req.Mood)
	if !entities.IsValidMood(mood) {
		mlh.writeError(w, http.StatusBadRequest, "Invalid mood. Valid moods: happy, sad, neutral, angry, excited")
		return
	}

	entry, err := mlh.container.JournalUseCase.CreateEntry(req.Title, req.Content, mood)
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create journal entry: %v", err))
		return
	}

	mlh.writeSuccess(w, entry, "Journal entry created successfully")
}

// HandleListJournal handles listing journal entries
func (mlh *MindloopHandler) HandleListJournal(w http.ResponseWriter, r *http.Request) {
	entries, err := mlh.container.JournalUseCase.GetAllEntries()
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch journal entries: %v", err))
		return
	}

	mlh.writeSuccess(w, entries, "Journal entries retrieved successfully")
}

// HandleGetJournal handles getting a specific journal entry
func (mlh *MindloopHandler) HandleGetJournal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid journal entry ID")
		return
	}

	entry, err := mlh.container.JournalUseCase.GetEntry(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusNotFound, fmt.Sprintf("Journal entry not found: %v", err))
		return
	}

	mlh.writeSuccess(w, entry, "Journal entry retrieved successfully")
}

// HandleUpdateJournal handles updating a journal entry
func (mlh *MindloopHandler) HandleUpdateJournal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid journal entry ID")
		return
	}

	var req JournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing entry
	entry, err := mlh.container.JournalUseCase.GetEntry(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusNotFound, fmt.Sprintf("Journal entry not found: %v", err))
		return
	}

	// Update fields
	if req.Title != "" {
		entry.Title = req.Title
	}
	if req.Content != "" {
		if err := entry.UpdateContent(req.Content); err != nil {
			mlh.writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid content: %v", err))
			return
		}
	}
	if req.Mood != "" {
		mood := entities.Mood(req.Mood)
		if !entities.IsValidMood(mood) {
			mlh.writeError(w, http.StatusBadRequest, "Invalid mood. Valid moods: happy, sad, neutral, angry, excited")
			return
		}
		if err := entry.UpdateMood(mood); err != nil {
			mlh.writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid mood: %v", err))
			return
		}
	}

	// Save changes
	if err := mlh.container.JournalUseCase.UpdateEntry(entry); err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update journal entry: %v", err))
		return
	}

	mlh.writeSuccess(w, entry, "Journal entry updated successfully")
}

// HandleDeleteJournal handles journal entry deletion
func (mlh *MindloopHandler) HandleDeleteJournal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid journal entry ID")
		return
	}

	err = mlh.container.JournalUseCase.DeleteEntry(uint(id))
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete journal entry: %v", err))
		return
	}

	mlh.writeSuccess(w, nil, "Journal entry deleted successfully")
}

// ===== SUMMARY HANDLERS =====

// HandleDailySummary handles daily summary generation
func (mlh *MindloopHandler) HandleDailySummary(w http.ResponseWriter, r *http.Request) {
	stats, err := mlh.container.SummaryUseCase.GetDailySummary()
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate daily summary: %v", err))
		return
	}

	mlh.writeSuccess(w, stats, "Daily summary generated successfully")
}

// HandleWeeklySummary handles weekly summary generation
func (mlh *MindloopHandler) HandleWeeklySummary(w http.ResponseWriter, r *http.Request) {
	stats, err := mlh.container.SummaryUseCase.GetWeeklySummary()
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate weekly summary: %v", err))
		return
	}

	mlh.writeSuccess(w, stats, "Weekly summary generated successfully")
}

// HandleMonthlySummary handles monthly summary generation
func (mlh *MindloopHandler) HandleMonthlySummary(w http.ResponseWriter, r *http.Request) {
	stats, err := mlh.container.SummaryUseCase.GetMonthlySummary()
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate monthly summary: %v", err))
		return
	}

	mlh.writeSuccess(w, stats, "Monthly summary generated successfully")
}

// HandleYearlySummary handles yearly summary generation
func (mlh *MindloopHandler) HandleYearlySummary(w http.ResponseWriter, r *http.Request) {
	stats, err := mlh.container.SummaryUseCase.GetYearlySummary()
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate yearly summary: %v", err))
		return
	}

	mlh.writeSuccess(w, stats, "Yearly summary generated successfully")
}

// HandleCustomSummary handles custom date range summary generation
func (mlh *MindloopHandler) HandleCustomSummary(w http.ResponseWriter, r *http.Request) {
	var req SummaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		mlh.writeError(w, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD")
		return
	}

	if startDate.After(endDate) {
		mlh.writeError(w, http.StatusBadRequest, "Start date cannot be after end date")
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(24*time.Hour - time.Second)

	stats, err := mlh.container.SummaryUseCase.GenerateSummary(startDate, endDate)
	if err != nil {
		mlh.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate custom summary: %v", err))
		return
	}

	mlh.writeSuccess(w, stats, "Custom summary generated successfully")
}
