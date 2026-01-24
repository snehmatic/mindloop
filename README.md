# Mindloop

_P.S. A constant WIP_

**Mindloop** is a comprehensive productivity suite designed for local-first workflow management. It operates as a CLI tool, a local API, and a UI, utilizing a local SQLite database with BYODB (Bring Your Own Database) support.

## Features

* **Intents:** Track single-threaded work items to maintain focus.
* **Focus Sessions:** Break intents into deep work chunks, deep work.
* **Habits:** Daily and weekly habit tracking.
* **Journal:** Reflection channel. Use it as you like, end-of-day closure or gratitude.
* **Summary:** High-level metrics and bird's-eye view of productivity data. With support for filtering by date range.

## Architecture

* **Core:** Go backend with local SQLite (or BYODB).
* **Interfaces:** CLI, Local REST API, Web UI (vibe coded with Gemini, not a frontend fanatic here).

---

## Motivation

Mindloop is a productivity suite, for getting started with intentions, habits, journals or focus sessions. 
As a developer with attention problems, Mindloop started as a personal CLI tool to manage my daily work routine. `Intents` to set one single intention or work item to track and focus on, and `focus sessions` to break that intention down into chuncks of high quality deep work frames. 
This workflow worked well for starters. Alongside, the additional `habits` and `journal` features were added just because. But weirdly enough those caught up and I incorporated them into my workflow as well. I'd start my day ticking off habits that I had set (daily and weekly) and by the end of the day, write down a mini journal to just close the loop, end work and go touch some grass. Funny story, the last part never happened. Anywho, now that I had so much data on my productivity, it only made sense to have a summary to have a birds eye point of view over all this.
