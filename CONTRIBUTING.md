# Contributing to Mindloop

Thanks for helping improve Mindloop. This project is local-first and intentionally small enough to reason about, so clear issues and focused pull requests matter.

## Before You Start

- Search existing issues and pull requests for related work.
- For user-facing changes, identify whether the change affects the CLI, Web UI, API, database/storage, docs, or release flow.
- If a feature is experimental or partially implemented, describe it that way. Routines/Rituals are currently CLI-only and may change.
- If you are using an AI coding agent or running automated PR sweeps, provide it with `AGENTS.md` as project context before making changes.

## Development

- Use `make build` to build the project.
- Use `make test` for the Go test suite.
- Use `make lint` when `golangci-lint` is installed.
- Use `make fmt` after changing Go code.
- Keep changes focused. Avoid unrelated refactors in feature or docs pull requests.

## Documentation

Update documentation whenever behavior changes:

- `README.md` for top-level positioning, install, and important feature highlights.
- `docs/FEATURES.md` for user-facing feature descriptions.
- `docs/CLI_USAGE.md` for CLI commands and flags.
- `docs/web_ui.md` for Web UI behavior and screenshots.
- `docs/ROADMAP.md` for planned, experimental, or recently completed work.
- `AGENTS.md` for guidance that helps AI coding agents and maintainers navigate the project.

## Pull Requests

Please include:

- A short summary of what changed.
- Testing or verification performed.
- Documentation updates, if applicable.
- Screenshots or terminal output for UI/CLI behavior changes.
- Any known limitations or follow-up work.

For AI-assisted contributions, mention which areas were inspected and which commands were run. Do not claim checks were run if they were skipped.
