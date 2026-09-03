# Mindloop Architecture

Mindloop follows a clean architecture pattern, separating core business logic from interfaces and data storage.

## Table of Contents
- [Core Services (Business Logic)](#core-services-business-logic)
- [Interfaces (Presentation Layer)](#interfaces-presentation-layer)
- [Data Layer](#data-layer)

## Core Services (Business Logic)
Located in `internal/core`, this layer isolates the rules for each domain:
*   `intent`: Manages single-threaded work goals and lifecycles.
*   `quest`: Handles ad-hoc "side quests" that shadow main intents.
*   `focus`: Handles deep work session timers and tracking.
*   `habit`: Logic for daily/weekly habit tracking and streaks.
*   `journal`: Manages daily reflections and mood logging.
*   `summary`: Aggregates data for productivity reporting.

## Interfaces (Presentation Layer)
*   **CLI (`cmd/cli`):** Built with [Cobra](https://github.com/spf13/cobra), this interface interacts directly with the local database for low-latency command-line usage. Use `mindloop --version` to check your current version. Check out the [CLI Usage Guide](CLI_USAGE.md) for detailed command instructions.
*   **Web Server (`cmd/cli/server.go`):** A Go HTTP server exposed through the `mindloop server` subcommand.
*   **Web UI:** Server-Side Rendered (SSR) HTML templates (`web/templates`) utilizing vanilla CSS/JS. "Vibe coded" with Gemini (backend-focused developer approach).

## Data Layer
*   Uses [GORM](https://gorm.io/) for ORM capabilities.
*   **Default:** Zero-config SQLite (`mindloop_local.db`).
*   **Universal Storage:** By default, data is stored in `~/.mindloop/` to ensure persistence across different working directories. It also checks for a local `mindloop_local.db` in the current directory for project-specific overrides.
*   **Optional:** Supports PostgreSQL via configuration (BYODB mode).
