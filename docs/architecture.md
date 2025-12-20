# LazyCurl Architecture

## 1) High-level design

### Key architectural goals

- **Single source of truth state** (like lazygit): selection, focus, active env, last run, errors.
- **Command-driven side effects**: keypress → action → command(s) → state mutation + async jobs.
- **Non-blocking UI**: request execution and load testing never run on the render thread.
- **Deterministic curl compilation**: internal `RequestSpec` → `CurlInvocation` (string + args) → execution → parsed result.
- **Config-first**: YAML config for keybindings, theme, defaults, custom commands later.

### Core layers

1. **UI** (Bubble Tea model/update/view; views are pure render of state)
2. **App State** (domain models + UI state: focus, selection, filters)
3. **Actions/Commands** (lazygit-style: action handlers create commands)
4. **Engine** (curl runner + load runner + parsers)
5. **Persistence** (history/session storage + imports)
6. **Interop** (Postman import, export later)

---

## 2) Process topology

### Single process, multi-goroutine

- Main goroutine runs TUI event loop.
- Worker goroutines:

  - **curl execution** jobs
  - **load test** workers + aggregator
  - **history writes** (async flush)

Communication via:

- Bubble Tea messages (preferred) OR internal event bus channel that converts to tea.Msg.

---

## 3) Suggested repository structure (Go)

```
lazycurl/
  cmd/
    root.go
    tui.go
    run.go
    load.go
    session.go
    import_postman.go

  internal/
    app/
      model.go          // AppModel: full UI/app state
      update.go         // Bubble Tea update: routes msgs -> reducers/actions
      view.go           // Root layout + pane composition
      navigation.go     // Focus + selection + keybinding routing
      keymap/           // Default keymaps + config mapping
      commands/         // Action handlers -> command constructors (lazygit style)

    ui/
      panes/
        leftnav/        // tree lists: requests/history/collections
        editor/         // request editor tabs
        response/       // response viewer: body/headers/timing/logs
        statusbar/      // bottom bar + notifications
      components/       // reusable widgets, tables, spinners, search, etc.
      theme/

    domain/
      request.go        // RequestSpec, AuthSpec, BodySpec, Vars
      response.go       // Response, Timing, Errors
      session.go        // Session, Env, Vars
      history.go        // HistoryEntry
      load.go           // LoadProfile, LoadResult
      collection.go     // Collection tree nodes

    engine/
      curl/
        build.go        // RequestSpec -> CurlInvocation
        exec.go         // run curl process; streaming; context/timeouts
        parse.go        // parse stdout/stderr/write-out blocks
        metrics.go      // normalize timing fields, size, status
      load/
        runner.go       // orchestrates N workers + rate limiter
        worker.go       // worker loop using engine/curl
        stats.go        // histograms, percentiles, rolling stats
        report.go       // final summaries + JSON output
      jsonx/
        detect.go
        pretty.go
        extract.go      // MVP: simple dot-path extraction

    store/
      sqlite/ or bolt/
        db.go
        history_repo.go
        session_repo.go
        request_repo.go
      migrate/

    interop/
      postman/
        import.go       // parse collection v2.1
        map.go          // postman -> domain.RequestSpec
        env.go          // environments -> vars

    config/
      config.go         // load/merge defaults + YAML + env
      defaults.go

    util/
      errx/
      logx/
      osx/
      textx/
      timex/

  docs/
    project-brief.md
    prd.md
    ux.md (optional)
  go.mod
  main.go
```

**Why this shape works**

- UI code doesn’t know about curl details.
- Engine doesn’t know about panes.
- Commands/actions glue UI → engine → store.

---

## 4) App state model (single source of truth)

### `internal/app/model.go` (conceptual)

- `Focus` (left/editor/response)
- `LeftNavState` (active section: requests/history/collections/load; filter string; selection path)
- `EditorState` (active tab; draft edits; dirty flag)
- `ResponseState` (active subview: body/headers/timing/logs; scroll offsets; last response)
- `RunState`

  - running bool
  - spinner/progress
  - last error

- `LoadState`

  - running bool
  - live stats snapshot

- `ActiveSession` (env vars, base URL, auth defaults)
- `Config` (theme, keybindings, defaults)

Everything is rendered from this state.

---

## 5) Action/Command architecture (lazygit-inspired)

### Pattern

- Keypress → `Action` → returns a list of `Command`s:

  - `StateCommand` (pure reducer)
  - `ExecCommand` (starts async engine job; emits messages)
  - `StoreCommand` (persist history/session)
  - `NavCommand` (focus changes, selection updates)

### Why

- Keeps `update()` small and testable.
- Side effects are isolated and mockable.

### Example flow: `r` Run request

1. Action: `RunSelectedRequest`
2. Commands:

   - `SetRunning(true)`
   - `ExecCurl(requestSpecResolved)`

3. Engine returns `CurlResultMsg`
4. Reducer updates `ResponseState` + `History` + `SetRunning(false)`

---

## 6) Curl execution subsystem

### Deterministic curl compilation

`domain.RequestSpec` → `engine/curl.CurlInvocation{Args []string, Env map, Redacted string}`

- Use `--write-out` to append a **unique delimiter block** at end of output so parsing is reliable.
- Prefer:

  - `--silent --show-error --location` (configurable)
  - `--dump-header -` optional (or `-D -`)
  - `--output -` capture body
  - `--compressed`

### Parsing strategy (robust)

- Output format:

  - stdout: body + `\n<DELIM>\n` + JSON metrics line
  - stderr: error lines (capture as well)

- Parse by splitting on delimiter from the end (safe if body contains junk).
- Normalize timing keys into `domain.Timing`.

### Concurrency guardrails

- Default max workers (e.g., 32) to avoid fork bombs.
- Load testing uses worker pool with queue; doesn’t spawn unbounded curl processes.

---

## 7) Load testing engine

### Components

- `load.Runner` orchestrates:

  - worker pool
  - optional rate limiter
  - shared context cancellation

- `load.Stats` maintains:

  - rolling counters
  - latency histogram (fixed buckets) to compute p50/p95/p99 efficiently

- UI updates:

  - stats snapshot published at fixed tick (e.g., 200ms) to avoid flooding UI

### Data model

- `LoadProfile{Concurrency, Duration, Ramp, RPSLimit, Warmup}`
- `LoadSnapshot{RPS, ErrRate, P50, P95, P99, InFlight, Total}`
- `LoadReport{Summary, Histogram, ErrorsByType}`

---

## 8) Persistence layer

### What to persist

- Requests (saved)
- Sessions (env + defaults)
- History entries (request + response metadata)

### Store choice

- **SQLite**: best for query/search/history filtering
- **BoltDB/bbolt**: simpler embed, but less query power

Recommendation:

- Start with **SQLite** if you want great `/` search in history quickly.
- Start with **bbolt** if you want minimal dependencies.

(Architecture supports either; define `store` interfaces.)

### Repos (interfaces)

- `HistoryRepo.Append(entry)`
- `HistoryRepo.Search(filter, limit)`
- `SessionRepo.Save(session)`
- `RequestRepo.Save(request)`

---

## 9) Postman import subsystem

### MVP mapping rules

- Collection tree → left pane tree nodes
- Each request becomes `domain.RequestSpec`
- Variables:

  - collection vars + env vars → `Vars` resolver chain

- Unsupported in MVP:

  - pre-request/test scripts
  - auth flows beyond basic translation

### Variable resolution order (clear + predictable)

1. Runtime overrides (CLI flags / quick set)
2. Active session env vars
3. Postman environment vars
4. Postman collection vars
5. Request-local vars

---

## 10) Config + keybindings

### Config load order

1. Built-in defaults
2. `~/.config/lazycurl/config.yml` (or OS-specific)
3. Project-local `.lazycurl.yml` (optional)
4. Env vars overrides
5. CLI flags override

### Keybinding mapping

- `config.Keybindings` maps action IDs → keys
- UI binds keys to actions via keymap layer (like lazygit)

---

## 11) Performance & responsiveness rules (non-negotiable)

- UI renders from state only; no blocking I/O in `View()`
- Engine jobs always emit messages back; never mutate UI state directly
- Stats updates are throttled (tick-based)
- Curl output parsing must be linear-time; avoid regex-heavy parsing in hot loops

---

## 12) Implementation sequence (architect’s “thin slices”)

1. **TUI skeleton**

   - panes + focus + status bar + help overlay

2. **Request model + curl builder**

   - `run` action works in TUI; show response body + timing

3. **History persistence**

   - append last run; browse in left pane

4. **Load runner**

   - concurrency + duration + snapshots + right-pane dashboard

5. **Postman import**

   - populate collections tree; run imported request

6. **Config keymaps/theme**

   - lazygit-like customization
