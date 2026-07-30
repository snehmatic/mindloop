package cli

import (
	"fmt"
	"math"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cfg "github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/core/focus"
	"github.com/snehmatic/mindloop/internal/core/habit"
	"github.com/snehmatic/mindloop/internal/core/intent"
	"github.com/snehmatic/mindloop/models"
	"github.com/spf13/cobra"
)

var (
	titleStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 2).MarginBottom(1)
	intentStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF79C6")).MarginBottom(1)
	focusStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F1FA8C")).MarginBottom(1)
	habitStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	selectedHabitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	doneStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4")).Strikethrough(true)
)

type tickMsg time.Time

type dashboardModel struct {
	activeIntent string
	activeFocus  *models.FocusSession
	habits       []models.Habit
	habitLogs    map[uint]bool // habitID -> completed today
	cursor       int
	uc           cfg.UserConfig
}

func initialModel() dashboardModel {
	hService := habit.NewService(gdb)
	iService := intent.NewService(gdb)
	fService := focus.NewService(gdb)

	var activeIntent string
	intents, _ := iService.ListActiveIntents()
	if len(intents) > 0 {
		activeIntent = intents[0].Name
	} else {
		activeIntent = "No active intent. Use 'mindloop intent start' to set one."
	}

	activeFocus, _ := fService.GetActiveSession()

	habits, _ := hService.ListHabits(models.Daily)
	logs, _ := hService.ListHabitLogs(models.Daily)

	habitLogs := make(map[uint]bool)
	for _, l := range logs {
		if l.ActualCount >= l.TargetCount {
			habitLogs[l.HabitID] = true
		}
	}

	userCfg := cfg.UserConfig{}
	_ = userCfg.ReadFromYAML()

	return dashboardModel{
		activeIntent: activeIntent,
		activeFocus:  activeFocus,
		habits:       habits,
		habitLogs:    habitLogs,
		uc:           userCfg,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m dashboardModel) Init() tea.Cmd {
	return tickCmd()
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.habits)-1 {
				m.cursor++
			}
		case " ":
			if len(m.habits) > 0 {
				h := m.habits[m.cursor]
				hService := habit.NewService(gdb)
				idStr := fmt.Sprintf("%d", h.ID)
				if m.habitLogs[h.ID] {
					_, _ = hService.UnlogHabit(idStr)
					m.habitLogs[h.ID] = false
				} else {
					_, _, _, _ = hService.LogHabit(idStr, m.uc.PointsConfig.Habit)
					m.habitLogs[h.ID] = true
				}
			}
		}
	}
	return m, nil
}

func (m dashboardModel) View() string {
	s := "\n" + titleStyle.Render("Mindloop Dashboard") + "\n\n"

	s += "Active Intent:\n"
	s += intentStyle.Render(m.activeIntent) + "\n\n"

	if m.activeFocus != nil {
		duration := time.Since(m.activeFocus.CreatedAt)
		mins := int(math.Floor(duration.Minutes()))
		secs := int(duration.Seconds()) % 60
		timer := fmt.Sprintf("[%02d:%02d] %s", mins, secs, m.activeFocus.Title)
		s += "Active Focus Session:\n"
		s += focusStyle.Render(timer) + "\n\n"
	}

	s += "Daily Habits (Space to toggle):\n"
	for i, h := range m.habits {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.habitLogs[h.ID] {
			checked = "x"
		}

		line := fmt.Sprintf("%s [%s] %s", cursor, checked, h.Title)
		if m.habitLogs[h.ID] {
			s += doneStyle.Render(line) + "\n"
		} else if m.cursor == i {
			s += selectedHabitStyle.Render(line) + "\n"
		} else {
			s += habitStyle.Render(line) + "\n"
		}
	}

	s += "\nPress q to quit."
	return s
}

var dashCmd = &cobra.Command{
	Use:   "dash",
	Short: "Open the interactive TUI dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running dashboard: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(dashCmd)
}
