package comp

import (
	"feast/types"
	"feast/ui/comp/help"
	"feast/ui/logger"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

type MenuItem struct {
	T, D, P string
}

func (i MenuItem) Title() string       { return i.T }
func (i MenuItem) Description() string { return i.D }
func (i MenuItem) FilterValue() string { return i.T }

type SubCategoryModel struct {
	list   list.Model
	help   help.Model
	height int
	width  int
}

func NewSubCategoryModel(items []list.Item) tea.Model {
	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.SetShowHelp(false)
	listModel.SetShowTitle(false)
	listModel.SetShowStatusBar(false)
	helpModel := help.NewHelpModel(help.KeyMap{
		SHelp: []string{"help", "quit"},
		LHelp: [][]string{
			{"enter", "up", "down"},
			{"help", "quit"},
		},
		KeyBindings: map[string]key.Binding{
			"up": key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			"down": key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			"enter": key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("↵enter", "enter"),
			),
			"help": key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("?", "toggle help"),
			),
			"quit": key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
	})
	return SubCategoryModel{listModel, helpModel, 0, 0}
}

func (m SubCategoryModel) Init() tea.Cmd {
	return nil
}

func (m SubCategoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// change view
			path := m.list.Items()[m.list.Index()].(MenuItem).P
			logger.Logger.Debug("menu update, key enter", zap.String("P", path))
			return m, types.RouteCmd(path)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m SubCategoryModel) View() string {
	s := lipgloss.NewStyle().Width(m.width).Height(m.height)
	helpView := m.help.View()
	helpViewHeight := lipgloss.Height(helpView)
	m.list.SetSize(m.width, m.height-helpViewHeight)
	listView := m.list.View()
	return s.Render(lipgloss.JoinVertical(lipgloss.Left, listView, helpView))
}
