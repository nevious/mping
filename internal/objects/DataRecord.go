package objects

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/nevious/mping/internal/pinger"
	"github.com/nevious/mping/internal/utils"
)

// Structure for a single row within the table.
// added to bubbletea rootModel via records-array
type DataRecord struct {
	Address string
	Sent int
	Failures int
	Loss float64
	PingReturn utils.IcmpReply
	LastMessage string
	FailState lipgloss.Style
}

// refresh a single dataRecord within m.records
func (d DataRecord) Refresh() DataRecord {
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

func (d DataRecord) Render() []string {
	return []string{
		d.Address,
		fmt.Sprintf("%d", d.Sent),
		fmt.Sprintf("%d", d.Failures),
		fmt.Sprintf("%2.3f%%", d.Loss),
		fmt.Sprintf(d.FailState.Render("⚫")),
		d.LastMessage,
	}
}


func MakeTableRows(addrs []string) []DataRecord {
	var result []DataRecord
	for _, element := range addrs {
		result = append(result, DataRecord{element, 0, 0.0, 0, utils.IcmpReply{}, "-", lipgloss.NewStyle()})
	}

	return result
}
