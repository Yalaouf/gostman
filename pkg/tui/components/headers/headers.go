package headers

import (
	"github.com/Yalaouf/gostman/pkg/tui/style"
)

type Model struct {
	Headers map[string]string
	Focused bool
}

func New() Model {
	return Model{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
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
	content := "Content-Type: application/json\n"
	return style.SectionBox("Headers", content, m.Focused, width)
}
