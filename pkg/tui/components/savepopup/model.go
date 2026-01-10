package savepopup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	visible bool
	input   textinput.Model
	err     string
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Request name"
	ti.CharLimit = 64
	ti.Width = 30

	return Model{
		input: ti,
	}
}

func (m *Model) Show() tea.Cmd {
	m.visible = true
	m.err = ""
	m.input.SetValue("")
	m.input.Focus()
	return textinput.Blink
}

func (m *Model) Hide() {
	m.visible = false
	m.input.Blur()
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m *Model) SetError(err string) {
	m.err = err
}
