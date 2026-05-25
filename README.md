# Mindloop

<p align="center">
  <img src="web/static/images/logo-readme.svg" width="200" alt="Mindloop Logo">
</p>

![CI Status](https://github.com/snehmatic/mindloop/actions/workflows/ci.yml/badge.svg)
![Release Version](https://img.shields.io/github/v/release/snehmatic/mindloop)
![License](https://img.shields.io/github/license/snehmatic/mindloop)
![Go Version](https://img.shields.io/github/go-mod/go-version/snehmatic/mindloop)

**Mindloop** is a comprehensive productivity suite designed for local-first workflow management. It operates as a **CLI tool**, a **local API**, and a **UI**, utilizing a local SQLite database with BYODB (Bring Your Own Database) support.

## Table of Contents
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Build Commands](#build-commands)
- [Devcontainers](#devcontainers)
- [Documentation](#documentation)
- [Configuration](#configuration)
- [Motivation](#motivation)

## Getting Started

### Prerequisites

*   [Go](https://go.dev/) 1.26+
*   [Make](https://www.gnu.org/software/make/)

### Installation

#### Option 1: Homebrew (macOS/Linux)
The easiest way to install and keep Mindloop updated.
```bash
brew tap snehmatic/mindloop
brew install mindloop
```

> **Note:** If you are not seeing the latest version, run `brew update` to refresh the tap.

**Run as a background service:**
You can use Homebrew Services to run the Mindloop server in the background:
```bash
brew services start mindloop
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

    **Live Reloading (Development):**
    For active development, you can run the server with live reloading. This requires [Air](https://github.com/air-verse/air).
    First, install `air` globally:
    ```bash
    go install github.com/air-verse/air@latest
    ```
    Then, start the development server:
    ```bash
    make dev
    ```
    The server will automatically rebuild and restart when you save changes to Go, HTML, CSS, or JS files.

### Build Commands

The project includes a `Makefile` to simplify common tasks:

*   `make build`: Build both the CLI and Server binaries.
*   `make build-cli`: Build only the CLI binary.
*   `make dev`: Run the server with live-reloading (requires `air`).
*   `make start-server`: Build and run the server in the background (daemon).
*   `make kill-server`: Stop the background server.
*   `make test`: Run unit tests.
*   `make clean`: Remove build artifacts.

## Devcontainers

This project includes a [Devcontainer](https://containers.dev/) configuration for a seamless development experience using [GitHub Codespaces](https://github.com/features/codespaces), [Visual Studio Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers), or any other compatible devcontainer tool.

**Development with Air:**

Development uses [Air](https://github.com/air-verse/air) for live-reloading. Air is automatically installed in the devcontainer on setup and runs as the `postStartCommand`. When you make changes to Go, HTML, CSS, or JS files, the server automatically rebuilds and restarts — no manual rebuild step required.

**Configuration Gotcha:**

> After the first devcontainer run, the Mindloop configuration file (and the Air binary) lives in the devcontainer environment's user home folder (e.g., `/home/vscode`), which is **not synced** with your local machine. You will need to manually run the build and any configuration commands each time you open the devcontainer **fresh**.

## Documentation

* **[Features](docs/FEATURES.md)**: Explore the different features Mindloop offers like Intents, Focus Sessions, Habits, and more.
* **[Architecture](docs/ARCHITECTURE.md)**: Learn about Mindloop's clean architecture, interfaces, and data layer.
* **[CLI Usage Guide](docs/CLI_USAGE.md)**: Detailed instructions on how to use Mindloop directly from the command line.
* **[Web UI Documentation](docs/web_ui.md)**: A visual tour of the web interface and responsive design.
* **[Releasing Guide](docs/RELEASES.md)**: Instructions on how to create and publish new releases.

## Configuration

By default, Mindloop runs in **Local Mode** using SQLite.

To use an external database (e.g., PostgreSQL), you can configure the application via environment variables.

1.  Copy the example env file:
    ```bash
    cp example.env .env
    ```
2.  Edit `.env` with your database credentials (DB_HOST, DB_PORT, etc.).

## Motivation

**Mindloop** is a productivity suite for getting started with intentions, habits, journals, and focus sessions.
As a developer with attention problems, Mindloop originally started as a very personal CLI tool to manage my daily work routine and reduce some of the mental clutter that comes with context switching all day.
The initial workflow was intentionally simple:
* Intents to define one single thing I wanted to work on or pay attention to
* Focus sessions to break that intention down into smaller chunks of high quality, intentional deep work

That alone worked surprisingly well.
Instead of planning an ideal day or maintaining some giant productivity system that I'd eventually ignore, this gave me something much smaller and more immediate: decide what matters right now, work on it, and close the loop.
Along the way, additional features like `habits` and `journal` were added almost as an afterthought. But weirdly enough, those ended up becoming part of the workflow too.
I’d usually start my day by ticking off a few habits I cared about (daily and weekly), then work through an active intention using focus sessions, and by the end of the day, write a tiny journal entry just to mentally close work for the day, shut the laptop, and go touch some grass (or did I?).

Over time, **Mindloop** evolved beyond just a CLI experiment.
What started as a tool to help me focus on work gradually became a slightly more opinionated system for managing attention, reflection, and intentionality. Tasks, subtasks, notes, summaries, side quests, AI overviews, and other features naturally grew around that core.
But the philosophy has stayed mostly the same.
_Mindloop is not meant to be a system for planning the perfect day, optimizing every hour, or building the most aesthetic productivity dashboard known to mankind (I despise frontend).
It is intentionally a bit more present and reflective._

The core idea is simple:
* pick something meaningful to focus on
* work on it in deliberate chunks
* keep track of a few recurring systems that matter
* reflect a little
* and finally, move on

This is also why some seemingly obvious productivity features have intentionally been avoided or delayed.
For example, routines/rituals or heavily pre-planned workflows sound useful on paper (and honestly, I still think they are), but they introduce a strange overlap with the intention-centric design of Mindloop.
If everything is pre-planned, templated, stacked, drag-and-dropped, and optimized into perfect routines, then the question becomes: what purpose do intentions serve anymore?
That starts pulling Mindloop in a different direction.

At least for now, Mindloop is less about designing a perfect life system, and more about helping reduce mental chaos in the present.
It exists to answer simpler questions:
* What am I doing right now?
* Am I actually paying attention to it?
* Did I close the day with some amount of intention?
That’s really where Mindloop started, and honestly, still the _soul_ of the project.
