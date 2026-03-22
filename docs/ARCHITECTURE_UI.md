# UI Architecture & Modernization

This document outlines the architectural principles and technical details of the Mindloop web interface, specifically following the "Frontend Modernization" (Issue #49).

## Table of Contents
- [Core Philosophy: The "Snappy Monolith"](#core-philosophy-the-snappy-monolith)
- [HTMX Integration](#htmx-integration)
  - [Partial Templates](#partial-templates)
  - [Backend Handling](#backend-handling)
- [CSS Design System](#css-design-system)
  - [Key Variables](#key-variables)
  - [Dark Mode](#dark-mode)
- [Iconography](#iconography)
- [Module Specifics](#module-specifics)
  - [Focus Module](#focus-module)
  - [Habits Module](#habits-module)

## Core Philosophy: The "Snappy Monolith"

Mindloop adheres to a **Server-Side Rendered (SSR)** architecture enhanced with **HTMX**. This allows us to maintain a single binary deployment (Go) while providing a modern, Single-Page-Application (SPA) feel.

*   **Language:** Go (Golang) with `html/template`
*   **Interactivity:** [HTMX](https://htmx.org/)
*   **Styling:** Native CSS (Utility + Component classes)
*   **Icons:** [Lucide Icons](https://lucide.dev/)

## HTMX Integration

We use HTMX to perform partial page updates, avoiding full reloads for common actions like checking off a habit or stopping a timer.

### Partial Templates
To support HTMX, standard templates are broken down into smaller, reusable "fragments" or "partials".

*   **Convention:** Partials that represent a specific UI component (like a card) are prefixed with `_` (e.g., `_habit_card.html`). Partials that represent a section of a page are named descriptively (e.g., `focus_active_timer.html`).
*   **Location:** All templates reside in `web/templates/`.

### Backend Handling
The Go backend (`api/v1/handlers_impl.go`) detects HTMX requests via the `HX-Request` header.

*   **Standard Request:** Renders the full page (Layout + Content).
*   **HTMX Request:** Renders **only** the relevant partial(s) and returns them to the client.
*   **Error Handling:** If an error occurs during an HTMX request, the backend sets the `HX-Redirect` header to trigger a full page reload/redirect to the error page/state, preventing nested layout issues.

## CSS Design System

Styles are consolidated in `web/static/css/style.css`. We use CSS Variables for theming and consistency.

### Key Variables
*   **Colors:** `--primary`, `--primary-light`, `--bg-body`, `--text-main`
*   **Shadows:** `--shadow-sm`, `--shadow-md`, `--shadow-lg`, `--shadow-glow` (Layered depth system)
*   **Gradients:** `--gradient-primary`

### Dark Mode
Dark mode is implemented via the `html.dark-mode` class, which overrides the CSS variables. This preference is persisted in `localStorage`.

## Iconography

We use **Lucide Icons** for a consistent, professional visual language.
*   **Implementation:** `<i data-lucide="icon-name"></i>`
*   **Loading:** The `lucide.createIcons()` function is called on page load and re-triggered on `htmx:afterSwap` events to ensure icons render correctly in dynamically loaded content.

## Module Specifics

### Focus Module
*   **Files:** `focus.html`, `focus_active_timer.html`, `focus_session_list.html`
*   **Behavior:** Starting/Stopping a timer updates the active timer view and the session list independently without reloading the page.

### Habits Module
*   **Files:** `habits.html`, `_habit_card.html`
*   **Behavior:** "Check-in", "Undo", and "Delete" actions swap only the specific habit card (`outerHTML` swap), preserving the user's scroll position and context.
