# Issue 2: [CLI] Interactive TUI Dashboard (`mindloop dash`)

**Labels:** `enhancement`, `cli`, `tui`

### Rationale & Market Fit
Currently, the CLI is transactional (run command, get output, exit). Many developers prefer to keep a persistent terminal pane open with their daily dashboard (like `k9s` or `lazydocker`). A TUI (Text User Interface) brings the beautiful Web UI experience back to the terminal.

### Expected Outcomes
- Running `mindloop dash` opens a full-screen terminal dashboard.
- The dashboard displays the Active Intent (large and bold), Daily Habits, and a timer for any active Focus Session.
- Keyboard shortcuts allow interacting with the dashboard without leaving it.

### Implementation Details
1. **Framework:** Use [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) and [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) for layout and styling.
2. **Architecture:**
   - Add a new cobra command in `cmd/cli/dash.go`.
   - The TUI should fetch data via the existing `internal/core` domain logic (no need to hit the HTTP API).
3. **Features:**
   - Show Active Intent prominently.
   - List Habits for today (toggle with `Space`).
   - Quit with `q`.
