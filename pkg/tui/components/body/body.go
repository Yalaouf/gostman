package body

import (
	"github.com/Yalaouf/gostman/pkg/tui/style"
)

type Model struct {
	Content string
	Focused bool
}

func New() Model {
	return Model{
		Content: "{ }",
		Focused: false,
	}
}

func (m *Model) Focus() {
	m.Focused = true
}

func (m *Model) Blur() {
	m.Focused = false
}

func (m Model) View(width int) string {
	content := m.Content + "\n"
	return style.SectionBox("Body", content, m.Focused, width)
}
