package views

import (
	"fmt"
	"strings"
	"time"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nevious/mping/internal/utils"
	"github.com/nevious/mping/internal/pinger"
)

type traceModel struct {
	hops []hopRecord
	rootModel *rootModel
	dst string
}

type hopRecord struct {
	Address string
	Duration time.Duration
}

func (t *traceModel) runTrace() traceModel {
	if len(t.hops) > 0 {
		return *t
	}
	destination, err := net.ResolveIPAddr("ip4", t.dst)
	if err != nil {
		panic("Unable to resolve address")
	}

	for ttl := 1; ttl < 32 ; ttl++ {
		resp, err := pinger.SendICMPEcho(destination.String(), ttl)
		var hop hopRecord

		if err !=  nil {
			hop = hopRecord{Address: "*", Duration: resp.Duration}
		} else {
			hop = hopRecord{Address: resp.Peer.String(), Duration: resp.Duration}
		}

		t.hops = append(t.hops, hop)

		if hop.Address == destination.String() {
			return *t
		}
	}
	return *t
}

func (t *traceModel) SetDestination(dst string) string {
	t.dst = dst
	return dst
}

func (t traceModel) View() string {
	if t.dst == "" {
		//  should never happend
		return fmt.Sprint("No trace destination selected")
	}
	
	var lines = []string{fmt.Sprintf("Trace to %s", t.dst), "==========================================="}
	for index, element := range t.hops {
		lines = append(lines, fmt.Sprintf("%02d\t%-20s\t%-20v", index, element.Address, element.Duration))
	}

	return textStyle.Render(strings.Join(lines, "\n"))

}

func (t traceModel) Init() tea.Cmd { return nil }

func (t traceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "q":
					return t, tea.Quit
				case "r":
					t.hops = nil
					return t, nil
				case "R":
					t.rootModel.table.ClearRows()
					var x []string
					for _, element := range t.hops {
						if element.Address == "*" {
							continue
						}

						x = append(x, element.Address)
					}
					records := makeTableRows(x)
					t.rootModel.records = records
					return *t.rootModel, nil
				case "esc":
					return *t.rootModel, nil
			}
		case utils.SecondTickMsg:
			return t.runTrace(), utils.SecondTick()
	}

	return t, nil
}

func NewTrace(r *rootModel) *traceModel {
	return &traceModel{rootModel: r}
}
