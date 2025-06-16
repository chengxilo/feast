package main

import (
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
		at:    "/",
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
	map_ := make(map[string]tea.Model)
	map_["/"] = menu.NewModel()
	map_["/file"] = file.NewModel()
	map_["/system"] = notyet.Model{}
	map_["/network"] = notyet.Model{}
	map_["/application"] = notyet.Model{}

	p := tea.NewProgram(newModel(map_), tea.WithAltScreen())

	help.New()
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
