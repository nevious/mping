package views

import (
	"os"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styling
	renderer = lipgloss.NewRenderer(os.Stdout)
	baseStyle = renderer.NewStyle().Padding(0, 1).Foreground(
		lipgloss.Color("#FFFEEE"),
	)
	headerStyle = baseStyle.Bold(true).Width(5)

	// views
	helpView helpModel
	traceView traceModel

	// table index
	record_index int = 1

	textStyle = lipgloss.NewStyle().MarginLeft(2).MarginRight(2).MarginTop(1)
)
