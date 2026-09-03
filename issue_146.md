The frontend dashboard (`/`) currently acts only as a static summary and navigation hub. It requires users to navigate away to perform basic actions like logging a daily habit or ticking off a pending task.

**Proposed Changes:**
1. **Interactive Habit Tracker:** Show the actual list of today's habits on the dashboard with HTMX-powered checkboxes to log them instantly.
2. **Pending Tasks Widget:** Display up to 5 pending tasks on the dashboard with 1-click completion.
3. **Active Focus Indicator:** If a focus session is active, show the timer or a prominent banner to stop/pause it.
4. **Remove Redundant Quick Access:** Since we have the Command Palette (`Cmd+K`) and the Navbar, the 6 "Quick Access" cards take up valuable vertical real estate. Replace them with the functional widgets mentioned above.

**Implementation Details:**
- Update `api/v1/handlers.go` (`HandleHome`) to fetch and pass `[]models.HabitView`, `[]models.TaskView`, and the `CurrentFocus` session.
- Update `web/templates/home.html` to render these widgets. Use HTMX for instant, no-reload interactions (e.g. `hx-post="/tasks/complete"`).
