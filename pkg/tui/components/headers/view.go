package headers

import (
	"fmt"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) updateViewportContent() {
	var content string
	if len(m.Headers) == 0 {
		content = style.Unselected.Render("No headers (press 'a' to add)")
	} else {
		for i, h := range m.Headers {
			line := m.renderHeaderLine(i, h)
			content += line + "\n"
		}
	}
	m.viewport.SetContent(content)
}

func (m Model) renderHeaderLine(index int, h Header) string {
	isCursor := index == m.cursor && m.Focused

	check := "[ ]"
	if h.Enabled {
		check = "[X]"
	}

	var key, value string
	if m.EditMode && isCursor && m.fieldFocus == 0 {
		key = h.Key.View()
	} else {
		key = h.Key.Value()
	}

	if m.EditMode && isCursor && m.fieldFocus == 1 {
		value = h.Value.View()
	} else {
		value = h.Value.Value()
	}

	line := fmt.Sprintf("%s %s: %s", check, key, value)

	if h.Auto {
		line += style.Unselected.Render(" (auto)")
	}

	if !h.Enabled {
		line = style.Unselected.Render(line)
	} else if isCursor {
		line = lipgloss.NewStyle().Background(style.ColorSurface).Foreground(style.ColorText).Render(line)
	}

	return line
}

func (m Model) viewPresets(width int) string {
	content := "Select a preset:\n\n"

	for i, p := range CommonPresets {
		line := fmt.Sprintf("%s: %s", p.Key, p.Value)
		if i == m.presetCursor {
			line = lipgloss.NewStyle().
				Background(style.ColorSurface).
				Foreground(style.ColorText).
				Render("> " + line)
		} else {
			line = " " + line
		}

		content += line + "\n"
	}

	content += "\n" + style.Unselected.Render("[enter] select		[esc] cancel")
	return style.SectionBox("Headers - Presets", content, m.Focused, width)
}

func (m Model) View(width int) string {
	if m.showPresets {
		return m.viewPresets(width)
	}

	topContent := m.viewport.View()
	footer := style.Unselected.Render(
		"[a]dd [d]el [p]resets [space]toggle [tab]key<>value [esc/enter]validate",
	)

	// Combine viewport and footer, constrained to available height
	content := topContent + "\n" + footer

	return style.SectionBox("Headers", content, m.Focused, width, m.height-4)
}
