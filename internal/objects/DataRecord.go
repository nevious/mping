package objects

import (
	"fmt"
	"math"
	"time"
	"errors"
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
	Last time.Duration
	Max time.Duration
	Min time.Duration
	Average time.Duration
	StD time.Duration
	durations []time.Duration
}

// Calculate the average latency
func (d *DataRecord) calc_average() time.Duration {
	sum := time.Duration(0)
	for i := 0 ; i <= len(d.durations)-1; i++ {
		sum += d.durations[i]
	}

	avg := float64(sum) / float64(len(d.durations))
	return time.Duration(avg)
}

// calculate standard deviation
// square root of [ (1/N) times Sigma i=1 to N of (xi - mu)^2 ]
func (d *DataRecord) calc_std() (time.Duration, error) {
	if d.Average == time.Duration(0) {
		return time.Duration(0), errors.New("cannot calc stdev with average of 0 or nil")
	}

	avg_f := float64(d.Average)
	n := float64(len(d.durations))
	var sigma_result float64

	for _, element := range d.durations {
		value := float64(element)
		sigma_result = sigma_result + math.Pow(value - avg_f, 2)
	}

	result := math.Sqrt((1/n)*sigma_result)
	return time.Duration(result*float64(time.Nanosecond)), nil
}

// refresh a single dataRecord within m.records
func (d *DataRecord) Refresh() *DataRecord {
	answer, err := pinger.SendICMPEcho(d.Address, 64)
	if err != nil {
		d.PingReturn = *answer
		d.LastMessage = fmt.Sprintf("%+v", err)
		d.Failures = d.Failures + 1
		d.FailState = d.FailState.Foreground(lipgloss.Color("#F18C8E"))
	} else {
		d.PingReturn = *answer
		d.LastMessage = fmt.Sprintf(
			"Peer: %-25v Check: %5d \t Proto: %2d",
			answer.Peer, answer.Checksum, answer.IcmpProto,
		)
		d.FailState = d.FailState.Foreground(lipgloss.Color("#4CAEA3"))
		d.Last = answer.Duration

		switch {
			case answer.Duration > d.Max:
				d.Max = answer.Duration
			case answer.Duration < d.Min:
				d.Min = answer.Duration
		}

		d.durations = append(d.durations, answer.Duration)

		d.Average = d.calc_average()
		d.StD, _ = d.calc_std()
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
		fmt.Sprintf(d.FailState.Render("⚫")),
		d.LastMessage,
		fmt.Sprintf("%2.3f%%", d.Loss),
		fmt.Sprintf("%s", d.Last),
		fmt.Sprintf("%s", d.Max),
		fmt.Sprintf("%s", d.Min),
		fmt.Sprintf("%s", d.Average),
		fmt.Sprintf("%s", d.StD),
	}
}

func MakeTableRows(addrs []string) []DataRecord {
	var result []DataRecord
	for _, element := range addrs {
		result = append(result, DataRecord{
			element, 0, 0.0, 0, utils.IcmpReply{}, "-", lipgloss.NewStyle(), 0, time.Nanosecond, time.Second, 0, 0, nil},
		)
	}

	return result
}
