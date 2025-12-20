package tui

import (
	"lazycurl/internal/load"

	tea "github.com/charmbracelet/bubbletea"
)

// WaitForStats produces a command that waits for a message on the channel.
func WaitForStats(ch chan load.StatsMsg) tea.Cmd {
	return func() tea.Msg {
		if msg, ok := <-ch; ok {
			return msg
		}
		return nil
	}
}
