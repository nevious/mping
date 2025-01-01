package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/nevious/goping/pinger"
)

/* ---
 * Type and struct definitions
 */
type model struct {
	table *table.Table
	records []DataRecord
}

type tickMsg time.Time

type DataRecord struct {
	Address string
	Sent int
	Failures int
	Loss float64
	LastMessage string
	FailState lipgloss.Style
}

/* ---
 * Method definitions for DataRecord
*/
func (d DataRecord) Refresh() DataRecord {
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

func (d DataRecord) GetRow() []string {
	return []string{
		d.Address,
		fmt.Sprintf("%d", d.Sent),
		fmt.Sprintf("%d", d.Failures),
		fmt.Sprintf("%2.3f%%", d.Loss),
		fmt.Sprintf(d.FailState.Render("⚫")),
		d.LastMessage,
	}
}

func (m model) UpdateRecords() model {
	var rows [][]string

	for index, element := range m.records {
		m.records[index] = element.Refresh()
		rows = append(rows, element.GetRow())
	}
	m.table.ClearRows()
	m.table.Rows(rows...)

	return m
}

/* ---
 * Internals
*/
var (
	renderer = lipgloss.NewRenderer(os.Stdout)
	baseStyle = renderer.NewStyle().Padding(0, 1).Foreground(
		lipgloss.Color("#FFFEEE"),
	)
	headerStyle = baseStyle.Bold(true).Width(5)
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func NewDataRecord(addr string) DataRecord {
	return DataRecord{
		addr, 0, 0.0, 0, "-", lipgloss.NewStyle(), 
	}
}

/* ---
 * Model Methods
*/
func(m model) Init() tea.Cmd {
	return tick()
}

func(m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type) {
		case tickMsg:
			return m.UpdateRecords(), tick()
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

/* ---
 * Main function
*/
func main() {
	records := []DataRecord{
		NewDataRecord("1.1.1.1"),
		NewDataRecord("1.1.1.2"),
		NewDataRecord("185.178.192.107"),
		NewDataRecord("111.1.11.1"),
		NewDataRecord("192.168.50.1"),
		NewDataRecord("192.168.50.22"),
		NewDataRecord("2a00:1450:400a:802::200e"),
		NewDataRecord("google.com"),
		NewDataRecord("nevious.ch"),
		NewDataRecord("gmail.com"),
		NewDataRecord("facebook.com"),
		NewDataRecord("hosttech.ch"),
	}

	s := table.New().Border(
			lipgloss.NormalBorder(),
		).BorderStyle(
			renderer.NewStyle().Foreground(lipgloss.Color("#2196F3")),
		).Headers(
			"Address", "Total", "Failed", "%-Loss", "Up", "Last Message",
		).StyleFunc(func(row, col int) lipgloss.Style {
			var style = baseStyle

			switch {
				case col == 1 && row == -1:
					style = headerStyle.Width(5)
				case col == 2 && row == -1:
					style = headerStyle.Width(5)
				case col == 3 && row == -1:
					style = headerStyle.Width(5)
				case col == 4 && row == -1:
					style = headerStyle.Width(1).AlignHorizontal(lipgloss.Center)
				case col == 5 && row == -1:
					style = headerStyle.Width(200)
				case col == 4:
					style = style.AlignHorizontal(lipgloss.Center)
			}

			return style
		}).BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).BorderColumn(false)

	m := model{s, records}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		slog.Warn("Unable to run programm", err)
	}
}

func _main() {
	//answer, err := pinger.Ping("111.1.11.1")
	answer, err := pinger.Ping("1.1.1.1")
	slog.Info("answer", answer, "error", err)
}
