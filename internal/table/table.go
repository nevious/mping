package table

import (
	"time"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/nevious/mping/internal/pinger"
)

var (
	renderer = lipgloss.NewRenderer(os.Stdout)
	baseStyle = renderer.NewStyle().Padding(0, 1).Foreground(
		lipgloss.Color("#FFFEEE"),
	)
	headerStyle = baseStyle.Bold(true).Width(5)
)

// General model for bubbletea to interact with
type model struct {
	table *table.Table
	records []dataRecord
}

// Model for a single row within the table.
// added to bubbletea model via records-array
type dataRecord struct {
	Address string
	Sent int
	Failures int
	Loss float64
	LastMessage string
	FailState lipgloss.Style
}

// basically emit a tea command every time.Second
type tickMsg time.Time
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// initialize the model with a tick
func (m model) Init() tea.Cmd {
	return tick()
}

// any sort of event shoudl trigger this method
func(m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type) {
		case tickMsg:
			return m.updateRecords(), tick()
		case tea.WindowSizeMsg:
			m.table = m.table.Width(msg.Width)
		case tea.KeyMsg:
			switch msg.String() {
				case "q", "ctrl+c":
					return m, tea.Quit
			}
	}

	return m, nil
}

func(m model) View() string {
	return m.table.String()
}

// triggered by Update() if the tick ran out
// replace rows with new values
func (m model) updateRecords() model {
	var rows [][]string

	for index, element := range m.records {
		m.records[index] = element.refresh()
		rows = append(rows, element.render())
	}
	m.table.ClearRows()
	m.table.Rows(rows...)

	return m
}

// refresh a single dataRecord within m.records
func (d dataRecord) refresh() dataRecord {
	answer, err := pinger.Ping(d.Address)
	if err != nil {
		d.LastMessage = fmt.Sprintf("%+v", err)
		d.Failures = d.Failures + 1
		d.FailState = d.FailState.Foreground(lipgloss.Color("#F18C8E"))
	} else {
		d.LastMessage = answer
		d.FailState = d.FailState.Foreground(lipgloss.Color("#4CAEA3"))
	}

	d.Sent = d.Sent + 1
	d.Loss = float64(d.Failures)/float64(d.Sent)*100.0
	return d
}

func (d dataRecord) render() []string {
	return []string{
		d.Address,
		fmt.Sprintf("%d", d.Sent),
		fmt.Sprintf("%d", d.Failures),
		fmt.Sprintf("%2.3f%%", d.Loss),
		fmt.Sprintf(d.FailState.Render("⚫")),
		d.LastMessage,
	}
}


func makeTableRows(addrs []string) []dataRecord {
	var result []dataRecord
	for _, element := range addrs {
		result = append(result, dataRecord{element, 0, 0.0, 0, "-", lipgloss.NewStyle()})
	}

	return result
}

func MakeTable(records []string) model {
	rows := makeTableRows(records)
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

			return style
		}).BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).BorderColumn(false)

	return model{s, rows}
}

func LaunchTablePing(m model) (tea.Model, error) {
	if programm, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return nil, err
	} else {
		return programm, nil
	}
}
