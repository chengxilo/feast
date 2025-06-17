package main

import (
	"feast/ui/root"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	p := tea.NewProgram(root.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
