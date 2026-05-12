# Instruction Restructure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Consolidate and update project instructions into a new `AGENT.md` file and setup symlinks for `GEMINI.md` and `CLAUDE.md`.

**Architecture:** Create a single source of truth (`AGENT.md`) for AI agents and symlink platform-specific instruction files to it. This ensures all agents have the latest project context and guidance.

**Tech Stack:** Markdown, Shell (Symlinks)

---

### Task 1: Create AGENT.md

**Files:**
- Create: `AGENT.md`

- [ ] **Step 1: Write the content for AGENT.md**

```markdown
# Mindloop Project Context (AGENT.md)

## Project Overview
**Mindloop** is a comprehensive productivity suite designed for local-first workflow management. It provides tools for tracking intents, focus sessions, habits, and journals.

The project operates as a dual-interface application:
1.  **CLI:** A command-line interface for low-latency interaction.
2.  **Web Server/UI:** A local web server providing a visual interface and REST API.

**Key Technologies:**
*   **Language:** Go (Golang) 1.26
*   **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
*   **ORM:** [GORM](https://gorm.io/)
*   **Database:** SQLite (Default) or PostgreSQL (BYODB mode)
*   **Web:** Go `net/http` with HTML templates (Server-Side Rendered)
*   **Logging:** Zerolog

## Architecture
The project follows a clean architecture pattern:
*   `cmd/`: Application entry points.
    *   `cmd/cli/`: CLI command definitions (Cobra).
    *   `cmd/server/`: Web server entry point.
*   `internal/core/`: Business logic domain (Focus, Habit, Intent, Journal, Summary, Quest, Note, Routine).
*   `api/v1/`: HTTP handlers for the web server.
*   `db/`: Database connection and schema management.
*   `web/`: Static assets and HTML templates.

## Building and Running

The project uses a `Makefile` for build automation.

### Build
*   **Build All:** `make build` (Generates the single `mindloop` binary)

### Run
*   **CLI:**
    *   Run the binary directly: `./mindloop <command>` (e.g., `./mindloop help`)
*   **Server:**
    *   Run via subcommand: `./mindloop server`
    *   Run locally (foreground): `make run-server` (Default port: 8765)
    *   Start in background: `make start-server`
    *   Stop background server: `make kill-server`

### Testing & Verification
*   **Run Unit Tests:** `make test`
*   **Linting:** `make lint`
*   **Formatting:** `make fmt`

## Configuration & Persistence
*   **Storage:** Data is stored by default in `~/.mindloop/mindloop_local.db`.
*   **Overrides:** Checks for a local `mindloop_local.db` in the current directory.
*   **BYODB Mode:** Can be configured to use an external PostgreSQL database via `mindloop configure`.

## Agent Guidance
*   **Business Logic:** Look in `internal/core/` for domain rules.
*   **API Handlers:** Look in `api/v1/` for web interface logic.
*   **CLI Commands:** Check `cmd/cli/` and `docs/CLI_USAGE.md`.
*   **Style:** Always run `make fmt` after modifying Go code.
```

- [ ] **Step 2: Commit AGENT.md**

Run: `git add AGENT.md && git commit -m "docs: create AGENT.md as primary instruction file"`

### Task 2: Setup Symlinks

**Files:**
- Modify: `GEMINI.md` (Delete and recreate as symlink)
- Create: `CLAUDE.md` (Symlink)

- [ ] **Step 1: Delete existing GEMINI.md**

Run: `rm GEMINI.md`

- [ ] **Step 2: Create symlinks**

Run: `ln -s AGENT.md GEMINI.md && ln -s AGENT.md CLAUDE.md`

- [ ] **Step 3: Verify symlinks**

Run: `ls -la GEMINI.md CLAUDE.md`
Expected: Output showing `GEMINI.md -> AGENT.md` and `CLAUDE.md -> AGENT.md`.

- [ ] **Step 4: Commit symlinks**

Run: `git add GEMINI.md CLAUDE.md && git commit -m "docs: symlink GEMINI.md and CLAUDE.md to AGENT.md"`

### Task 3: Final Validation

- [ ] **Step 1: Verify AGENT.md content**

Run: `cat AGENT.md`
Expected: Full content of the new AGENT.md.

- [ ] **Step 2: Verify GEMINI.md content (via link)**

Run: `cat GEMINI.md`
Expected: Full content of the new AGENT.md.
