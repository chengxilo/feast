package menu

import (
	_const "feast/const"
	"feast/types"
	"feast/ui/comp"
	"feast/ui/logger"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"go.uber.org/zap"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var log = logger.Logger.With(zap.String("model", "menu"))

const (
	sidebarWidth = 9
)

var (
	lisStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, true, true).
			Width(sidebarWidth)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Background(lipgloss.Color("62"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
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
			return selectedItemStyle.Render(strings.Join(s, "useless"))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	list          list.Model
	help          comp.Model
	choice        string
	height, width int
}

func NewModel() *Model {
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
	l.Styles.PaginationStyle = paginationStyle
	return &Model{list: l}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		lisStyle = lisStyle.Height(m.height - 4)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":
			return m, types.RouteCmd(_const.RouteHome)

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			log.Debug("enter key", zap.String("route", i.route))
			if ok {
				return m, types.RouteCmd(i.route)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	helpView := comp.NewHelpModel(comp.KeyMap{
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
	}).View()

	helpViewHeight := lipgloss.Height(helpView)
	m.list.SetHeight(m.height - 2 - helpViewHeight)
	styledList := lisStyle.Render(m.list.View())
	sideBar := lipgloss.JoinVertical(lipgloss.Top, "┌─[FEAST]─┐", styledList)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sideBar)

	return lipgloss.JoinVertical(lipgloss.Left, mainView, helpView)
}
