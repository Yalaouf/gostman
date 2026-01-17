package body

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
)

func (m Model) View(width int) string {
	tabs := m.renderTabs()

	var content string
	if m.BodyType == TypeNone {
		content = style.Unselected.Render("No body")
	} else if m.EditMode {
		content = m.Editor.View()
	} else {
		content = m.Viewport.View()
	}

	footer := style.Unselected.Render("[tab]switch type [enter]edit mode [esc]exit edit")

	body := tabs + "\n" + content + "\n" + footer

	return style.SectionBox("Body", body, m.Focused, width, m.height-4)
}

func (m Model) renderTabs() string {
	var tabs []string

	for _, t := range AllTypes {
		label := t.String()
		if t == m.BodyType {
			tabs = append(tabs, style.Selected.Render("["+label+"]"))
		} else {
			tabs = append(tabs, style.Unselected.Render(" "+label+" "))
		}
	}

	return strings.Join(tabs, " ")
}
