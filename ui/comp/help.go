package comp

import (
	"feast/logger"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

var log = logger.Logger.With(zap.String("at", "comp/help"))

// KeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type KeyMap struct {
	KeyBindings map[string]key.Binding
	SHelp       []string
	LHelp       [][]string
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	shortHelp := make([]key.Binding, len(k.SHelp))
	for i, v := range k.SHelp {
		shortHelp[i] = k.KeyBindings[v]
	}
	return shortHelp
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	longHelp := make([][]key.Binding, len(k.LHelp))
	for i, v := range k.LHelp {
		sub := make([]key.Binding, len(v))
		for j, b := range v {
			sub[j] = k.KeyBindings[b]
		}
		longHelp[i] = sub
	}
	return longHelp
}

type HelpModel struct {
	keys     KeyMap
	help     help.Model
	quitting bool
}

func NewHelpModel(keys KeyMap) HelpModel {
	return HelpModel{
		keys: keys,
		help: help.New(),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) SetKeyMap(k KeyMap) {
	m.keys = k
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.KeyBindings["help"]) {
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	m.help, cmd = m.help.Update(msg)
	return m, cmd
}

func (m HelpModel) View() string {
	return m.help.View(m.keys)
}
