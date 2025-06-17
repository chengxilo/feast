package notyet

import (
	"feast/types"
	"feast/ui"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	focused bool
	height  int
	width   int
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, types.RouteCmd("/")
	}
	return m, nil
}

func (m *Model) View() string {
	s := fmt.Sprint("Building, press any key to go exit...")
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(s)
}

func (m *Model) SetWidth(width int)   { m.width = width }
func (m *Model) SetHeight(height int) { m.height = height }
func (m *Model) Focus() {
	m.focused = true
}
func (m *Model) Blur() {
	m.focused = false
}
func (m *Model) IsFocused() bool {
	return m.focused
}

func NewNotYet() ui.Model {
	mdl := &Model{focused: false}
	return mdl
}
