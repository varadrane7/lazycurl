package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the application.
type KeyMap struct {
	// Global
	Quit     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Run      key.Binding
	Help     key.Binding

	// Requests Pane
	Up     key.Binding
	Down   key.Binding
	New    key.Binding
	Delete key.Binding

	// Editor Pane
	EditEnter key.Binding // Enter edit mode
	EditEsc   key.Binding // Exit edit mode
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q", "quit"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next pane"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev pane"),
		),
		Run: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "run request"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new request"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d", "delete"), // 'd' might conflict if we allow typing, but in nav mode it is fine
			key.WithHelp("d", "delete"),
		),
		EditEnter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "edit field"),
		),
		EditEsc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit edit"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Tab, k.Run}
}

// FullHelp returns keybindings for the expanded help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.New, k.Delete},    // Requests
		{k.Tab, k.ShiftTab, k.Run, k.Quit}, // Global
	}
}
