package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func New() Model {
	t := textinput.New()
	t.Placeholder = "http://localhost:3000"
	t.Cursor.Blink = true
	t.Width = 300
	t.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))

	return Model{
		urlInput:     t,
		focusSection: URL,
	}
}
