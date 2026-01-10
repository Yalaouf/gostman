package url

import "github.com/Yalaouf/gostman/pkg/tui/style"

func (m Model) View() string {
	return m.Input.View()
}

func (m Model) ViewWithMethod(methodView string, width int) string {
	content := methodView + " " + m.Input.View()
	return style.SectionBox("Request", content, m.Focused, width)
}
