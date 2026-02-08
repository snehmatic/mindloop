# Mindloop

![CI Status](https://github.com/snehmatic/mindloop/actions/workflows/ci.yml/badge.svg)
![Release Version](https://img.shields.io/github/v/release/snehmatic/mindloop)
![License](https://img.shields.io/github/license/snehmatic/mindloop)
![Go Version](https://img.shields.io/github/go-mod/go-version/snehmatic/mindloop)

**Mindloop** is a comprehensive productivity suite designed for local-first workflow management. It operates as a **CLI tool**, a **local API**, and a **UI**, utilizing a local SQLite database with BYODB (Bring Your Own Database) support.

## Features

### 🎯 Intents
Set a single-threaded goal to maintain absolute focus. **One thing at a time.**
- Track lifecycle of your intents (Active, Completed, Failed).
- Prevent context switching by committing to a single objective.

### ⏱️ Focus Sessions
Break your intents into deep work chunks with a built-in timer.
- Track duration and frequency of deep work.
- Associate sessions directly with your active Intent.

### ✅ Habits
Track daily and weekly habits to build consistency.
- Simple check-in system.
- Track streaks and completion rates.

### 📓 Journal
A dedicated space for reflections, end-of-day closure, or gratitude logging.
- Markdown support for rich text entries.
- "Close the loop" at the end of your workday.

### 📊 Summary
High-level metrics and a bird's-eye view of your productivity data.
- Visualize time spent in focus.
- Review completed intents and habit consistency.
- Filter data by date ranges.

## Architecture

Mindloop follows a clean architecture pattern, separating core business logic from interfaces and data storage.

*   **Core Services (Business Logic):**
    Located in `internal/core`, this layer isolates the rules for each domain:
    *   `intent`: Manages single-threaded work goals and lifecycles.
    *   `focus`: Handles deep work session timers and tracking.
    *   `habit`: Logic for daily/weekly habit tracking and streaks.
    *   `journal`: Manages daily reflections and mood logging.
    *   `summary`: Aggregates data for productivity reporting.

*   **Interfaces (Presentation Layer):**
    *   **CLI (`cmd/cli`):** Built with [Cobra](https://github.com/spf13/cobra), this interface interacts directly with the local database for low-latency command-line usage. Check out the [CLI Usage Guide](docs/CLI_USAGE.md) for detailed command instructions.
    *   **Web Server (`cmd/server`):** A Go HTTP server exposing a REST API (`api/v1`).
    *   **Web UI:** Server-Side Rendered (SSR) HTML templates (`web/templates`) utilizing vanilla CSS/JS. "Vibe coded" with Gemini (backend-focused developer approach).

*   **Data Layer:**
    *   Uses [GORM](https://gorm.io/) for ORM capabilities.
    *   **Default:** Zero-config SQLite (`mindloop_local.db`) for a true local-first experience.
    *   **Optional:** Supports PostgreSQL via configuration (BYODB mode).

---

## Getting Started

### Prerequisites

*   [Go](https://go.dev/) 1.24+
*   [Make](https://www.gnu.org/software/make/)

### Installation

#### Option 1: Homebrew (macOS/Linux)
The easiest way to install and keep Mindloop updated.
```bash
brew tap snehmatic/mindloop
brew install mindloop
```

#### Option 2: Go Install
If you have Go installed, you can install the latest version directly:
```bash
go install github.com/snehmatic/mindloop@latest
```

#### Option 3: Download Binaries
Download the latest pre-compiled binary for your OS from the [Releases Page](https://github.com/snehmatic/mindloop/releases).

#### Option 4: Build from Source
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/snehmatic/mindloop.git
    cd mindloop
    ```

2.  **Build:**
    ```bash
    make build
    ```
    This generates `mindloop` (CLI) and `mindloop-server` (Server) binaries.

3.  **Run the Server (Quick Start):**
    For local development with the UI:
    ```bash
    make run-server
    ```
    This will start the server on port `8765` (default).

    To use a custom port or mode:
    ```bash
    make run-server PORT=9000 MODE=byodb
    ```
    
    Access the UI at [http://localhost:8765](http://localhost:8765)

### Build Commands

The project includes a `Makefile` to simplify common tasks:

*   `make build`: Build both the CLI and Server binaries.
*   `make build-cli`: Build only the CLI binary.
*   `make start-server`: Build and run the server in the background (daemon).
*   `make kill-server`: Stop the background server.
*   `make test`: Run unit tests.
*   `make clean`: Remove build artifacts.

## Configuration

By default, Mindloop runs in **Local Mode** using SQLite.

To use an external database (e.g., PostgreSQL), you can configure the application via environment variables.

1.  Copy the example env file:
    ```bash
    cp example.env .env
    ```
2.  Edit `.env` with your database credentials (DB_HOST, DB_PORT, etc.).

## Motivation

Mindloop is a productivity suite, for getting started with intentions, habits, journals or focus sessions.

As a developer with attention problems, Mindloop started as a personal CLI tool to manage my daily work routine. `Intents` to set one single intention or work item to track and focus on, and `focus sessions` to break that intention down into chunks of high quality deep work frames.

This workflow worked well for starters. Alongside, the additional `habits` and `journal` features were added just because. But weirdly enough those caught up and I incorporated them into my workflow as well. I'd start my day ticking off habits that I had set (daily and weekly) and by the end of the day, write down a mini journal to just close the loop, end work and go touch some grass.