package v1

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/internal/core/journal"
	"github.com/snehmatic/mindloop/internal/core/summary"
	"github.com/snehmatic/mindloop/internal/utils"
)

type MindloopHandler struct {
	config  *config.Config
	journal *journal.Service
	habit   *habit.Service
	focus   *focus.Service
	intent  *intent.Service
	summary *summary.Service
}

func NewMindloopHandler(
	journal *journal.Service,
	habit *habit.Service,
	focus *focus.Service,
	intent *intent.Service,
	summary *summary.Service,
) *MindloopHandler {
	return &MindloopHandler{
		config:  config.GetConfig(),
		journal: journal,
		habit:   habit,
		focus:   focus,
		intent:  intent,
		summary: summary,
	}
}

func (mlh *MindloopHandler) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	cwd, _ := filepath.Abs(".")
	// Define the base layout and the specific template
	basePath := filepath.Join(cwd, "web/templates")
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// Try going up levels for tests
		if _, err := os.Stat(filepath.Join(cwd, "../web/templates")); err == nil {
			cwd = filepath.Join(cwd, "..")
		} else if _, err := os.Stat(filepath.Join(cwd, "../../web/templates")); err == nil {
			cwd = filepath.Join(cwd, "../..")
		}
	}

	files := []string{
		filepath.Join(cwd, "web/templates/layout.html"),
		filepath.Join(cwd, "web/templates/", tmpl),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing templates")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Error().Err(err).Msg("Error executing template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (mlh *MindloopHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	// Gather Dashboard Stats
	// 1. Active Habits
	habits, _ := mlh.habit.ListHabits("")
	activeHabits := len(habits)

	// 2. Focus Time Today
	now := time.Now()
	todayStart := now.Truncate(24 * time.Hour)
	focusStats, _ := mlh.summary.GetFocusStats(todayStart, now)

	// 3. Last Journal Mood
	entries, _ := mlh.journal.ListEntries()
	lastMood := "N/A"
	if len(entries) > 0 {
		lastMood = entries[0].Mood // Assuming sorted by desc
	}

	mlh.renderTemplate(w, "home.html", map[string]interface{}{
		"Title": "Home",
		"Stats": map[string]interface{}{
			"ActiveHabits": activeHabits,
			"FocusTime":    focusStats.TotalDuration,
			"LastMood":     lastMood,
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
