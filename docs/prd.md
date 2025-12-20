# Product Requirements Document (PRD): LazyCurl

## 1. Product Overview

### Product Name

**LazyCurl**

### Description

LazyCurl is a **terminal-based UI (TUI) and CLI application** for API exploration, testing, and stress testing, built in **Golang** using **Cobra** and powered by **curl** as the execution engine.

Inspired heavily by **lazygit**, LazyCurl provides a **keyboard-first, pane-based interface** that enables fast, repeatable API workflows with minimal friction.

### Primary Value Proposition

> “Everything you like about curl’s power, with the flow and ergonomics of lazygit.”

---

## 2. Target Users

### Primary Users

- Backend / API engineers
- SRE / Platform engineers
- QA / Automation engineers

### Secondary Users

- Developers migrating from Postman to terminal-first workflows
- CI / automation users who want reproducible API runs

---

## 3. Goals & Success Criteria

### Goals

1. Provide a **lazygit-like UX** for API workflows
2. Enable **faster iteration loops** than GUI tools like Postman
3. Support **high-performance API load testing**
4. Maintain **curl compatibility and flexibility**
5. Offer **first-class JSON ergonomics**

### Success Metrics

- Time to repeat a request (edit → run → inspect) is significantly faster than Postman
- TUI remains responsive during concurrent load tests
- Postman collections import and run with minimal manual fixes
- Users can complete common API workflows without leaving the terminal

---

## 4. Non-Goals (MVP Scope Guardrails)

- Full Postman feature parity (scripts, monitors, mocks)
- GUI or web-based UI
- Replacing curl with a custom HTTP engine in MVP
- Cloud sync or team collaboration features

---

## 5. User Experience Requirements (TUI-first)

### Layout (Initial MVP)

- **Left Pane**: Requests / Collections / History (tree view)
- **Middle Pane**: Request editor (tabs for Params, Headers, Body, Auth)
- **Right Pane**: Response viewer (Body / Headers / Timing)
- **Bottom Bar**: Key hints, active environment, status messages

### Interaction Model

- Keyboard-first navigation (hjkl / arrows / tab)
- Focus-based panes (actions depend on focused pane)
- Contextual keybindings (similar to lazygit)
- Global actions:

  - Help
  - Search
  - Quit
  - Command palette (stretch goal)

### Accessibility

- Works over SSH
- No mouse required
- Clear focus indicators

---

## 6. Functional Requirements

### 6.1 CLI Commands (Cobra)

#### Core Commands

- `lazycurl tui`
  Launch interactive TUI (primary experience)

- `lazycurl run`
  Run a single request (scriptable, non-interactive)

- `lazycurl load`
  Run a load/stress test

- `lazycurl import postman`
  Import Postman collection(s)

- `lazycurl session`
  Save, list, load sessions/environments

---

### 6.2 Request Execution (Curl Backbone)

- All requests compile into deterministic curl commands
- Curl invoked with:

  - silent mode
  - structured `--write-out` for metrics

- Capture:

  - HTTP status
  - total time
  - DNS/connect/TLS/TTFB timings (when available)
  - response size

- Support:

  - headers
  - body
  - auth
  - query params
  - cookies
  - TLS/proxy flags

---

### 6.3 JSON Ergonomics

- Auto-detect JSON responses
- Pretty-print JSON by default (configurable)
- Response summary:

  - status
  - latency
  - size

- Basic JSON inspection:

  - collapse/expand
  - top-level key summary

- (Future): jq/JSONPath-style filtering

---

### 6.4 History & Sessions

- Persist:

  - request definitions
  - resolved variables
  - response metadata

- History is:

  - searchable
  - navigable from TUI

- Sessions include:

  - base URL
  - auth
  - headers
  - environment variables

---

### 6.5 Load / Stress Testing

#### Configuration

- Concurrency (number of workers)
- Duration
- Optional ramp-up
- Optional RPS cap (token bucket)

#### Metrics

- Throughput
- Latency percentiles (p50 / p95 / p99)
- Error rates (network vs HTTP)

#### UX

- Live updating metrics in TUI
- Final summary table
- JSON output for CI

---

### 6.6 Postman Interoperability (MVP)

- Import Postman Collection v2.1
- Map:

  - requests
  - folders
  - variables

- Explicitly **ignore**:

  - pre-request scripts
  - test scripts (MVP)

- Allow:

  - run a single request
  - run a folder sequentially

---

## 7. Configuration System

### Config File (YAML)

- Theme
- Keybindings
- Default headers/auth
- Curl flags
- UI preferences

### Philosophy

- Opinionated defaults
- Everything overridable (lazygit-style)

---

## 8. Technical Constraints

- Cross-platform (macOS, Linux, Windows)
- Curl must be available on system
- Go-native TUI library (Bubble Tea–style architecture recommended)
- Efficient process spawning and worker pooling

---

## 9. Risks & Mitigations

| Risk                     | Mitigation                          |
| ------------------------ | ----------------------------------- |
| Curl process overhead    | Worker pool, concurrency limits     |
| UI lag during load       | Decouple render loop from execution |
| Postman complexity       | Strict MVP scope                    |
| Terminal inconsistencies | Conservative rendering patterns     |

---

## 10. Milestones

1. TUI shell + pane layout + navigation
2. Request execution + response viewer
3. JSON pretty + timing metrics
4. History + session persistence
5. Load testing engine + live dashboard
6. Postman import
7. Config + keybinding customization

---

## 11. Open Questions (Deferred)

- Should LazyCurl ever switch from curl to libcurl or native HTTP?
- How extensible should custom commands/macros be in v1?
- Plugin system vs config-only extensibility?
