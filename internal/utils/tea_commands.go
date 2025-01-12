package utils

import (
	"time"
	tea "github.com/charmbracelet/bubbletea"
)

// Primitive type and function to emit this type
// every second
type SecondTickMsg time.Time
func SecondTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return SecondTickMsg(t)
	})
}
