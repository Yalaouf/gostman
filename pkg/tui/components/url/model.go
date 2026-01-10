package url

import (
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Input   textinput.Model
	Focused bool
}

func New() Model {
	t := textinput.New()
	t.Placeholder = "http://localhost:3000"
	t.Cursor.Blink = true
	t.TextStyle = style.TextInput

	return Model{
		Input:   t,
		Focused: false,
	}
}

func (m Model) Value() string {
	return m.Input.Value()
}

func (m *Model) SetValue(value string) {
	m.Input.SetValue(value)
}

func (m *Model) SetWidth(width int) {
	m.Input.Width = width
}

func (m *Model) Focus() tea.Cmd {
	m.Focused = true
	m.Input.Focus()
	return textinput.Blink
}

func (m *Model) Blur() {
	m.Focused = false
	m.Input.Blur()
}

func (m Model) IsFocused() bool {
	return m.Input.Focused()
}
