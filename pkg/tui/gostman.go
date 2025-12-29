package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Gostman() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
