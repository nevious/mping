package utils

import (
	"time"
	tea "github.com/charmbracelet/bubbletea"
)

// basically emit a tea command every time.Second
type SecondTickMsg time.Time
func SecondTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return SecondTickMsg(t)
	})
}
