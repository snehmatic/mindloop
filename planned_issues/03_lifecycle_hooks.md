# Issue 3: [Backend] Lifecycle Event Hooks (Extensibility)

**Labels:** `enhancement`, `backend`, `automation`

### Rationale & Market Fit
As a local-first application, Mindloop has a superpower: it can talk to your operating system. When a user starts a Focus Session, they shouldn't have to manually turn on macOS "Do Not Disturb" or block Reddit. Mindloop should trigger local shell scripts on lifecycle events.

### Expected Outcomes
- Users can drop executable scripts into `~/.mindloop/hooks/`.
- Scripts named `on-focus-start`, `on-focus-end`, `on-intent-complete` are automatically executed by Mindloop when those events occur.

### Implementation Details
1. **Core Logic (`internal/core/`):**
   - Create a new utility function (e.g., `ExecuteHook(eventName string)`).
   - Check if `~/.mindloop/hooks/{eventName}` exists and is executable.
   - Use `os/exec` to run the script asynchronously.
2. **Integration:**
   - Call `ExecuteHook("on-focus-start")` inside the focus session start logic.
   - Call `ExecuteHook("on-focus-end")` when a session finishes or is stopped.
3. **Documentation:** Update `docs/FEATURES.md` to explain how to use hooks (e.g., a simple bash script to toggle macOS DND using `shortcuts`).
