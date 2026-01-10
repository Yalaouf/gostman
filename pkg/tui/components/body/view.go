package body

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/utils"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View(width int) string {
	tabs := m.renderTabs()

	var content string
	editorHeight := m.height - 6
	if m.BodyType == TypeNone {
		content = lipgloss.NewStyle().
			Height(editorHeight).
			Render(style.Unselected.Render("No body"))
	} else if m.EditMode {
		content = m.Editor.View()
	} else {
		raw := m.Editor.Value()
		if utils.IsJSON(raw) {
			content = utils.HighlightJSON(raw)
		} else {
			content = m.Editor.View()
		}
	}

	topContent := tabs + "\n" + content
	footer := style.Unselected.Render("[tab]switch type [enter]edit mode [esc]exit edit")

	innerHeight := m.height - 4
	body := lipgloss.Place(
		width-6,
		innerHeight,
		lipgloss.Left,
		lipgloss.Bottom,
		footer,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
	)
	body = topContent + "\n" + body

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
