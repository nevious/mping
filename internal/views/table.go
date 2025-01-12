package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/nevious/mping/internal/utils"
	"github.com/nevious/mping/internal/objects"
)


// General rootModel for bubbletea to interact with
type rootModel struct {
	table *table.Table
	records []objects.DataRecord
}

// initialize the rootModel with a tick
func (m rootModel) Init() tea.Cmd {
	helpView = *NewHelp(&m)
	traceView = *NewTrace(&m)

	return utils.SecondTick()
}

// any sort of event shoudl trigger this method
func(m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type) {
		case utils.SecondTickMsg:
			return m.updateRecords(), utils.SecondTick()
		case tea.WindowSizeMsg:
			m.table = m.table.Width(msg.Width)
		case tea.KeyMsg:
			switch msg.String() {
				case "q", "ctrl+c":
					return m, tea.Quit
				case "?":
					return helpView, nil
				case "j":
					if record_index == len(m.records) -1 {
						record_index = 0
					} else {
						record_index = record_index+1
					}
					return m, nil
				case "k":
					if record_index == 0 {
						record_index = len(m.records) - 1
					} else {
						record_index = record_index-1
					}
					return m, nil
				case "t":
					traceView.SetDestination(m.records[record_index].Address)
					return traceView, nil
			}
	}

	return m, nil
}

func(m rootModel) View() string {
	return m.table.String()
}

// triggered by Update() if the tick ran out
// replace rows with new values
func (m rootModel) updateRecords() rootModel {
	var rows [][]string

	for index, element := range m.records {
		m.records[index] = element.Refresh()
		rows = append(rows, element.Render())
	}
	m.table.ClearRows()
	m.table.Rows(rows...)

	return m
}
func MakeTable(records []string) rootModel {
	rows := objects.MakeTableRows(records)
	s := table.New().Border(
			lipgloss.NormalBorder(),
		).BorderStyle(
			renderer.NewStyle().Foreground(lipgloss.Color("#2196F3")),
		).Headers(
			"Address", "Total", "Failed", "%-Loss", "Up", "Last Message",
		).StyleFunc(func(row, col int) lipgloss.Style {
			var style = baseStyle

			switch {
				case col == 0 && row == -1:
					style = headerStyle
				case col == 1 && row == -1:
					style = headerStyle
				case col == 2 && row == -1:
					style = headerStyle
				case col == 3 && row == -1:
					style = headerStyle
				case col == 4 && row == -1:
					style = headerStyle.Width(5).AlignHorizontal(lipgloss.Center)
				case col == 5 && row == -1:
					style = headerStyle.Width(200)
				case col == 4:
					style = style.AlignHorizontal(lipgloss.Center)
			}

			if row == record_index {
				style = style.Foreground(lipgloss.Color("#111111")).Background(lipgloss.Color("#2196f3"))
			}

			return style
		}).BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).BorderColumn(false)

	return rootModel{s, rows}
}

func LaunchTablePing(m rootModel) (tea.Model, error) {
	if programm, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return nil, err
	} else {
		return programm, nil
	}
}
