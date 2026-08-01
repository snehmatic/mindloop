package v1

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/backup"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/note"
	"github.com/snehmatic/mindloop/internal/core/quest"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/core/task"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/snehmatic/mindloop/web"
)

type HabitView struct {
	models.Habit
	ActualCount int
	ProgressPct int
	Streak      int
}

var templateFuncs = template.FuncMap{
	"json": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
	"title": func(v any) string {
		s := fmt.Sprint(v)
		if len(s) == 0 {
			return ""
		}
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	},
	"iso8601": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	"split": func(s, sep string) []string {
		if s == "" {
			return nil
		}
		parts := strings.Split(s, sep)
		var trimmed []string
		for _, p := range parts {
			t := strings.TrimSpace(p)
			if t != "" {
				trimmed = append(trimmed, t)
			}
		}
		return trimmed
	},
	"asset": func(path string) string {
		cfg := config.GetConfig()
		if cfg.Mode == config.Local {
			return fmt.Sprintf("%s?v=%d", path, time.Now().Unix())
		}
		return fmt.Sprintf("%s?v=1.0.0", path)
	},
}

type MindloopHandler struct {
	config  *config.Config
	journal *journal.Service
	note    *note.Service
	habit   *habit.Service
	focus   *focus.Service
	intent  *intent.Service
	quest   *quest.Service
	summary *summary.Service
	backup  *backup.Service
	task    *task.Service
}

func NewMindloopHandler(
	journal *journal.Service,
	note *note.Service,
	habit *habit.Service,
	focus *focus.Service,
	intent *intent.Service,
	quest *quest.Service,
	summary *summary.Service,
	backup *backup.Service,
	task *task.Service,
) *MindloopHandler {
	return &MindloopHandler{
		config:  config.GetConfig(),
		journal: journal,
		note:    note,
		habit:   habit,
		focus:   focus,
		intent:  intent,
		quest:   quest,
		summary: summary,
		backup:  backup,
		task:    task,
	}
}

func (mlh *MindloopHandler) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	files := []string{
		"templates/layout.html",
		"templates/" + tmpl,
	}

	// Also include partials that might be used by the main template
	if tmpl == "focus.html" {
		files = append(files, "templates/focus_active_timer.html", "templates/focus_session_list.html")
	}
	if tmpl == "habits.html" {
		files = append(files, "templates/_habit_card.html")
	}

	ts := template.New("layout.html").Funcs(templateFuncs)

	ts, err := ts.ParseFS(web.WebFS, files...)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing templates")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var buf strings.Builder

	// Inject global view data (like UserName) if data is a map
	if d, ok := data.(map[string]interface{}); ok {
		if mlh.config.UserName != "" {
			d["UserName"] = mlh.config.UserName
		}
		if _, exists := d["Config"]; !exists {
			uc := config.UserConfig{}
			_ = uc.ReadFromYAML()
			d["Config"] = uc
		}
	}

	err = ts.Execute(&buf, data)
	if err != nil {
		log.Error().Err(err).Msg("Error executing template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprint(w, buf.String())
}

func (mlh *MindloopHandler) renderPartial(w http.ResponseWriter, tmpl string, data interface{}) {
	files := []string{
		"templates/" + tmpl,
	}

	// We use the filename as the template name for single file parsing
	ts := template.New(tmpl).Funcs(templateFuncs)

	ts, err := ts.ParseFS(web.WebFS, files...)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing partial template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var buf strings.Builder
	err = ts.Execute(&buf, data)
	if err != nil {
		log.Error().Err(err).Msg("Error executing partial template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprint(w, buf.String())
}

func (mlh *MindloopHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	// 1. Ongoing Intent (Active or Paused)
	currentIntent, _ := mlh.intent.GetOngoingIntent()

	// 1b. Active Side Quest
	var currentQuest *models.SideQuest
	if q, err := mlh.quest.GetActiveQuest(); err == nil {
		currentQuest = q
	}

	// 2. Focus Time Today
	now := time.Now()
	todayStart := now.Truncate(24 * time.Hour)
	focusStats, _ := mlh.summary.GetFocusStats(todayStart, now)
	focusMinutes := int(focusStats.RawDuration)

	// 3. Habit Progress (Completed Today / Total Active Daily Habits)
	habits, _ := mlh.habit.ListHabits(models.Daily)
	habitLogs, _ := mlh.habit.ListHabitLogs(models.Daily)

	completedHabits := 0
	totalHabits := 0

	for _, h := range habits {
		totalHabits++
		// Check if completed today
		for _, l := range habitLogs {
			if l.HabitID == h.ID && l.CreatedAt.After(todayStart) && l.ActualCount >= h.TargetCount {
				completedHabits++

				break
			}
		}
	}

	// Build Active Habits for Dashboard
	var activeHabits []HabitView
	for _, h := range habits {
		actual := 0
		for _, log := range habitLogs {
			if log.HabitID == h.ID && log.CreatedAt.After(todayStart) {
				actual = log.ActualCount
				break
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
		activeHabits = append(activeHabits, HabitView{
			Habit:       h,
			ActualCount: actual,
			ProgressPct: pct,
			Streak:      streak,
		})
	}

	// 4. Pending Tasks
	allTasks, _ := mlh.task.ListTasks()
	var pendingTasks []models.TaskView
	for _, t := range allTasks {
		if t.Status == "pending" {
			pendingTasks = append(pendingTasks, models.ToTaskView(t))
		}
	}

	// 5. Active Focus
	activeFocus, _ := mlh.focus.GetActiveSession()
	mlh.renderTemplate(w, "home.html", map[string]interface{}{
		"Title": "Home",
		"Dashboard": map[string]interface{}{
			"CurrentIntent":   currentIntent,
			"CurrentQuest":    currentQuest,
			"FocusMinutes":    focusMinutes,
			"CompletedHabits": completedHabits,
			"TotalHabits":     totalHabits,
			"ActiveHabits":    activeHabits,
			"PendingTasks":    pendingTasks,
			"ActiveFocus":     activeFocus,
		},
	})
}

func (mlh *MindloopHandler) HandleJournalList(w http.ResponseWriter, r *http.Request) {
	entries, err := mlh.journal.ListEntries()
	if err != nil {
		log.Error().Err(err).Msg("Error listing journal entries")
		http.Error(w, "Error fetching entries", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Journal",
		"Entries": entries,
	}

	if success := r.URL.Query().Get("success"); success != "" {
		switch success {
		case "true":
			data["SuccessMessage"] = "Action completed successfully"
		case "done":
			data["SuccessMessage"] = "Journal entry saved! Great reflection!"
		case "milestone":
			data["SuccessMessage"] = "🏆 MILESTONE REACHED! You are amazing! 🏆"
		}
	}

	mlh.renderTemplate(w, "journal.html", data)
}

func (mlh *MindloopHandler) HandleJournalCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/journal", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	mood := r.FormValue("mood")

	uc := config.UserConfig{}
	_ = uc.ReadFromYAML()

	milestoneReached, err := mlh.journal.CreateEntry(title, content, mood, uc.PointsConfig.Journal)
	if err != nil {
		log.Error().Err(err).Msg("Error creating journal entry")
		// In a real app, we'd pass the error back to the template
	}

	if uc.FeatureFlags.Gamification {
		if r.Header.Get("HX-Request") == "true" {
			if milestoneReached {
				w.Header().Set("HX-Trigger", "{\"milestone\": {}}")
			} else {
				w.Header().Set("HX-Trigger", "{\"confetti\": {}}")
			}
		} else {
			successType := "done"
			if milestoneReached {
				successType = "milestone"
			}
			http.Redirect(w, r, "/journal?success="+successType, http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/journal", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse([]byte("OK"), w, http.StatusOK)
}

// --- Quest Handlers ---

func (mlh *MindloopHandler) HandleQuestStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")

	// 1. Pause Intent
	currentIntent, _ := mlh.intent.GetOngoingIntent()
	if currentIntent != nil && currentIntent.Status == "active" {
		_, _ = mlh.intent.PauseIntent(currentIntent.ID)
	}

	// 2. Pause Focus
	activeFocus, _ := mlh.focus.GetActiveSession()
	if activeFocus != nil {
		_, _ = mlh.focus.PauseSession(activeFocus.ID)
	}

	// 3. Start Quest
	_, _ = mlh.quest.StartQuest(title)

	http.Redirect(w, r, "/intent", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleQuestStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	note := r.FormValue("note")

	uc := config.UserConfig{}
	_ = uc.ReadFromYAML()

	_, milestoneReached, _ := mlh.quest.StopQuest(uint(id), note, uc.PointsConfig.Quest)

	// Auto-resume intent if one is paused
	currentIntent, _ := mlh.intent.GetOngoingIntent()
	if currentIntent != nil && currentIntent.Status == "paused" {
		_, _ = mlh.intent.ResumeIntent(currentIntent.ID)
	}

	if uc.FeatureFlags.Gamification {
		if r.Header.Get("HX-Request") == "true" {
			if milestoneReached {
				w.Header().Set("HX-Trigger", "{\"milestone\": {}}")
			} else {
				w.Header().Set("HX-Trigger", "{\"confetti\": {}}")
			}
		} else {
			successType := "done"
			if milestoneReached {
				successType = "milestone"
			}
			http.Redirect(w, r, "/intent?success="+successType, http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/intent", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleQuestDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	_ = mlh.quest.DeleteQuest(uint(id))

	http.Redirect(w, r, "/intent", http.StatusSeeOther)
}

func (mlh *MindloopHandler) HandleIntentResume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/intent", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	// 1. Resume the intent
	_, _ = mlh.intent.ResumeIntent(uint(id))

	// 2. Automatically complete any active side quest
	activeQuest, _ := mlh.quest.GetActiveQuest()
	if activeQuest != nil {
		uc := config.UserConfig{}
		_ = uc.ReadFromYAML()
		_, _, _ = mlh.quest.StopQuest(activeQuest.ID, "Resumed main intent", uc.PointsConfig.Quest)
	}

	http.Redirect(w, r, "/intent", http.StatusSeeOther)
}
