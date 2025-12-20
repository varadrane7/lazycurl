# Project Brief: LazyCurl — a Lazygit-inspired “Fancy curl” for APIs

## 1) Summary

**LazyCurl** is a **terminal UI (TUI) + CLI** application built in **Golang**, using **Cobra** for command structure and **curl as the execution backbone**, designed to make API exploration and stress testing feel as fluid as using **lazygit** for git.

LazyCurl emphasizes:

- **Keyboard-first**, pane-based UI with predictable workflows
- **Fast iteration loops** for API requests (“edit → run → inspect → repeat”)
- **Load/stress testing** with live throughput/latency visualization
- **JSON-native ergonomics**
- **Session/history** and **Postman Collection interoperability**

## 2) Inspiration and Product Ethos (Lazygit-style)

LazyCurl borrows heavily from lazygit’s product principles:

- **No ceremony**: minimal friction, instant feedback
- **Opinionated defaults** that work well, but deeply configurable
- **Contextual actions** tied to the currently focused pane/item
- **Stateful TUI** where navigation and actions are muscle-memory friendly
- **Config-driven customization** (YAML), including keybindings and custom commands
- **Composable “commands”** as the unit of behavior (actions are predictable + scriptable)

## 3) Problem

Postman-class tools are powerful but often slow/heavy for power users and less natural for shell workflows/CI. Curl is fast and ubiquitous but verbose and awkward for:

- repeated workflows
- JSON inspection
- maintaining history/sessions
- load testing
- running collections

## 4) Target Users

Primary:

- Backend/API engineers
- SRE/Platform engineers
- QA/automation engineers

Secondary:

- Developers migrating from Postman to terminal-first workflows

## 5) Goals

- **Lazygit-like experience**, but for HTTP APIs:

  - pane-based navigation
  - contextual keybindings/actions
  - fast “flow state”

- **Performance-first**:

  - quick startup
  - high throughput load generation
  - efficient rendering (no lag during stress tests)

- **Curl-powered** MVP:

  - keep compatibility and feature breadth of curl

- **Ergonomic JSON tooling**:

  - pretty output, structured inspection, light query/extract

- **Interop**:

  - import Postman Collections (and environments when feasible)

## 6) Non-Goals (initially)

- Full Postman parity (scripting, monitors, mocks, cloud sync)
- Replacing curl entirely with a custom HTTP stack (MVP stays curl-backed)
- GUI app

## 7) Core Use Cases

1. **Explore APIs in a lazily productive way**

   - browse saved requests, environments, and history
   - edit request parts quickly (headers/body/query/auth)

2. **Run and inspect**

   - run request; see status/time breakdown; inspect body/headers

3. **Stress test**

   - configure concurrency/duration/ramp/RPS; view live metrics and final report

4. **Session workflows**

   - save “workspaces” (base URL, auth, headers, variables)

5. **Postman import**

   - load a collection; run a folder/request; map vars to env

## 8) UX Model (TUI-first)

### Pane layout (initial concept)

- **Left pane:** Collections / Requests (tree)
- **Middle pane:** Request editor / Params / Headers / Auth tabs
- **Right pane:** Response viewer (Body / Headers / Timing / Logs)
- **Bottom bar:** key hints + status + active environment

### Interaction model (lazygit-inspired)

- **Focus-based navigation** (hjkl/arrows/tab cycling)
- **Contextual actions** (keys change meaning by focused pane)
- **Global actions** (search, command palette, help, quit)
- **Action confirmation pattern** for destructive ops (e.g., clearing history)

## 9) Architectural Goal: “Lazygit-like” internal structure

LazyCurl will mirror lazygit-style layering:

- **Models/state**: app state, focused pane, selected request, active environment
- **UI layer**: panes/views render from state
- **Controllers/handlers**: keybindings/events map to actions
- **Command system**: actions produce commands that mutate state and/or execute side effects (run curl, parse output, store history)
- **Config system**: YAML config for:

  - keybindings
  - themes
  - custom commands/macros (stretch goal for MVP)

### Curl execution backbone

- A normalized internal `Request` model compiles into a deterministic curl invocation.
- Curl returns:

  - response body (optional capture)
  - headers (optional)
  - structured timing metrics (via `--write-out`)

- The engine parses and streams results into:

  - live UI updates (TUI)
  - structured output (CLI mode)

## 10) Key Features (MVP)

- Cobra CLI:

  - `lazycurl tui` (primary experience)
  - `lazycurl run` (single request; script-friendly)
  - `lazycurl load` (stress test)
  - `lazycurl import postman`
  - `lazycurl session` (save/list/load)

- JSON ergonomics:

  - auto pretty print
  - response summary (status, size, latency)

- History + sessions:

  - searchable recent runs
  - saved environments/variables

- Load testing:

  - concurrency + duration (MVP)
  - aggregated latency percentiles + error breakdown

## 11) Success Metrics

- “Time to repeat run” (edit → run → inspect) feels comparable to lazygit’s speed for git actions
- Load tests run with stable UI responsiveness
- Postman collection import works for common real-world collections with minimal edits

## 12) Risks

- Curl process spawning overhead for very high concurrency (mitigate via worker caps + scheduling)
- Postman import complexity (variable scoping, scripts) → strictly scope MVP
- Cross-platform terminal quirks (rendering, key handling)

## 13) Milestones (high level)

1. TUI skeleton (panes + navigation + keybinding help) + `run` integration
2. Response viewer with timing + JSON pretty
3. History + session persistence
4. `load` engine + live dashboard
5. Postman import and runnable requests
6. Configurable keybindings/themes (lazygit-like polish pass)
