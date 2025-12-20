package tui

import (
	"lazycurl/internal/load"
	"lazycurl/internal/model"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global Keybindings (only when NOT editing, or special keys)
		if !m.IsEditing {
			if key.Matches(msg, m.KeyMap.Quit) {
				return m, tea.Quit
			}
			if key.Matches(msg, m.KeyMap.Tab) {
				m.ActivePane = (m.ActivePane + 1) % 3
				return m, nil
			}
			if key.Matches(msg, m.KeyMap.ShiftTab) {
				m.ActivePane = (m.ActivePane - 1 + 3) % 3
				return m, nil
			}
			if key.Matches(msg, m.KeyMap.Run) {
				// Special handling if in Load Tab
				if m.ActivePane == PaneEditor && m.ActiveEditorTab == TabLoad {
					m.SyncRequestToEditor()

					// Parse inputs
					conc, _ := strconv.Atoi(m.LoadConfig.Concurrency.Value())
					if conc <= 0 {
						conc = 1
					}
					dur, err := time.ParseDuration(m.LoadConfig.Duration.Value())
					if err != nil {
						dur = 5 * time.Second
					}

					m.LoadState.IsRunning = true
					m.LoadState.Stats = load.NewStats()

					runner := load.NewRunner()
					req := m.Requests[m.SelectedReqIdx]

					// Start
					ch := runner.Run(req, conc, dur)
					m.LoadState.Sub = ch

					m.ActivePane = PaneResponse // Switch to dashboard

					return m, WaitForStats(ch)
				} else {
					m.SyncRequestToEditor()
					return m, m.RunRequestCmd
				}
			}
		} else {
			// Special handling while editing
			if key.Matches(msg, m.KeyMap.EditEsc) {
				m.IsEditing = false
				// Blur all
				m.EditorInputs[0].Blur()
				m.EditorInputs[1].Blur()
				m.EditorBody.Blur()
				for i := range m.HeaderInputs {
					m.HeaderInputs[i].Key.Blur()
					m.HeaderInputs[i].Value.Blur()
				}
				m.LoadConfig.Concurrency.Blur()
				m.LoadConfig.Duration.Blur()

				m.SyncRequestToEditor() // Save on exit edit
				return m, nil
			}
		}

		// Pane Specific Handling
		switch m.ActivePane {
		case PaneRequests:
			m, cmd = m.updateRequests(msg)
			cmds = append(cmds, cmd)
		case PaneEditor:
			m, cmd = m.updateEditor(msg)
			cmds = append(cmds, cmd)
		case PaneResponse:
			// Scrolling logic could go here
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Help.Width = msg.Width
	case load.StatsMsg:
		m.LoadState.Stats = msg.Stats
		if msg.Done {
			m.LoadState.IsRunning = false
			m.LoadState.Sub = nil // Clear channel
		} else {
			return m, WaitForStats(m.LoadState.Sub) // Wait for next
		}
	case model.Response:
		m.Response = &msg
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateRequests(msg tea.KeyMsg) (Model, tea.Cmd) {
	if key.Matches(msg, m.KeyMap.Up) {
		if m.SelectedReqIdx > 0 {
			m.SelectedReqIdx--
			m.SyncEditorToRequest()
		}
	} else if key.Matches(msg, m.KeyMap.Down) {
		if m.SelectedReqIdx < len(m.Requests)-1 {
			m.SelectedReqIdx++
			m.SyncEditorToRequest()
		}
	} else if key.Matches(msg, m.KeyMap.New) {
		newReq := model.NewRequest()
		newReq.URL = "https://example.com/new"
		m.Requests = append(m.Requests, newReq)
		m.SelectedReqIdx = len(m.Requests) - 1
		m.SyncEditorToRequest()
	} else if key.Matches(msg, m.KeyMap.Delete) {
		if len(m.Requests) > 1 {
			m.Requests = append(m.Requests[:m.SelectedReqIdx], m.Requests[m.SelectedReqIdx+1:]...)
			if m.SelectedReqIdx >= len(m.Requests) {
				m.SelectedReqIdx = len(m.Requests) - 1
			}
			m.SyncEditorToRequest()
		}
	}
	return m, nil
}

func (m Model) updateEditor(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Editing Mode
	if m.IsEditing {
		// Hand off to bubbles
		switch m.FocusedField {
		case FieldMethod:
			m.EditorInputs[0], cmd = m.EditorInputs[0].Update(msg)
		case FieldURL:
			m.EditorInputs[1], cmd = m.EditorInputs[1].Update(msg)
		case FieldContent:
			if m.ActiveEditorTab == TabBody {
				m.EditorBody, cmd = m.EditorBody.Update(msg)
			} else if m.ActiveEditorTab == TabHeaders {
				// Header Editing
				idx := m.FocusedHeaderIdx
				if idx < len(m.HeaderInputs) {
					if m.FocusedHeaderKey {
						m.HeaderInputs[idx].Key, cmd = m.HeaderInputs[idx].Key.Update(msg)
					} else {
						m.HeaderInputs[idx].Value, cmd = m.HeaderInputs[idx].Value.Update(msg)
					}
				}
			} else if m.ActiveEditorTab == TabLoad {
				// Load Config Editing
				if m.FocusedHeaderIdx == 0 {
					m.LoadConfig.Concurrency, cmd = m.LoadConfig.Concurrency.Update(msg)
				} else {
					m.LoadConfig.Duration, cmd = m.LoadConfig.Duration.Update(msg)
				}
			}
		}
		return m, cmd
	}

	// Navigation Mode
	// Tab Switching (Global shortcuts inside editor)
	if key.Matches(msg, m.KeyMap.Tab) || key.Matches(msg, m.KeyMap.ShiftTab) {
		// Let the parent update loop handle pane switching
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.FocusedField == FieldContent && m.ActiveEditorTab == TabHeaders {
			if m.FocusedHeaderIdx > 0 {
				m.FocusedHeaderIdx--
				return m, nil // Handled inside content
			}
		}
		if m.FocusedField == FieldContent && m.ActiveEditorTab == TabLoad {
			if m.FocusedHeaderIdx > 0 {
				m.FocusedHeaderIdx--
				return m, nil
			}
		}
		if m.FocusedField > FieldMethod {
			m.FocusedField--
		}
	case "down", "j":
		if m.FocusedField == FieldContent && m.ActiveEditorTab == TabHeaders {
			// Sub-navigation inside headers list
			if m.FocusedHeaderIdx < len(m.HeaderInputs)-1 {
				m.FocusedHeaderIdx++
				return m, nil // Handled inside content
			}
		}
		if m.FocusedField == FieldContent && m.ActiveEditorTab == TabLoad {
			if m.FocusedHeaderIdx < 1 { // Only 2 inputs
				m.FocusedHeaderIdx++
				return m, nil
			}
		}
		if m.FocusedField < FieldContent {
			m.FocusedField++
		}
	case "left", "h":
		if m.FocusedField == FieldTabs {
			if m.ActiveEditorTab == TabHeaders {
				// Already leftmost
			} else if m.ActiveEditorTab == TabBody {
				m.ActiveEditorTab = TabHeaders // Move left
			} else if m.ActiveEditorTab == TabLoad {
				m.ActiveEditorTab = TabBody // Move left
			}
		} else if m.FocusedField == FieldContent && m.ActiveEditorTab == TabHeaders {
			m.FocusedHeaderKey = true
		}
	case "right", "l":
		if m.FocusedField == FieldTabs {
			if m.ActiveEditorTab == TabHeaders {
				m.ActiveEditorTab = TabBody // Move right
			} else if m.ActiveEditorTab == TabBody {
				m.ActiveEditorTab = TabLoad // Move right
			}
		} else if m.FocusedField == FieldContent && m.ActiveEditorTab == TabHeaders {
			m.FocusedHeaderKey = false
		}
	case "enter":
		if m.FocusedField == FieldTabs {
			// Toggle active tab via cycling?
			m.ActiveEditorTab = (m.ActiveEditorTab + 1) % 3
		} else {
			m.IsEditing = true
			if m.FocusedField == FieldMethod {
				cmd = m.EditorInputs[0].Focus()
			}
			if m.FocusedField == FieldURL {
				cmd = m.EditorInputs[1].Focus()
			}
			if m.FocusedField == FieldContent {
				if m.ActiveEditorTab == TabBody {
					cmd = m.EditorBody.Focus()
				} else if m.ActiveEditorTab == TabHeaders {
					// Focus active header input
					idx := m.FocusedHeaderIdx
					if idx < len(m.HeaderInputs) {
						if m.FocusedHeaderKey {
							cmd = m.HeaderInputs[idx].Key.Focus()
						} else {
							cmd = m.HeaderInputs[idx].Value.Focus()
						}
					}
				} else if m.ActiveEditorTab == TabLoad {
					if m.FocusedHeaderIdx == 0 {
						cmd = m.LoadConfig.Concurrency.Focus()
					} else {
						cmd = m.LoadConfig.Duration.Focus()
					}
				}
			}
		}
	case "n":
		// Add new header logic
		if m.ActiveEditorTab == TabHeaders && m.FocusedField == FieldContent {
			kInput := textinput.New()
			kInput.Placeholder = "Header"
			vInput := textinput.New()
			vInput.Placeholder = "Value"
			m.HeaderInputs = append(m.HeaderInputs, InputPair{Key: kInput, Value: vInput})
			m.FocusedHeaderIdx = len(m.HeaderInputs) - 1
			m.FocusedHeaderKey = true
			m.IsEditing = true
			cmd = m.HeaderInputs[m.FocusedHeaderIdx].Key.Focus()
		}
	case "d":
		// Delete header logic
		if m.ActiveEditorTab == TabHeaders && m.FocusedField == FieldContent {
			if len(m.HeaderInputs) > 0 {
				m.HeaderInputs = append(m.HeaderInputs[:m.FocusedHeaderIdx], m.HeaderInputs[m.FocusedHeaderIdx+1:]...)
				if m.FocusedHeaderIdx >= len(m.HeaderInputs) {
					m.FocusedHeaderIdx = len(m.HeaderInputs) - 1
				}
				if m.FocusedHeaderIdx < 0 {
					m.FocusedHeaderIdx = 0
				}
				// If empty, add one back
				if len(m.HeaderInputs) == 0 {
					kInput := textinput.New()
					kInput.Placeholder = "Header"
					vInput := textinput.New()
					vInput.Placeholder = "Value"
					m.HeaderInputs = append(m.HeaderInputs, InputPair{Key: kInput, Value: vInput})
				}
			}
		}
	}

	return m, cmd
}

// RunRequestCmd executes the current request.
func (m Model) RunRequestCmd() tea.Msg {
	if len(m.Requests) == 0 {
		return nil
	}
	return m.Executor.Execute(m.Requests[m.SelectedReqIdx])
}
