# Issue 4: [Core/UI] NLP Smart Parsing for Tasks and Intents

**Labels:** `enhancement`, `core`, `nlp`

### Rationale & Market Fit
When quick-capturing an idea, users don't want to fiddle with date pickers or dropdowns. Natural Language Processing (NLP) allows users to type "Review PRs tomorrow at 10am" and automatically sets the correct due date.

### Expected Outcomes
- The Web UI Task creation input and CLI `mindloop task add` commands support natural language dates.
- Seamless, intuitive input parsing.

### Implementation Details
1. **Parsing Logic:**
   - Integrate a lightweight NLP date parser for Go, such as [tj/go-naturaldate](https://github.com/tj/go-naturaldate) or similar.
   - When a task title is submitted, scan it for time-based keywords.
2. **Data Layer:**
   - Extract the date and strip it from the title. 
   - E.g. "Review PRs tomorrow" -> Title: "Review PRs", Due Date: [Tomorrow's Date].
3. **UI Feedback:**
   - If using the Web UI, optionally highlight the parsed date in the input field as the user types (stretch goal), but backend parsing on submit is perfectly fine for v1.
