package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/nevious/mping/internal/pinger"
	"github.com/nevious/mping/internal/utils"
)


// General rootModel for bubbletea to interact with
type rootModel struct {
	table *table.Table
	records []dataRecord
}


// Model for a single row within the table.
// added to bubbletea rootModel via records-array
type dataRecord struct {
	Address string
	Sent int
	Failures int
	Loss float64
	PingReturn pinger.IcmpReply
	LastMessage string
	FailState lipgloss.Style
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
		m.records[index] = element.refresh()
		rows = append(rows, element.render())
	}
	m.table.ClearRows()
	m.table.Rows(rows...)

	return m
}

// refresh a single dataRecord within m.records
func (d dataRecord) refresh() dataRecord {
	answer, err := pinger.SendICMPEcho(d.Address, 64)
	if err != nil {
		d.PingReturn = *answer
		d.LastMessage = fmt.Sprintf("%+v", err)
		d.Failures = d.Failures + 1
		d.FailState = d.FailState.Foreground(lipgloss.Color("#F18C8E"))
	} else {
		d.PingReturn = *answer
		d.LastMessage = fmt.Sprintf(
			"Peer: %v - Checksum: %d - Proto: %d - Duration: %v",
			answer.Peer, answer.Checksum, answer.IcmpProto, answer.Duration,
		)
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
		result = append(result, dataRecord{element, 0, 0.0, 0, pinger.IcmpReply{}, "-", lipgloss.NewStyle()})
	}

	return result
}

func MakeTable(records []string) rootModel {
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
