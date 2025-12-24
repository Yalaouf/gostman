package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Gostman() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
