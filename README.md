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

Mindloop is a productivity suite, for getting started with intentions, habits, journals or focus sessions.

As a developer with attention problems, Mindloop started as a personal CLI tool to manage my daily work routine. `Intents` to set one single intention or work item to track and focus on, and `focus sessions` to break that intention down into chunks of high quality deep work frames.

This workflow worked well for starters. Alongside, the additional `habits` and `journal` features were added just because. But weirdly enough those caught up and I incorporated them into my workflow as well. I'd start my day ticking off habits that I had set (daily and weekly) and by the end of the day, write down a mini journal to just close the loop, end work and go touch some grass.