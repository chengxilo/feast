package root

import (
	_const "feast/const"
	"feast/logger"
	"feast/types"
	"feast/ui"
	"feast/ui/comp"
	"feast/ui/file"
	"feast/ui/notyet"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"go.uber.org/zap"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Gray         = lipgloss.Color("240")
	FocusClr     = lipgloss.Color("15")
	sidebarWidth = 9
)

var (
	sideBarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, true, true).
			BorderForeground(Gray).
			Width(sidebarWidth)
	sideBarFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(FocusClr).
				Width(sidebarWidth)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).PaddingRight(1).
				Background(lipgloss.Color("62"))
	contentStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, true, true, false).
			BorderForeground(Gray)
	contentFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(FocusClr)
)

type item struct {
	route string
	title string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	itm, ok := listItem.(item)
	if !ok {
		return
	}
	str := itm.title

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, ""))
		}
	}

	fmt.Fprint(w, fn(str))
}

type UI struct {
	route          map[string]ui.Model
	at             string
	list           list.Model
	sideBarFocused bool
	help           comp.HelpModel
	choice         string
	height, width  int
}

func NewModel() *UI {
	mp := make(map[string]ui.Model)
	mp[_const.RouteHome] = notyet.NewNotYet()
	mp[_const.RouteFile] = file.NewModel()
	mp[_const.RouteSystem] = notyet.NewNotYet()
	mp[_const.RouteNetWork] = notyet.NewNotYet()
	mp[_const.RouteApplication] = notyet.NewNotYet()
	items := []list.Item{
		item{_const.RouteNetWork, _const.TitleNetWork},
		item{_const.RouteSystem, _const.TitleSystem},
		item{_const.RouteFile, _const.TitleFile},
	}

	l := list.New(items, itemDelegate{}, sidebarWidth, 10)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	helpModel := comp.NewHelpModel(comp.KeyMap{
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
	return &UI{list: l, route: mp, at: _const.RouteHome, sideBarFocused: true, help: helpModel}
}

func (m *UI) Init() tea.Cmd {
	return nil
}

func (m *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	logger.Logger.Debug("root update", zap.Any("msg", msg))
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		sideBarStyle = sideBarStyle.Height(m.height - 4)
		return m, nil

	case types.RouteMsg:
		// update where we are
		m.at = msg.Path
		// update its window size
		var mdl tea.Model
		mdl, cmd = m.route[m.at].Update(msg)
		m.route[m.at] = mdl.(ui.Model)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		keypress := msg.String()
		if keypress == "ctrl+c" {
			return m, tea.Quit
		}
		if keypress == "tab" {
			if m.sideBarFocused {
				m.sideBarFocused = false
				m.route[m.at].Focus()
			} else {
				m.sideBarFocused = true
				m.route[m.at].Blur()
			}
		}
		switch keypress {
		case "enter":
			if m.sideBarFocused {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					cmds = append(cmds, types.RouteCmd(i.route))
				}
			}
		}
	}

	// Only update m.list when it is not key message
	if _, ok := msg.(tea.KeyMsg); !ok || m.sideBarFocused {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	// run the update method of where we are
	var mdl tea.Model
	mdl, cmd = m.route[m.at].Update(msg)
	m.route[m.at] = mdl.(ui.Model)
	cmds = append(cmds, cmd)

	mdl, cmd = m.help.Update(msg)
	m.help = mdl.(comp.HelpModel)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *UI) View() string {
	helpView := m.help.View()

	helpViewHeight := lipgloss.Height(helpView)
	m.list.SetHeight(m.height - 2 - helpViewHeight)
	var sideBar string
	if m.sideBarFocused {
		sideBar = sideBarFocusedStyle.Render(m.list.View())
	} else {
		sideBar = sideBarStyle.Render(m.list.View())
	}

	m.route[m.at].SetHeight(m.height - helpViewHeight - 2)
	m.route[m.at].SetWidth(m.width - lipgloss.Width(sideBar) - 2)

	var contentView string
	if m.route[m.at].IsFocused() {
		contentView = contentFocusedStyle.Render(m.route[m.at].View())
	} else {
		contentView = contentStyle.Render(m.route[m.at].View())
	}
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sideBar, contentView)

	return lipgloss.JoinVertical(lipgloss.Left, mainView, helpView)
}
