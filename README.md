# ⚡ LazyCurl

> **The Terminal Interface for Curl, Stress Testing, and API Exploration.**

LazyCurl is a powerful Terminal User Interface (TUI) for HTTP API interaction, written in Go. Inspired by `lazygit`, it brings the ergonomics of a modern TUI to `curl`, adding robust session management, visual request editing, and a built-in concurrent load testing engine.

## ✨ Features

- **Keyboard-First Design**: Navigate, edit, and execute requests without touching your mouse.
- **Interactive Editor**: Edit Method, URL, Headers, and JSON Body with full TUI support.
- **Project-Based**: Organizes requests into a session list for easy switching.
- **🚀 Load Testing Engine**: 
    - Built-in concurrent runner.
    - **Real-time Latency Dashboard**: Watch response times plotted live in ASCII.
    - Live statistics (Status codes, throughput, max/avg latency).
- **Curl-Powered**: Uses native `curl` under the hood for maximum compatibility and reliability.

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/varadrane7/lazycurl.git
cd lazycurl

# Built with Go 1.21+
go build -o lazycurl main.go

# Run
./lazycurl tui
```

## 🎮 Controls

### Global
- `Tab` / `Shift+Tab`: Switch Panes (Requests <-> Editor <-> Response)
- `q` / `Ctrl+C`: Quit

### Requests Pane (Left)
- `j` / `k` (or Arrows): Navigate requests
- `n`: Create new request
- `d`: Delete request
- `Enter`: Select request to edit

### Editor Pane (Middle)
- **Navigation**:
    - `j` / `k`: Move between Method, URL, Tabs, and Content.
    - `Tab Bar`: Use `h` / `l` (Left/Right) to switch between **[Headers]**, **[Body]**, and **[Load]**.
- **Editing**:
    - `Enter`: Enter Edit Mode (Focus field).
    - `Esc`: Exit Edit Mode (Save & Blur).
- **Headers Tab**:
    - `n`: Add new header.
    - `d`: Delete header.
- **Load Tab**:
    - Set **Concurrency** (workers) and **Duration** (e.g., `10s`).

### Execution
- `r`: **Run Request** (or Start Load Test if in Load Tab).

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **CLI**: [Cobra](https://github.com/spf13/cobra)
- **Graphing**: [AsciiGraph](https://github.com/guptarohit/asciigraph)

## 📄 License

MIT
