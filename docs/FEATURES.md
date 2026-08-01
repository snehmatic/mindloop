# Mindloop Features

<p align="center">
  <img src="../web/static/images/logo-readme.svg" width="160" alt="Mindloop Logo">
</p>

## Table of Contents
- [Web & TUI Dashboard](#web--tui-dashboard)
- [Intents](#intents)
- [Focus Sessions](#focus-sessions)
- [Tasks & Sub-tasks](#tasks--sub-tasks)
- [Habits](#habits)
- [Write (Journal & Notes)](#write-journal--notes)
- [Command Palette & Navigation](#command-palette--navigation)
- [AI-Assisted Journal Generation](#ai-assisted-journal-generation)
- [The Void](#the-void)
- [Summary](#summary)
- [Gamification & Points](#gamification--points)
- [Settings & Configuration](#settings--configuration)

## Web & TUI Dashboard
The "Anti-Dashboard" vision puts friction-less capture first.
- **Interactive Web UI Dashboard**: A true command center. Start focus sessions, tick off pending tasks, complete daily habits, and quick capture notes directly from the home grid without navigating away.
- **TUI Dashboard (`mindloop dash`)**: A beautiful terminal interface powered by Bubbletea that features a live-updating focus session timer and interactive habit toggling.

## Intents
Set a single-threaded goal to maintain absolute focus. **One thing at a time.**
- **Natural Language Dates**: Type "Finish report by next Friday" to automatically set a due date.
- Track lifecycle of your intents (Active, Completed, Failed).
- **Side Quest Modal**: Manage ad-hoc tasks that interrupt your flow without losing context.

## Focus Sessions
Break your intents into deep work chunks with a built-in timer.
- Track duration and frequency of deep work.
- Associate sessions directly with your active Intent.

## Tasks & Sub-tasks
Manage independent to-do items and link them to your bigger goals.
- **Natural Language Dates**: Just like Intents, tasks support natural language date parsing (e.g., "Review PR tomorrow").
- **Full Web UI**: Create, complete, and delete tasks directly from the Tasks page.
- **Sub-tasks**: Break tasks into smaller, actionable sub-tasks with independent completion tracking.
- **Link to Intents & Focus Sessions**: Optionally attach tasks to an active Intent or Focus Session when creating them.
- **Inline Task Creation**: Add tasks directly from within the Intent or Focus Session views via a toggleable form.
- **Drag & Drop Reordering**: Reorder tasks and sub-tasks with persistent drag-and-drop.
- **Status Filtering**: Filter the task list by All, Pending, or Completed status.
- **Visual Polish**: Completed tasks are dimmed with a subtle green tint; pill-shaped action buttons with hover effects.

## Habits
Track daily and weekly habits to build consistency.
- **Activity Heatmap**: Visualize your consistency with a GitHub-style activity grid.
- Simple check-in system with streak tracking.

## Write (Journal & Notes)
A dedicated space for your thoughts.
- **Journal**: A space for reflections, end-of-day closure, or gratitude logging. Features mood tracking and rich markdown support. "Close the loop" at the end of your workday.
- **Notes**: Quickly capture thoughts, ideas, or longer-form reference material that doesn't fit into a daily journal entry. Use the "Quick Note" widget on the Dashboard to capture instantly.

## Command Palette & Navigation
Frictionless navigation across the entire Web UI.
- **Omnibar (`Cmd+K`)**: Hit `Cmd+K` anywhere in the Web UI to instantly open a spotlight-style command palette to jump between intents, tasks, habits, or the void.

## AI-Assisted Journal Generation
Generate a reflective journal entry from your activity data.
- **CLI support**: Use `mindloop journal generate` with daily, weekly, or yearly ranges.
- **Web UI support**: Use the Journal auto-generate flow, or the Summary page's AI Overview shortcut.
- **Summary-backed**: The generated entry is based on Mindloop summary data such as intents, focus sessions, habits, and completed work.
- **Optional save flow**: Review the generated text before choosing whether to save it into the journal.
- **Configurable AI settings**: Manage model and prompt settings from the Settings page.

## The Void
A specialized feature to write down thoughts, worries, or frustrations and let them disappear. Perfect for clearing your mind before deep work.
- **Web UI**: A place to release thoughts without keeping them as permanent notes.
- **CLI support**: `mindloop void [minutes]` starts a timed breathing session. The default duration is 5 minutes.

## Summary
High-level metrics and a bird's-eye view of your productivity data.
- Visualize time spent in focus.
- Review completed intents and habit consistency.
- Filter data by date ranges.
- Use summaries as input for AI-assisted journal generation.

## Gamification & Points
Turn your productivity into a game with a built-in reward system.
- **Earn Points**: Get rewarded for completing focus sessions, habits, intents, journals, tasks, sub-tasks, and side quests.
- **Milestones**: Reach configurable point milestones (200 pts by default) to trigger special full-screen celebration screens.
- **Customizable Rewards**: Define your own point values and milestone interval in the Settings.
- **Visual Progress**: Track your points over time with a dedicated chart in the Summary report.
- **Celebrations**: Enjoy confetti animations when you finish tasks and reach new heights.


## Lifecycle Event Hooks (Extensibility)
Extend Mindloop by running custom scripts when specific events occur.
- **Hook Scripts**: Place executable scripts in `~/.mindloop/hooks/`.
- **Supported Events**:
  - `focus_start`: Triggered when a new focus session starts. Passes `MINDLOOP_FOCUS_TITLE`.
  - `focus_stop`: Triggered when a focus session ends. Passes `MINDLOOP_FOCUS_TITLE` and `MINDLOOP_FOCUS_DURATION`.
- **Environment**: Scripts are passed context variables in their environment.
- **Async Execution**: Hooks run asynchronously so they do not block the app.

## Settings & Configuration
- **Stacked Layout**: Clean, organized settings page with vertical stacking for better focus. Includes segmented controls and toggle switches.
- **Data Management**: Easy backup, restore, and reset options.
- **BYODB**: Bring Your Own Database (PostgreSQL) support.

> **Note:** For a visual tour of the interface and responsive design, check out the [Web UI Documentation](web_ui.md).
