# Mindloop Roadmap 2026

> **Status Update**: Development is currently in a **maintenance mode**. New features will be added slowly as the core stability and user experience are prioritized. Expect a slower release cycle.

## 🚀 CLI Roadmap
The command-line interface is the heart of Mindloop. Future work will focus on speed, scriptability, and deeper integrations.

- [ ] **Interactive Dashboard**: A TUI (Text User Interface) mode listing active intents and valid commands.
- [ ] **Natural Language Processing**: Basic NLP for intent parsing (e.g., `mindloop intent "Do the thing tomorrow"`).
- [ ] **Plugins/Extensions**: Allow user scripts to hook into Mindloop events (e.g., specific hooks on `intent start`).
- [ ] **Sync/Export**: Better export formats (JSON, CSV, Markdown) for data portability.

## 🎨 UI Roadmap
The Web UI brings Mindloop to a wider audience. The goal is a premium, "vibe-coded" experience that feels native.

- [ ] **PWA Support**: Full Progressive Web App capabilities for install-on-mobile.
- [ ] **Offline Mode**: Service Workers to cache assets and allow basic functionality without a server connection (sync later).
- [ ] **Themes**: User-customizable color themes beyond just Light/Dark.
- [ ] **Keyboard Shortcuts**: Global hotkeys for quick actions (e.g., `Cmd+K` for a command palette).
- [ ] **Visualizations**: deeper analytics on Focus Sessions and Habits (graphs, trends).

## 🏗 Architecture Roadmap
Building a robust foundation for long-term maintainability.

- [ ] **Plugin System**: A defined interface for community plugins.
- [ ] **API Expansion**: Full REST API coverage for all internal features.
- [ ] **Testing**: Increase unit test coverage for core logic (Intents, Habits).
- [ ] **Performance**: Optimize SQLite queries and template rendering for larger datasets.
- [ ] **Containerization**: Official Docker image for easier deployment.

## ✅ Recently Completed
- [x] **UI Polish**: Massive consistency update (Alignment, Icons, Spacing).
- [x] **Visuals**: Dark Mode, Glassmorphism elements.
- [x] **Features**: Side Quest Modal, Habit Heatmaps, stacked Settings.
