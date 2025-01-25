package views

import (
	"fmt"
	"time"
	"os/user"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/nevious/mping/internal/utils"
	"github.com/nevious/mping/internal/objects"
)

var record_index int = 1
var helpView *helpModel
var traceView *traceModel

// General rootModel for bubbletea to interact with
type rootModel struct {
	table *table.Table
	records []objects.DataRecord
}

// initialize the rootModel with a tick
func (m rootModel) Init() tea.Cmd {
	helpView = NewHelp(&m)
	traceView = NewTrace(&m, helpView)

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
	username, err := user.Current()
	if err != nil {
		username = &user.User{}
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	
	v_header := timestampStyle.Render(
		fmt.Sprintf(
			"%s@%s\t%v\tPID: %d", username.Username, hostname,
			time.Now().Format(time.RFC850), os.Getpid(),
		),
	)
	
	output := fmt.Sprintf("%v\n%v", v_header, m.table.String())
	return output
}

// called by Update() if the tick ran out
// replace rows with new values
func (m *rootModel) updateRecords() *rootModel {
	// [row][cells] - if that makes sense
	var rows [][]string

	for index, element := range m.records {
		m.records[index] = *element.Refresh()
		// Rows must be []string or for the lipgloss.table
		// object to be able to handle it. Therefore the replace logic
		rows = append(rows, m.records[index].Render())
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
			"Address", "Total", "Failed", "Up", "Last Message",
			"%-Loss", "Last", "Max", "Min", "Avg", "StD",
		).StyleFunc(func(row, col int) lipgloss.Style {
			var style = baseStyle

			switch {
				case col == 4 && row == -1:
					style = baseStyle.Width(200)
				case row == record_index:
					style = highlightStyle
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
