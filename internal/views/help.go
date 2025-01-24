package views

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nevious/mping/internal/utils"
)

const help_text = `
Mping Help
============================================

Root View:
	q: Quit				?: Show this Help

Help View:
	q: Quit				<esc> : Return to Root View


Main options
	-a [addr,...]		List of addresses to ping
	--version 			Show Version and exit
	--help				Show usage and exit
`

type helpModel struct {
	text string
	rootModel *rootModel
}

func (o helpModel) View() string {
	return textStyle.Render(o.text)
}

func (o helpModel) Init() tea.Cmd { return nil }

func (o helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "q":
					return o, tea.Quit
				case "esc":
					return o.rootModel, utils.SecondTick()
			}
	}
	return o, nil
}

func NewHelp(r *rootModel) *helpModel {
	return &helpModel{help_text, r}
}
