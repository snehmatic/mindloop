# Mindloop Release Guide

This document outlines the process for releasing a new version of Mindloop. Mindloop uses [GoReleaser](https://goreleaser.com/) combined with GitHub Actions to automate the build and release process.

## Table of Contents
- [Triggering a Release](#triggering-a-release)
  - [Steps to Release](#steps-to-release)
- [What Happens Next?](#what-happens-next)
- [Testing Releases Locally (Snapshot)](#testing-releases-locally-snapshot)
- [Troubleshooting](#troubleshooting)

## Triggering a Release

Releases are triggered automatically when a new Git tag starting with `v` is pushed to the repository.

### Steps to Release

1. **Ensure your local `main` branch is up to date:**
   ```bash
   git checkout main
   git pull origin main
   ```

2. **Create a new annotated tag for the release:**
   Use semantic versioning for your tags (e.g., `v1.0.0`, `v1.0.1`, `v1.1.0`).
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   ```

3. **Push the tag to GitHub:**
   ```bash
   git push origin v1.0.0
   ```

## What Happens Next?

Once the tag is pushed to GitHub, the following automated steps occur:

1. The `.github/workflows/release.yml` GitHub Action is triggered.
2. The Action checks out the code and sets up the Go environment.
3. **GoReleaser** runs based on the configuration in `.goreleaser.yaml`.
   - It builds binaries for Linux, macOS (Darwin), and Windows across multiple architectures (`amd64`, `arm64`).
   - It creates compressed archives (`tar.gz` for most, `zip` for Windows).
   - It updates the Homebrew tap repository (`snehmatic/homebrew-mindloop`).
   - It publishes a new GitHub Release with the compiled binaries and an auto-generated changelog.

## Testing Releases Locally (Snapshot)

If you want to test the release build process locally without actually publishing anything to GitHub or updating Homebrew, you can run GoReleaser in snapshot mode:

```bash
goreleaser release --snapshot --clean
```
This will build the binaries and archives in the `dist/` directory without publishing them.

## Troubleshooting

If a release fails, you can check the logs of the `Release` GitHub Action workflow on the [Actions page](https://github.com/snehmatic/mindloop/actions) of the repository.
