# Issue 1: [Web UI] Command Palette (Omnibar) for Frictionless Navigation & Entry

**Labels:** `enhancement`, `ui`, `ux`

### Rationale & Market Fit
Power users and developers expect keyboard-first navigation (popularized by Superhuman, Linear, Raycast). Mindloop's web UI currently requires mouse clicks to create tasks, switch pages, or start sessions. An omnibar will reduce action latency to near-zero.

### Expected Outcomes
- Users can press `Cmd+K` (macOS) or `Ctrl+K` (Windows/Linux) from anywhere in the Web UI to open a modal Command Palette.
- The palette features a sleek, glassmorphic design matching the current aesthetics.
- Users can search for actions and execute them instantly.

### Implementation Details
1. **Frontend (`web/static/js/main.js` or similar):**
   - Add a global `keydown` event listener for `Cmd+K` / `Ctrl+K`.
   - Build a modal overlay for the palette.
2. **Supported Commands:**
   - Navigation: `Go to Summary`, `Go to Tasks`, `Go to Settings`, `Go to Journal`
   - Actions: `Start Focus Session`, `Create Intent`, `Create Task`
3. **UI/UX:**
   - Use an input field that auto-focuses.
   - List matching commands below it. Keyboard arrow navigation and `Enter` to select.
   - Follow the `ARCHITECTURE_UI.md` guidelines for styling (CSS variables, clean layout).
