# AI-Powered Journal Generation - Feature Discovery & PRD

## 1. Goal
Enable users to automatically generate comprehensive daily, weekly, or yearly journal entries summarizing their accomplishments. This feature will leverage an LLM to digest activity data (Focus Sessions, Habits, Tasks, Intents, and Points) and produce a cohesive narrative or statistical summary, which can then be saved as a standard Journal entry. 

The prompt will be highly optimized, analytical, and ADHD-friendly to help users understand exactly what they did and accomplished across all features.

## 2. User Experience / Flow

### 2.1 Configuration & Token Management
- **CLI:** Users can provide their token via an environment variable (e.g., `MINDLOOP_AI_TOKEN`).
- **Web UI:** Users configure their AI provider and token in the `Settings` page.
- **Token Storage:** To avoid requiring the user to repeatedly enter the API key in the UI, the token will be encrypted and stored securely in the database.
- **Provider & Model Selection:** Users will select their preferred provider (e.g., Google, OpenAI, Anthropic, local/Ollama). The system will use provider-specific endpoints to dynamically fetch and display the models available to the user based on their API key.

### 2.2 Generation Flow (CLI)
- Trigger generation using flags: `mindloop journal generate -d` (daily), `-w` (weekly), or `-y` (yearly).
- **Recommendation in Summary:** Running `mindloop summary -d` (or `-w`, `-y`) will include a recommendation: *For an AI-generated overview, use `mindloop journal generate -d`*.
- **Interactive Save:** After displaying the generated summary, the CLI will interactively ask: `Would you like to save this into journal?`. If yes, it saves the entry with appropriate tags (e.g., `#ai-summary`, `#weekly-review`).

### 2.3 Generation Flow (Web UI)
- In the `Journal` feature (and optionally the `Summary` view), a new action (e.g., "✨ Auto-generate Entry") becomes available.
- The user selects the period (Daily, Weekly, Yearly).
- The generated text is presented in an editor. The user can review/edit it, and choose to save it as a journal entry with auto-generated tags.

## 3. Architecture & Data Flow

### 3.1 Context Gathering & Pre-aggregation
- Introduce a new service method `GatherContext(timeframe string)` that aggregates data from the selected period.
- To optimize for lower-end models (e.g., Claude Haiku, Gemini Flash, local Ollama models) and avoid exceeding token limits, the data will be pre-aggregated before sending to the LLM. For instance, yearly summaries will aggregate data into monthly chunks rather than sending thousands of raw task entries.
- Data gathered includes: FocusSessions, HabitLogs, Tasks, Intents, SideQuests, and Points.

### 3.2 Security & Encryption
- Implement encryption/decryption utilities for the API tokens before storing them in the database to ensure they are protected at rest.

### 3.3 Prompt Engineering
- A hardcoded, system-optimized prompt designed to be analytical and ADHD-friendly. It will instruct the LLM to provide a clear, structured overview of everything the user accomplished.

## 4. Out of Scope (V1)
- Multi-modal input (e.g., voice-to-text journaling).
- Automatic generation in the background without user initiation.
- Custom user-defined system prompts.
- Fine-tuning local models.
