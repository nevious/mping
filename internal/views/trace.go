package views

import (
	"net"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"

	"github.com/nevious/mping/internal/utils"
)

type traceModel struct {
	hops []net.IPAddr
	rootModel *rootModel
}

func (t traceModel) View() string {
	return "Nothing yet"
}

func (t traceModel) Init() tea.Cmd { return nil }

func (t traceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "q":
					return t, tea.Quit
				case "esc":
					return *t.rootModel, utils.SecondTick()
			}
	}
	return t, nil
}

func NewTrace(r *rootModel) *traceModel {
	return &traceModel{rootModel: r}
}
