package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
)

var (
	// Styles
	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")) // Purple

	blurredStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")) // Grey

	docStyle = lipgloss.NewStyle().Margin(1, 2)

	// Component Styles
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	labelStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginBottom(0)
	activeLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true).MarginBottom(0)
)

// View renders the TUI.
func (m Model) View() string {
	// Calculate pane widths
	// MVP: Equal split for 3 panes approx
	paneWidth := (m.Width - 6) / 3
	paneHeight := m.Height - 8 // Reserve space for help

	if paneWidth < 10 {
		paneWidth = 10
	}
	if paneHeight < 5 {
		paneHeight = 5
	}

	requestsView := m.viewRequests(paneWidth, paneHeight)
	editorView := m.viewEditor(paneWidth, paneHeight)
	responseView := m.viewResponse(paneWidth, paneHeight)

	// Join Main Panes
	mainPanes := lipgloss.JoinHorizontal(lipgloss.Top, requestsView, editorView, responseView)

	// Help Bar
	helpView := m.Help.View(m.KeyMap)

	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, mainPanes, helpView))
}

func (m Model) viewRequests(width, height int) string {
	style := blurredStyle
	if m.ActivePane == PaneRequests {
		style = focusedStyle
	}

	var items []string
	for i, req := range m.Requests {
		line := fmt.Sprintf("%s %s", req.Method, req.URL)
		if len(line) > width-4 {
			line = line[:width-4] + "..."
		}

		if i == m.SelectedReqIdx {
			items = append(items, selectedItemStyle.Render("> "+line))
		} else {
			items = append(items, itemStyle.Render("  "+line))
		}
	}

	content := strings.Join(items, "\n")
	if content == "" {
		content = "No requests. Press 'n' to create."
	}

	return style.
		Width(width).
		Height(height).
		Render(content)
}

func (m Model) viewEditor(width, height int) string {
	style := blurredStyle
	if m.ActivePane == PaneEditor {
		style = focusedStyle
	}

	// Helper to render a field with label
	renderField := func(field EditorField, label string, content string) string {
		lStyle := labelStyle
		if m.ActivePane == PaneEditor && m.FocusedField == field {
			lStyle = activeLabelStyle
		}
		return lStyle.Render(label) + "\n" + content + "\n"
	}

	methodView := renderField(FieldMethod, "Method", m.EditorInputs[0].View())
	urlView := renderField(FieldURL, "URL", m.EditorInputs[1].View())

	// Tabs View
	tabHeader := "Headers"
	tabBody := "Body"
	tabLoad := "Load"

	if m.ActiveEditorTab == TabHeaders {
		tabHeader = "[" + tabHeader + "]"
	} else if m.ActiveEditorTab == TabBody {
		tabBody = "[" + tabBody + "]"
	} else if m.ActiveEditorTab == TabLoad {
		tabLoad = "[" + tabLoad + "]"
	}
	tabsContent := fmt.Sprintf("%s  %s  %s", tabHeader, tabBody, tabLoad)
	tabsView := renderField(FieldTabs, "", tabsContent) // Empty label for tabs

	// Content View
	var contentView string
	if m.ActiveEditorTab == TabBody {
		contentView = renderField(FieldContent, "Body Content", m.EditorBody.View())
	} else if m.ActiveEditorTab == TabHeaders {
		// Headers List
		var sb strings.Builder
		sb.WriteString(labelStyle.Render("Headers List (n: new, d: del)") + "\n")

		for i, pair := range m.HeaderInputs {
			// Determine styles for Key vs Value
			kStyle := lipgloss.NewStyle()
			vStyle := lipgloss.NewStyle()

			if m.ActivePane == PaneEditor && m.FocusedField == FieldContent && i == m.FocusedHeaderIdx {
				if m.FocusedHeaderKey {
					kStyle = activeLabelStyle // Highlight Key
				} else {
					vStyle = activeLabelStyle // Highlight Value
				}
				sb.WriteString("> ")
			} else {
				sb.WriteString("  ")
			}

			// Adjust width of inputs slightly for the list
			pair.Key.Width = 15
			pair.Value.Width = 20

			sb.WriteString(kStyle.Render(pair.Key.View()))
			sb.WriteString(" : ")
			sb.WriteString(vStyle.Render(pair.Value.View()))
			sb.WriteString("\n")
		}
		contentView = sb.String()
	} else if m.ActiveEditorTab == TabLoad {
		// Load Config
		cStyle := labelStyle
		dStyle := labelStyle
		// If focused on Content, determine which sub-input is focused based on index
		if m.FocusedField == FieldContent {
			if m.FocusedHeaderIdx == 0 {
				cStyle = activeLabelStyle
			}
			if m.FocusedHeaderIdx == 1 {
				dStyle = activeLabelStyle
			}
		}

		contentView = fmt.Sprintf("%s\n%s\n\n%s\n%s",
			cStyle.Render("Concurrency"), m.LoadConfig.Concurrency.View(),
			dStyle.Render("Duration"), m.LoadConfig.Duration.View())
	}

	content := methodView + "\n" + urlView + "\n" + tabsView + "\n" + contentView

	return style.
		Width(width).
		Height(height).
		Render(content)
}

func (m Model) viewResponse(width, height int) string {
	style := blurredStyle
	if m.ActivePane == PaneResponse {
		style = focusedStyle
	}

	// Check if Load Mode (either running, or we have stats to show)
	showDashboard := m.LoadState.IsRunning
	if !showDashboard && m.LoadState.Stats != nil && m.LoadState.Stats.TotalRequests > 0 {
		showDashboard = true
	}

	var content string
	if showDashboard {
		content = m.viewDashboard(width, height)
	} else {
		content = "No response yet.\nPress 'r' to run."
		if m.Response != nil {
			if m.Response.Error != nil {
				content = fmt.Sprintf("Error:\n%v", m.Response.Error)
			} else {
				// Basic truncating for large bodies logic could go here
				body := m.Response.Body
				if len(body) > 2000 {
					body = body[:2000] + "\n... (truncated)"
				}

				content = fmt.Sprintf("Status: %d\nTime: %s\n\nBody:\n%s",
					m.Response.StatusCode, m.Response.TimeTaken, body)
			}
		}
	}

	return style.
		Width(width).
		Height(height).
		Render(content)
}

func (m Model) viewDashboard(width, height int) string {
	s := m.LoadState.Stats
	if s == nil {
		return "Initializing..."
	}

	// 1. Stats Column
	stats := fmt.Sprintf("Valid Reqs: %d\nElapsed: %s\nAvg Latency: %s\nMax Latency: %s\n\nStatus Codes:\n",
		s.TotalRequests, s.ElapsedTime.Round(time.Millisecond), s.AvgLatency.Round(time.Millisecond), s.MaxLatency.Round(time.Millisecond))

	for code, count := range s.StatusCodes {
		pct := 0.0
		if s.TotalRequests > 0 {
			pct = (float64(count) / float64(s.TotalRequests)) * 100
		}
		stats += fmt.Sprintf("  %d: %d (%.1f%%)\n", code, count, pct)
	}

	// 2. Graph
	// Prepare data
	data := s.LatencyPoints
	// Cap points to fit screen width roughly? Asciigraph handles auto-scaling but too many points might be slow/messy.
	// If > width, take tail
	plotData := data
	graphWidth := width - 30
	if graphWidth < 0 {
		graphWidth = 10
	}

	if len(data) > graphWidth {
		plotData = data[len(data)-graphWidth:]
	}
	if len(plotData) == 0 {
		plotData = []float64{0}
	}

	graph := asciigraph.Plot(plotData,
		asciigraph.Height(height/2),
		asciigraph.Width(graphWidth),
		asciigraph.Caption("Latency (ms)"),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, graph, "\n\n"+stats)
}
