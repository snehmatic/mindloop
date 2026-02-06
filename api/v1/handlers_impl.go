package v1

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/motivation"
	"github.com/snehmatic/mindloop/models"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// --- Quote Handler ---

func (mlh *MindloopHandler) HandleQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := motivation.FetchRandomQuote()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching quote")
		http.Error(w, "Failed to fetch quote", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		log.Error().Err(err).Msg("Error encoding quote response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// --- Habit Handlers ---

func (mlh *MindloopHandler) HandleHabitList(w http.ResponseWriter, r *http.Request) {
	interval := r.URL.Query().Get("interval")

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
		Streak      int
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

		streak, _ := mlh.habit.CalculateStreak(h.ID, h.Interval)

		habitViews = append(habitViews, HabitView{
			Habit:       h,
			ActualCount: actual,
			ProgressPct: pct,
			Streak:      streak,
		})
	}

	data := map[string]interface{}{
		"Title":           "Habits",
		"Habits":          habitViews,
		"CurrentInterval": interval,
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

func (mlh *MindloopHandler) HandleHabitUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/habits", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	targetCount, _ := strconv.Atoi(r.FormValue("target_count"))
	interval := r.FormValue("interval")

	h, err := mlh.habit.GetHabit(id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching habit for update")
		http.Redirect(w, r, "/habits?error=Habit not found", http.StatusSeeOther)
		return
	}

	h.Title = title
	h.TargetCount = targetCount
	h.Interval = models.IntervalType(interval)

	if err := mlh.habit.UpdateHabit(h); err != nil {
		log.Error().Err(err).Msg("Error updating habit")
		http.Redirect(w, r, "/habits?error=Failed to update habit", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/habits?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleHabitView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h, err := mlh.habit.GetHabit(id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching habit for view")
		http.Redirect(w, r, "/habits?error=Habit not found", http.StatusSeeOther)
		return
	}

	logs, err := mlh.habit.ListLogsForHabit(h.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching logs for habit")
	}

	// Prepare heatmap data: map[date_string]completion_ratio
	heatmap := make(map[string]float64)
	for _, log := range logs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		ratio := 0.0
		if log.TargetCount > 0 {
			ratio = float64(log.ActualCount) / float64(log.TargetCount)
		}
		if ratio > 1 {
			ratio = 1
		}
		heatmap[dateStr] = ratio
	}

	streak, _ := mlh.habit.CalculateStreak(h.ID, h.Interval)

	mlh.renderTemplate(w, "habit_view.html", map[string]interface{}{
		"Title":   "Habit: " + h.Title,
		"Habit":   h,
		"Heatmap": heatmap,
		"Streak":  streak,
	})
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

func (mlh *MindloopHandler) HandleIntentUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")

	i, err := mlh.intent.GetIntent(id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching intent for update")
		http.Redirect(w, r, "/intent?error=Intent not found", http.StatusSeeOther)
		return
	}

	i.Name = name

	if err := mlh.intent.UpdateIntent(i); err != nil {
		log.Error().Err(err).Msg("Error updating intent")
		http.Redirect(w, r, "/intent?error=Failed to update intent", http.StatusSeeOther)
		return
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

func (mlh *MindloopHandler) HandleIntentDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}
	id := r.FormValue("id")
	if err := mlh.intent.DeleteIntent(id); err != nil {
		log.Error().Err(err).Msg("Error deleting intent")
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

	data := map[string]interface{}{
		"Title":    "Focus",
		"Sessions": sessions,
	}

	if success := r.URL.Query().Get("success"); success == "true" {
		data["SuccessMessage"] = "Action completed successfully"
	}
	if errStr := r.URL.Query().Get("error"); errStr != "" {
		data["ErrorMessage"] = errStr
	}

	mlh.renderTemplate(w, "focus.html", data)
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
		http.Redirect(w, r, "/focus?error="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/focus?success=true", http.StatusSeeOther)
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
		http.Redirect(w, r, "/focus?error="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/focus?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleFocusDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/focus", http.StatusSeeOther)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	if err := mlh.focus.DeleteSession(id); err != nil {
		log.Error().Err(err).Msg("Error deleting focus session")
		http.Redirect(w, r, "/focus?error="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/focus?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleFocusUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/focus", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	title := r.FormValue("title")

	session, err := mlh.focus.GetSession(id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching focus session for update")
		http.Redirect(w, r, "/focus?error=Session not found", http.StatusSeeOther)
		return
	}

	session.Title = title

	if err := mlh.focus.UpdateSession(session); err != nil {
		log.Error().Err(err).Msg("Error updating focus session")
		http.Redirect(w, r, "/focus?error=Failed to update session", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/focus?success=true", http.StatusSeeOther)
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

	// Charts Data
	dailyFocus, labels, _ := mlh.summary.GetFocusSeries(start, now)
	dailyHabits, _ := mlh.summary.GetHabitSeries(start, now)

	mlh.renderTemplate(w, "summary.html", map[string]interface{}{
		"Title":  "Summary",
		"Report": report,
		"Charts": map[string]interface{}{
			"Labels":      labels,
			"DailyFocus":  dailyFocus,
			"DailyHabits": dailyHabits,
		},
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
		err5 := mlh.note.DeleteAll()
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
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
	case "notes":
		err = mlh.note.DeleteAll()
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

func (mlh *MindloopHandler) HandleJournalDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/journal", http.StatusSeeOther)
		return
	}
	id := r.FormValue("id")
	if err := mlh.journal.DeleteEntry(id); err != nil {
		log.Error().Err(err).Msg("Error deleting journal entry")
	}
	http.Redirect(w, r, "/journal", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleJournalUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/journal", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	mood := r.FormValue("mood")

	entry, err := mlh.journal.GetEntry(id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching journal entry for update")
		http.Redirect(w, r, "/journal?error=Entry not found", http.StatusSeeOther)
		return
	}

	entry.Title = title
	entry.Content = content
	entry.Mood = mood

	if err := mlh.journal.UpdateEntry(&entry); err != nil {
		log.Error().Err(err).Msg("Error updating journal entry")
		http.Redirect(w, r, "/journal?error=Failed to update entry", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/journal", http.StatusSeeOther)
}
// --- Note Handlers ---

func (mlh *MindloopHandler) HandleNoteList(w http.ResponseWriter, r *http.Request) {
	notes, err := mlh.note.ListNotes()
	if err != nil {
		log.Error().Err(err).Msg("Error listing notes")
		http.Error(w, "Error fetching notes", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Notes",
		"Notes": notes,
	}

	if success := r.URL.Query().Get("success"); success == "true" {
		data["SuccessMessage"] = "Action completed successfully"
	}
	if errStr := r.URL.Query().Get("error"); errStr != "" {
		data["ErrorMessage"] = errStr
	}

	mlh.renderTemplate(w, "notes.html", data)
}

func (mlh *MindloopHandler) HandleNoteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	labels := r.FormValue("labels")

	_, err := mlh.note.CreateNote(title, content, labels)
	if err != nil {
		log.Error().Err(err).Msg("Error creating note")
		http.Redirect(w, r, "/notes?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/notes?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleNoteView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.Atoi(idStr)

	n, err := mlh.note.GetNote(id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching note")
		http.Redirect(w, r, "/notes?error=Note not found", http.StatusSeeOther)
		return
	}

	htmlContent := mdToHTML([]byte(n.Content))

	mlh.renderTemplate(w, "note_view.html", map[string]interface{}{
		"Title":       "View Note",
		"Note":        n,
		"HTMLContent": template.HTML(htmlContent),
	})
}

func (mlh *MindloopHandler) HandleNoteDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)

	err := mlh.note.DeleteNote(id)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting note")
		http.Redirect(w, r, "/notes?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/notes?success=true", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	mlh.renderTemplate(w, "about.html", map[string]interface{}{
		"Title": "About",
	})
}

func (mlh *MindloopHandler) HandleVoid(w http.ResponseWriter, r *http.Request) {
	mlh.renderTemplate(w, "void.html", map[string]interface{}{
		"Title": "The Void",
	})
}

// --- Settings Handlers ---

func (mlh *MindloopHandler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	uc := config.UserConfig{}
	_ = uc.ReadFromYAML() // Ignore error if file doesn't exist

	data := map[string]interface{}{
		"Title":    "Settings",
		"Config":   uc,
		"AllModes": config.AllModes,
	}

	if success := r.URL.Query().Get("success"); success == "true" {
		data["SuccessMessage"] = "Settings updated successfully"
	}
	if errStr := r.URL.Query().Get("error"); errStr != "" {
		data["ErrorMessage"] = errStr
	}

	mlh.renderTemplate(w, "settings.html", data)
}

func (mlh *MindloopHandler) HandleSettingsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	mode := r.FormValue("mode")

	uc := config.UserConfig{
		Name: name,
		Mode: mode,
		FeatureFlags: config.FeatureFlags{
			FocusCloud:   r.FormValue("focus_cloud") == "on",
			HabitCloud:   r.FormValue("habit_cloud") == "on",
			IntentCloud:  r.FormValue("intent_cloud") == "on",
			JournalCloud: r.FormValue("journal_cloud") == "on",
			NoteCloud:    r.FormValue("note_cloud") == "on",
		},
	}

	if mode == "byodb" {
		uc.DbConfig = config.DBConfig{
			Host:     r.FormValue("db_host"),
			Port:     r.FormValue("db_port"),
			User:     r.FormValue("db_user"),
			Password: r.FormValue("db_pass"),
			Name:     r.FormValue("db_name"),
		}
	}

	uc.WriteToYAML()

	http.Redirect(w, r, "/settings?success=true", http.StatusSeeOther)
}

// --- Backup Handlers ---

func (mlh *MindloopHandler) HandleBackupExport(w http.ResponseWriter, r *http.Request) {
	tmpFile, err := os.CreateTemp("", "mindloop_backup_*.json")
	if err != nil {
		log.Error().Err(err).Msg("Error creating temp file for backup")
		http.Redirect(w, r, "/settings?error=Failed to create backup file", http.StatusSeeOther)
		return
	}
	defer os.Remove(tmpFile.Name())

	if err := mlh.backup.Export(tmpFile.Name()); err != nil {
		log.Error().Err(err).Msg("Error exporting backup")
		http.Redirect(w, r, "/settings?error=Failed to export data", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=mindloop_backup.json")
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, tmpFile.Name())
}

func (mlh *MindloopHandler) HandleBackupImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	file, _, err := r.FormFile("backup_file")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving uploaded file")
		http.Redirect(w, r, "/settings?error=No file uploaded", http.StatusSeeOther)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "mindloop_import_*.json")
	if err != nil {
		log.Error().Err(err).Msg("Error creating temp file for import")
		http.Redirect(w, r, "/settings?error=Import failed", http.StatusSeeOther)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Copy uploaded file to temp file
	data, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("Error reading uploaded file")
		http.Redirect(w, r, "/settings?error=Failed to read upload", http.StatusSeeOther)
		return
	}

	if err := os.WriteFile(tmpFile.Name(), data, 0644); err != nil {
		log.Error().Err(err).Msg("Error writing temp file")
		http.Redirect(w, r, "/settings?error=Import failed", http.StatusSeeOther)
		return
	}

	if err := mlh.backup.Import(tmpFile.Name()); err != nil {
		log.Error().Err(err).Msg("Error importing data")
		http.Redirect(w, r, "/settings?error=Restore failed: "+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/settings?success=true", http.StatusSeeOther)
}
