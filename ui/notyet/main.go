package notyet

import (
	"feast/types"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct{}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, types.RouteCmd("/")
	}
	return m, nil
}

func (m Model) View() string {
	s := fmt.Sprint("Building, press any key to go exit...")
	return s + "\n"
}
