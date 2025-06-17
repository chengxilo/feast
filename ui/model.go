package ui

import tea "github.com/charmbracelet/bubbletea"

type Model interface {
	tea.Model

	SetWidth(width int)
	SetHeight(height int)
}
