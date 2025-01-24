package views

import (
	"os"
	"github.com/charmbracelet/lipgloss"
)

var (
	renderer = lipgloss.NewRenderer(os.Stdout)

	baseStyle = renderer.NewStyle().Padding(0, 1).Foreground(
		lipgloss.Color("#FFFEEE"),
	)
	timestampStyle = baseStyle.MarginTop(1).MarginBottom(1).Bold(true)
	highlightStyle = baseStyle.Foreground(lipgloss.Color("#111111")).Background(lipgloss.Color("#a9d1f2"))

	textStyle = lipgloss.NewStyle().MarginLeft(2).MarginRight(2).MarginTop(1)
)
