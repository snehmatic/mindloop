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
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"github.com/snehmatic/mindloop/web"
)

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
	}
}

func (mlh *MindloopHandler) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	files := []string{
		"templates/layout.html",
		"templates/" + tmpl,
	}

	ts := template.New("layout.html").Funcs(template.FuncMap{
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		"title": func(v interface{}) string {
			s := fmt.Sprintf("%v", v)
			if s == "" {
				return ""
			}
			return strings.Title(s)
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
	})

	ts, err := ts.ParseFS(web.WebFS, files...)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing templates")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var buf strings.Builder
	err = ts.Execute(&buf, data)
	if err != nil {
		log.Error().Err(err).Msg("Error executing template")
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

	mlh.renderTemplate(w, "home.html", map[string]interface{}{
		"Title": "Home",
		"Dashboard": map[string]interface{}{
			"CurrentIntent":   currentIntent,
			"CurrentQuest":    currentQuest,
			"FocusMinutes":    focusMinutes,
			"CompletedHabits": completedHabits,
			"TotalHabits":     totalHabits,
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

	mlh.renderTemplate(w, "journal.html", map[string]interface{}{
		"Title":   "Journal",
		"Entries": entries,
	})
}

func (mlh *MindloopHandler) HandleJournalCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/journal", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	mood := r.FormValue("mood")

	if err := mlh.journal.CreateEntry(title, content, mood); err != nil {
		log.Error().Err(err).Msg("Error creating journal entry")
		// In a real app, we'd pass the error back to the template
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

	_, _ = mlh.quest.StopQuest(uint(id), note)

	// Auto-resume intent if one is paused
	currentIntent, _ := mlh.intent.GetOngoingIntent()
	if currentIntent != nil && currentIntent.Status == "paused" {
		_, _ = mlh.intent.ResumeIntent(currentIntent.ID)
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

	_, _ = mlh.intent.ResumeIntent(uint(id))

	http.Redirect(w, r, "/intent", http.StatusSeeOther)
}
