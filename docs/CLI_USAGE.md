# Mindloop CLI Usage Guide

This guide provides a detailed overview of the `mindloop` CLI commands and how to use them to manage your productivity.

## Getting Started

Before using the CLI, it is recommended to run the configuration command to set up your profile and database preferences.

```bash
mindloop configure
```

This interactive command will guide you through setting up:
*   **Username**: Your display name.
*   **Mode**: Choose between `local` (SQLite) or `byodb` (Bring Your Own Database, e.g., PostgreSQL).

---

## Intent Management

Intents help you set a single-threaded goal to maintain focus.

### Commands

*   **Start an intent:**
    ```bash
    mindloop intent start "Implement feature X"
    ```
    This sets your current active intent.

*   **Check current intent:**
    ```bash
    mindloop intent current
    ```
    Shows the intent you are currently working on.

*   **List all intents:**
    ```bash
    mindloop intent list
    ```
    Displays a history of all your logged intents.

*   **End an intent:**
    ```bash
    mindloop intent end <intent_id>
    ```
    Marks the specific intent as completed. You need to provide the ID of the intent (found via `current` or `list`).

---

## Side Quests

Side quests manage ad-hoc tasks that interrupt your main flow. Starting a quest automatically pauses your active intent and focus session.

### Commands

*   **Start a side quest:**
    ```bash
    mindloop quest start "Investigate outage"
    ```
    This pauses your current intent/focus and starts a new quest.

*   **List quests:**
    ```bash
    mindloop quest list
    ```

*   **Stop a quest:**
    ```bash
    mindloop quest stop "Fixed the issue"
    ```
    Completes the quest. You can optionally provide a closing note.

*   **Delete a quest:**
    ```bash
    mindloop quest delete <quest_id>
    ```

---

## Focus Sessions

Focus sessions allow you to track time spent on specific tasks, helping you stay in the "zone".

### Commands

*   **Start a session:**
    ```bash
    mindloop focus start "Deep work on API"
    ```
    Starts a timer for your focus session.

*   **List sessions:**
    ```bash
    mindloop focus list
    ```
    Shows a list of all focus sessions recorded.

*   **End a session:**
    ```bash
    mindloop focus end <session_id>
    ```
    Stops the timer and records the session duration.

*   **Rate a session:**
    ```bash
    mindloop focus rate <session_id> <rating>
    ```
    Rate the quality of your focus session from 0 to 10 (e.g., `mindloop focus rate 12 9`).

---

## Habit Tracking

Build consistency by tracking daily and weekly habits.

### Commands

*   **Add a habit:**
    ```bash
    # Interactive mode (Recommended)
    mindloop habit add -i

    # Quick add (Daily default)
    mindloop habit add "Read Book" "Read 10 pages" 1

    # Add a weekly habit
    mindloop habit add "Gym" "Workout leg day" 3 --weekly
    ```

*   **List habits:**
    ```bash
    mindloop habit list
    # Filter by interval
    mindloop habit list --daily
    mindloop habit list --weekly
    ```

*   **Log a habit (Mark as done):**
    ```bash
    mindloop habit log <habit_id>
    ```
    Increments the count for the habit for the current period. Use aliases `done` or `complete` if preferred.

*   **Unlog a habit:**
    ```bash
    mindloop habit unlog <habit_id>
    ```
    Decrements the count (undo).

*   **Check habit status:**
    ```bash
    mindloop habit show
    mindloop habit show --weekly
    ```
    Displays your progress (logs) for the current day or week.

*   **Update a habit:**
    ```bash
    mindloop habit update <habit_id>
    ```
    Enters interactive mode to modify habit details.

*   **Delete a habit:**
    ```bash
    mindloop habit delete <habit_id>
    ```

---

## Journaling

Capture your thoughts and reflections directly from the terminal.

### Commands

*   **New entry:**
    ```bash
    mindloop journal new "End of day reflection"
    # Optional mood flag
    mindloop journal new "Great day" --mood "happy"
    ```
    This command opens your default `$EDITOR` (e.g., vim, nano) to write the content.

*   **List entries:**
    ```bash
    mindloop journal list
    ```
    Shows a history of journal entries.

*   **View an entry:**
    ```bash
    mindloop journal view <entry_id>
    ```
    Displays the full content of a specific entry.

*   **Delete an entry:**
    ```bash
    mindloop journal delete <entry_id>
    ```

---

## Quick Notes

Capture ideas, snippets, and important information quickly without context switching.

### Commands

*   **Quick capture:**
    ```bash
    mindloop note "Meeting with team at 3 PM"
    ```

*   **New note (Editor):**
    ```bash
    mindloop note new --title "Brainstorm" --labels "ideas,work"
    ```
    Opens your default editor to write a longer note.

*   **List notes:**
    ```bash
    mindloop note list
    ```

*   **View a note:**
    ```bash
    mindloop note view <note_id>
    ```

*   **Edit a note:**
    ```bash
    mindloop note edit <note_id>
    ```

*   **Delete a note:**
    ```bash
    mindloop note delete <note_id>
    ```

---

## Tasks

Manage independent tasks and optionally link them to intents or focus sessions.

### Commands

*   **Add a task:**
    ```bash
    mindloop task add "Write documentation"
    # Link to an intent or focus session
    mindloop task add "Write API docs" --intent-id 12
    ```

*   **List tasks:**
    ```bash
    mindloop task list
    ```

*   **Complete a task:**
    ```bash
    mindloop task complete <task_id>
    ```

---

## Routines

Group multiple habits into a specific time of day.

### Commands

*   **Create a routine:**
    ```bash
    mindloop routine create "Morning Routine" --time "Morning"
    ```

*   **Add a habit to a routine:**
    ```bash
    mindloop routine add-habit <routine_id> <habit_id>
    ```

*   **List routines:**
    ```bash
    mindloop routine list
    ```

---

## Authentication

Manage local access to the Mindloop Web UI.

### Commands

*   **Setup or change password:**
    ```bash
    mindloop auth setup
    ```

*   **Disable password protection:**
    ```bash
    mindloop auth disable
    ```

---

## Summary

Get a bird's-eye view of your productivity metrics.

### Commands

*   **Generate summary:**
    ```bash
    mindloop summary
    ```
    By default, shows the summary for today.

*   **Time ranges:**
    ```bash
    mindloop summary --day    # Today (default)
    mindloop summary --week   # Current week
    mindloop summary --month  # Current month
    mindloop summary --year   # Current year
    ```

---

## Data Management (Backup & Restore)

Manage your local data with built-in export and import tools.

### Commands

*   **Backup data:**
    ```bash
    mindloop backup my_backup.json
    ```
    Exports all your data (Intents, Habits, Journal, etc.) to a JSON file.

*   **Restore data:**
    ```bash
    mindloop restore my_backup.json
    ```
    **WARNING**: This replaces your current database with the backup data.

---

## Help

To see the list of available commands and flags at any time:

```bash
mindloop help
# or for specific commands
mindloop habit help
```
