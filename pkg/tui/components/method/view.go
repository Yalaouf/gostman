package method

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
)

func (m Model) View() string {
	var b strings.Builder

	for i, method := range m.Methods {
		if i == m.Index {
			b.WriteString(style.Selected.Render("▸ " + string(method)))
		} else {
			b.WriteString(style.Unselected.Render("  " + string(method)))
		}
		b.WriteString("\n")
	}

	content := b.String()
	if m.Focused {
		return style.FocusedBorder.Render(content) + "\n"
	}

	return content + "\n"
}

func (m Model) ViewSelected() string {
	return style.Selected.Render(string(m.Selected()))
}
