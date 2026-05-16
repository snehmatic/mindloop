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
*   `internal/core/`: Business logic domain (Focus, Habit, Intent, Journal, Summary, Quest, Note).
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
