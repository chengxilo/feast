package main

import (
	_const "feast/const"
	"feast/types"
	"feast/ui/file"
	"feast/ui/logger"
	"feast/ui/menu"
	"feast/ui/notyet"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"os"
)

type model struct {
	route               map[string]tea.Model
	at                  string
	latestWindowSizeMsg tea.WindowSizeMsg
}

func newModel(map_ map[string]tea.Model) *model {
	return &model{
		route: map_,
		at:    _const.RouteHome,
	}
}
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	logger.Logger.Debug("main update", zap.String("types", fmt.Sprint(msg)), zap.String("at", m.at))
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.latestWindowSizeMsg = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case types.RouteMsg:
		// update where we are
		m.at = msg.Path
		logger.Logger.Debug("route update", zap.String("path", m.at), zap.Any("latest window size msg", m.latestWindowSizeMsg))
		// update its window size
		m.route[m.at], cmd = m.route[m.at].Update(m.latestWindowSizeMsg)
		cmds = append(cmds, cmd)
	}

	// run the update method of where we are
	m.route[m.at], cmd = m.route[m.at].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.route[m.at].View()
}

func main() {
	mp := make(map[string]tea.Model)
	mp[_const.RouteHome] = menu.NewModel()
	mp[_const.RouteFile] = file.NewModel()
	mp[_const.RouteSystem] = notyet.Model{}
	mp[_const.RouteNetWork] = notyet.Model{}
	mp[_const.RouteApplication] = notyet.Model{}

	p := tea.NewProgram(newModel(mp), tea.WithAltScreen())

	help.New()
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
