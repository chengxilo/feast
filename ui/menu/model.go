package menu

import (
	"feast/types"
	"feast/ui/component/help"
	"feast/ui/logger"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

var log = logger.Logger.With(zap.String("model", "menu"))

type item struct {
	title, desc, path string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	list   list.Model
	help   help.Model
	height int
	width  int
}

func NewModel() tea.Model {
	var items = []list.Item{
		item{title: "System", desc: "get system info", path: "/system"},
		item{title: "File", desc: "file explorer", path: "/file"},
		item{title: "Fire Wall", desc: "fire wall information", path: "/firewall"},
		item{title: "Application", desc: "application management", path: "/application"},
	}
	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.SetShowHelp(false)
	listModel.SetShowTitle(false)
	listModel.SetShowStatusBar(false)
	helpModel := help.NewHelpModel(help.KeyMap{
		SHelp: []string{"help", "quit"},
		LHelp: [][]string{
			{"up", "down", "left", "right"},
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
			"left": key.NewBinding(
				key.WithKeys("left", "h"),
				key.WithHelp("←/h", "move left"),
			),
			"right": key.NewBinding(
				key.WithKeys("right", "l"),
				key.WithHelp("→/l", "move right"),
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
	return Model{
		listModel, helpModel, 0, 0}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			path := m.list.Items()[m.list.Index()].(item).path
			logger.Logger.Debug("menu update, key enter", zap.String("path", path))
			return m, types.RouteCmd(path)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := lipgloss.NewStyle().Width(m.width).Height(m.height)
	helpView := m.help.View()
	helpViewHeight := lipgloss.Height(helpView)
	m.list.SetSize(m.width, m.height-helpViewHeight)
	listView := m.list.View()
	return s.Render(lipgloss.JoinVertical(lipgloss.Left, listView, helpView))
}
