package types

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type RouteMsg struct {
	Path string
}

func (m RouteMsg) String() string {
	return fmt.Sprintf("route: %s", m.Path)
}

func RouteCmd(path string) tea.Cmd {
	return func() tea.Msg {
		return RouteMsg{
			Path: path,
		}
	}
}

type ErrMsg struct {
	Err error
}

func (m ErrMsg) String() string {
	return fmt.Sprintf("error: %v", m.Err)
}

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrMsg{Err: err}
	}
}
