# Design Spec: Instruction Restructure (AGENT.md)

**Date:** 2026-05-05
**Topic:** Consolidating and updating project instructions for AI agents.

## Overview
Mindloop's project instructions are currently housed in `GEMINI.md`. To provide a consistent experience across different AI agents (Gemini, Claude, etc.), we will transition to an `AGENT.md` "source of truth" and symlink other agent-specific files to it.

## Goals
1. Create a comprehensive `AGENT.md` file that reflects the current state of the codebase.
2. Establish `AGENT.md` as the primary instruction file.
3. Symlink `GEMINI.md` and `CLAUDE.md` to `AGENT.md`.
4. Update instructions to reflect architectural changes (Unified binary, persistence conventions, new features).

## Proposed Content for `AGENT.md`

### 1. Project Identity
- **Name:** Mindloop
- **Vision:** Local-first productivity suite for intents, focus, habits, and journals.
- **Interfaces:** CLI (`mindloop`) and Web Server/UI (`mindloop server`).

### 2. Tech Stack
- **Language:** Go 1.26
- **CLI:** Cobra
- **Database:** GORM + SQLite (default) / PostgreSQL (BYODB)
- **Web:** Gorilla Mux + SSR (Go Templates)
- **Logging:** Zerolog

### 3. Key Architectures & Conventions
- **Clean Architecture:** 
  - `cmd/`: entry points.
  - `internal/core/`: domain logic (intent, quest, focus, habit, journal, summary, note).
  - `api/v1/`: HTTP handlers.
  - `db/`: GORM models and DB init.
  - `web/`: Assets and templates.
- **Persistence:** Default to `~/.mindloop/mindloop_local.db`. Support local directory override.
- **Vibe Coding:** Frontend is lean, SSR-based, and vanilla.

### 4. Development Workflow
- **Bootstrap:** `mindloop configure`
- **Build:** `make build` (Single binary `mindloop`)
- **Run Server:** `make run-server` or `./mindloop server`
- **Verification:** `make test`, `make lint`, `make fmt`

### 5. Agent-Specific Guidance
- **Symbol Discovery:** Use `internal/core` for business rules.
- **Command Discovery:** Check `cmd/cli` and `docs/CLI_USAGE.md`.
- **Modifications:** Always run `make fmt` after Go changes. Verify with `make build` and `make test`.

## Implementation Plan

1. **Scaffold AGENT.md:** Write the consolidated content.
2. **Setup Symlinks:**
   - Delete existing `GEMINI.md`.
   - Create symlink `GEMINI.md` -> `AGENT.md`.
   - Create symlink `CLAUDE.md` -> `AGENT.md`.
3. **Validation:** Ensure the symlinks are functional and the content is accurate.

## Risks & Mitigations
- **Link Breakage:** Symlinks are standard on Unix/macOS; if the project is moved to Windows without Git symlink support, they might become text files. *Mitigation:* Document the structure in `README.md` if necessary, but standard for these agents.
- **Duplication:** Content in `AGENT.md` might overlap with `README.md`. *Mitigation:* `AGENT.md` is specifically tuned for agent context (e.g., "Look in X for Y"), whereas `README.md` is for human users.
