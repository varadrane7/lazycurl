Below is a **lazygit-faithful UX specification**: panes, navigation, flows, and keybindings. This is written so you (or a future contributor) can implement directly with a Bubble Tea–style model/update/view loop.

---

# LazyCurl — UX Specification (TUI)

## 1. UX Principles (borrowed directly from lazygit)

1. **Everything is stateful**

   - The UI reflects app state, not command output
   - Actions mutate state; views re-render automatically

2. **Focus defines meaning**

   - Keys do different things depending on focused pane
   - User never needs to remember global flags mid-flow

3. **Minimal prompts, maximum flow**

   - Inline confirmations
   - Contextual help at bottom
   - No modal overload

4. **Predictable muscle memory**

   - hjkl / arrows
   - `enter` to run/open
   - `d` destructive, `e` edit, `r` run, `?` help

---

## 2. Default Screen Layout (MVP)

```
┌─────────────────────┬────────────────────────┬──────────────────────────┐
│ Requests / History  │ Request Editor         │ Response Viewer          │
│                     │                        │                          │
│ > GET /users        │ Method: GET            │ Status: 200 OK           │
│   POST /login       │ URL: /users            │ Time: 123ms              │
│   load-test-1       │ --------------------   │ Size: 2.4 KB             │
│                     │ Headers | Body | Auth  │ ----------------------   │
│                     │                        │ {                        │
│                     │                        │   "id": 1,               │
│                     │                        │   "name": "Alice"        │
│                     │                        │ }                        │
├─────────────────────┴────────────────────────┴──────────────────────────┤
│  r: run  e: edit  tab: switch pane  ?: help  env: dev                   │
└─────────────────────────────────────────────────────────────────────────┘
```

### Pane Roles

1. **Left Pane** – Navigation

   - Requests
   - Collections (Postman)
   - History
   - Load test configs

2. **Middle Pane** – Editor

   - Request configuration
   - Tabbed subviews

3. **Right Pane** – Output

   - Response body
   - Headers
   - Timing / metrics
   - Logs (load tests)

4. **Bottom Bar**

   - Key hints (context-sensitive)
   - Status messages
   - Active environment/session

---

## 3. Pane Focus & Navigation

### Global Navigation

| Key         | Action                 |
| ----------- | ---------------------- |
| `tab`       | Cycle focus forward    |
| `shift+tab` | Cycle focus backward   |
| `h/j/k/l`   | Move within pane       |
| `↑ ↓ ← →`   | Alternative navigation |
| `q`         | Quit                   |
| `?`         | Help overlay           |

Focused pane is visually highlighted (border + title color).

---

## 4. Left Pane: Requests / Collections / History

### Structure

- Tree-based

  - Collections

    - Folder

      - Request

  - Saved Requests
  - History (chronological)
  - Load Tests

### Keys (Left Pane)

| Key     | Action                  |
| ------- | ----------------------- |
| `enter` | Select item             |
| `r`     | Run selected request    |
| `e`     | Edit request            |
| `d`     | Delete (confirm inline) |
| `/`     | Search/filter           |
| `n`     | New request             |
| `i`     | Import (Postman)        |

Selection updates **middle pane** automatically.

---

## 5. Middle Pane: Request Editor

### Tabs

- **Params**
- **Headers**
- **Body**
- **Auth**
- **Settings**

Tabs are horizontal, like lazygit’s subviews.

### Keys (Middle Pane)

| Key         | Action             |
| ----------- | ------------------ |
| `tab`       | Next field / tab   |
| `shift+tab` | Previous field     |
| `e`         | Edit current field |
| `ctrl+s`    | Save request       |
| `r`         | Run request        |
| `esc`       | Cancel edit        |

### Editing Model

- Inline editing (no modal)
- JSON body:

  - auto-format on save
  - validation feedback inline

---

## 6. Right Pane: Response Viewer

### Subviews

- **Body** (default)
- **Headers**
- **Timing**
- **Errors / Logs**

### Body View

- Auto-detect JSON
- Pretty-printed
- Collapsible objects (MVP: fold large blocks)
- Scrollable

### Timing View

```
DNS:      12ms
Connect:  20ms
TLS:      30ms
TTFB:     80ms
Total:   123ms
```

### Keys (Right Pane)

| Key | Action                   |
| --- | ------------------------ |
| `b` | Body view                |
| `h` | Headers view             |
| `t` | Timing view              |
| `y` | Copy response            |
| `s` | Save response to history |

---

## 7. Load Test UX

### Load Test Editor (Middle Pane)

Fields:

- Concurrency
- Duration
- Ramp-up
- Optional RPS limit

### Load Test Viewer (Right Pane)

Live updating:

```
Requests/sec:  820
Errors:        0.3%
p50:           120ms
p95:           280ms
p99:           410ms
```

### Keys

| Key | Action      |
| --- | ----------- |
| `r` | Start test  |
| `x` | Stop test   |
| `s` | Save config |
| `l` | View logs   |

UI must **never freeze** during load tests.

---

## 8. Help & Discoverability

### Help Overlay (`?`)

- Shows:

  - Pane-specific keys
  - Global keys

- Closes with `esc` or `?`

### Bottom Bar

- Always shows **current pane actions**
- Example:

  ```
  r: run  e: edit  d: delete  /: search
  ```

---

## 9. Error & Confirmation UX

### Errors

- Inline, non-blocking
- Shown in bottom bar or right pane

### Destructive Actions

- Inline confirmation (lazygit-style)

  ```
  Delete request? (y/n)
  ```

No modal dialogs.

---

## 10. UX Acceptance Criteria

- User can:

  - run a request in ≤3 keystrokes
  - inspect headers/body/timing without leaving TUI
  - repeat a request faster than Postman

- UI remains responsive during load testing
- No mouse required
- No blocking popups

---

## 11. UX → Architecture Implications (Important)

To support this UX cleanly, the app **must** have:

- Central app state (focused pane, selected item)
- Event → command mapping (keys → actions)
- Views render **only from state**
- Execution (curl, load tests) runs async, streams updates

This maps cleanly to:

- Bubble Tea’s `model/update/view`
- Lazygit’s command-driven architecture

---
