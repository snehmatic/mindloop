package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/models"
)

// --- Habit Handlers ---

func (mlh *MindloopHandler) HandleHabitList(w http.ResponseWriter, r *http.Request) {
	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = string(models.Daily)
	}

	habits, err := mlh.habit.ListHabits(models.IntervalType(interval))
	if err != nil {
		log.Error().Err(err).Msg("Error listing habits")
		http.Error(w, "Error fetching habits", http.StatusInternalServerError)
		return
	}

	habitLogs, err := mlh.habit.ListHabitLogs(models.IntervalType(interval))
	if err != nil {
		log.Error().Err(err).Msg("Error listing habit logs")
	}

	// Calculate completion for UI
	type HabitView struct {
		models.Habit
		ActualCount int
		ProgressPct int
	}

	var habitViews []HabitView
	for _, h := range habits {
		actual := 0
		for _, log := range habitLogs {
			// Basic match for today/current interval - simplified logic for UI
			if log.HabitID == h.ID {
				// Check if the log is "current" (today for daily)
				// Simplify: just taking the log count if it matches.
				// In a real app, `ListHabitLogs` should filter by date range or we filter here.
				// For now, let's assume `ListHabitLogs` returns all, but we only really care about "current" status
				// This part ideally needs the service to return "HabitWithStatus".
				// Re-using service logic:
				// We'll iterate and find if there's a log for *today* (created_at)
				isToday := false
				if h.Interval == models.Daily {
					if log.CreatedAt.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour)) {
						isToday = true
					}
				} else {
					// Weekly check simplified
					// For weekly, we check if created_at is within this week.
					year, week := time.Now().ISOWeek()
					logYear, logWeek := log.CreatedAt.ISOWeek()
					if year == logYear && week == logWeek {
						isToday = true
					}
				}

				if isToday {
					actual = log.ActualCount
					break
				}
			}
		}
		pct := 0
		if h.TargetCount > 0 {
			pct = (actual * 100) / h.TargetCount
		}
		if pct > 100 {
			pct = 100
		}
		habitViews = append(habitViews, HabitView{
			Habit:       h,
			ActualCount: actual,
			ProgressPct: pct,
		})
	}

	data := map[string]interface{}{
		"Title":  "Habits",
		"Habits": habitViews,
	}

	// Pass query params as simple alerts
	if success := r.URL.Query().Get("success"); success == "true" {
		data["SuccessMessage"] = "Action completed successfully"
	}
	if errStr := r.URL.Query().Get("error"); errStr != "" {
		data["ErrorMessage"] = errStr
	}

	mlh.renderTemplate(w, "habits.html", data)
}

func (mlh *MindloopHandler) HandleHabitCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/habits", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	targetCount, _ := strconv.Atoi(r.FormValue("target_count"))
	interval := r.FormValue("interval")

	habit := &models.Habit{
		Title:       title,
		TargetCount: targetCount,
		Interval:    models.IntervalType(interval),
	}

	if err := mlh.habit.CreateHabit(habit); err != nil {
		log.Error().Err(err).Msg("Error creating habit")
		http.Redirect(w, r, "/habits?error=Failed to create habit", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/habits?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleHabitLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/habits", http.StatusSeeOther)
		return
	}

	habitID := r.FormValue("habit_id")
	_, _, err := mlh.habit.LogHabit(habitID)
	if err != nil {
		log.Error().Err(err).Msg("Error logging habit")
		http.Redirect(w, r, "/habits?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/habits?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleHabitUnlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/habits", http.StatusSeeOther)
		return
	}

	habitID := r.FormValue("habit_id")
	_, err := mlh.habit.UnlogHabit(habitID)
	if err != nil {
		log.Error().Err(err).Msg("Error Unlogging habit")
		http.Redirect(w, r, "/habits?error="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/habits?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleHabitDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/habits", http.StatusSeeOther)
		return
	}

	habitID := r.FormValue("habit_id")
	err := mlh.habit.DeleteHabit(habitID)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting habit")
		http.Redirect(w, r, "/habits?error="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/habits?success=true", http.StatusSeeOther)
}

// --- Intent Handlers ---

func (mlh *MindloopHandler) HandleIntent(w http.ResponseWriter, r *http.Request) {
	activeIntents, _ := mlh.intent.ListActiveIntents()
	allIntents, _ := mlh.intent.ListIntents()

	var currentIntent *models.Intent
	if len(activeIntents) > 0 {
		currentIntent = &activeIntents[0] // Just take the first active one
	}

	mlh.renderTemplate(w, "intent.html", map[string]interface{}{
		"Title":         "Intent",
		"CurrentIntent": currentIntent,
		"History":       allIntents,
	})
}

func (mlh *MindloopHandler) HandleIntentSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	_, err := mlh.intent.StartIntent(name)
	if err != nil {
		log.Error().Err(err).Msg("Error setting intent")
	}

	http.Redirect(w, r, "/intent", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleIntentComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}
	id := r.FormValue("id")
	_, err := mlh.intent.EndIntent(id)
	if err != nil {
		log.Error().Err(err).Msg("Error completing intent")
	}
	http.Redirect(w, r, "/intent", http.StatusSeeOther)
}

// --- Focus Handlers ---

func (mlh *MindloopHandler) HandleFocus(w http.ResponseWriter, r *http.Request) {
	sessions, _ := mlh.focus.ListSessions()
	// reverse order to show newest first
	for i, j := 0, len(sessions)-1; i < j; i, j = i+1, j-1 {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	}

	mlh.renderTemplate(w, "focus.html", map[string]interface{}{
		"Title":    "Focus",
		"Sessions": sessions,
	})
}

func (mlh *MindloopHandler) HandleFocusStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/focus", http.StatusSeeOther)
		return
	}
	title := r.FormValue("title")
	_, err := mlh.focus.StartSession(title)
	if err != nil {
		log.Error().Err(err).Msg("Error starting focus session")
	}
	http.Redirect(w, r, "/focus", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleFocusStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/focus", http.StatusSeeOther)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	_, err := mlh.focus.EndSession(id)
	if err != nil {
		log.Error().Err(err).Msg("Error ending focus session")
	}
	http.Redirect(w, r, "/focus", http.StatusSeeOther)
}

// --- Summary Handler ---

func (mlh *MindloopHandler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	start := now.AddDate(0, 0, -7) // Default to Last 7 days

	// Parse Custom Range
	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")

	if startParam != "" {
		if parsedStart, err := time.Parse("2006-01-02", startParam); err == nil {
			start = parsedStart
		}
	}
	if endParam != "" {
		if parsedEnd, err := time.Parse("2006-01-02", endParam); err == nil {
			// Set end to end of that day
			now = parsedEnd.Add(24*time.Hour - 1*time.Second)
		}
	}

	report, err := mlh.summary.GenerateSummary(start, now)
	if err != nil {
		log.Error().Err(err).Msg("Error generating summary")
		// Render with error message
		mlh.renderTemplate(w, "summary.html", map[string]interface{}{
			"Title":        "Summary",
			"ErrorMessage": "Failed to generate summary: " + err.Error(),
			"Report":       models.SummaryReport{DateRange: "Unavailable"},
		})
		return
	}

	mlh.renderTemplate(w, "summary.html", map[string]interface{}{
		"Title":  "Summary",
		"Report": report,
	})
}

func (mlh *MindloopHandler) HandleCleanSlate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	target := r.FormValue("type") // all, journal, habits, focus, intent
	var err error

	switch target {
	case "all", "": // Default to all if empty
		// Order matters for FK constraints if any, though we don't have many
		err1 := mlh.journal.DeleteAll()
		err2 := mlh.habit.DeleteAll()
		err3 := mlh.focus.DeleteAll()
		err4 := mlh.intent.DeleteAll()
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			err = errors.New("failed to delete all data") // Force non-nil error if any failed
			log.Error().Msg("Error in clean slate all")
		}
	case "journal":
		err = mlh.journal.DeleteAll()
	case "habits":
		err = mlh.habit.DeleteAll()
	case "focus":
		err = mlh.focus.DeleteAll()
	case "intent":
		err = mlh.intent.DeleteAll()
	default:
		// Unknown type
	}

	redirectURL := "/"
	if target != "all" {
		redirectURL = "/" + target // e.g. /journal, /habits
	}

	if err != nil {
		http.Redirect(w, r, redirectURL+"?error=Failed to reset data", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, redirectURL+"?success=Data cleared successfully", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	mlh.renderTemplate(w, "about.html", map[string]interface{}{
		"Title": "About",
	})
}
