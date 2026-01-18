package response

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/utils"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) renderTabs() string {
	var tabs []string
	for _, t := range AllTabs {
		label := t.String()
		if t == m.currentTab {
			tabs = append(tabs, style.Selected.Render("["+label+"]"))
		} else {
			tabs = append(tabs, style.Unselected.Render(" "+label+" "))
		}
	}

	return strings.Join(tabs, "")
}

func (m Model) View(width int) string {
	borderColor := style.ColorGray
	if m.Focused {
		borderColor = style.ColorPurple
	}

	tabs := m.renderTabs()

	var content string
	if m.Loading {
		content = style.Unselected.Render("Loading...")
	} else if m.Error != "" {
		content = style.Error.Render("Error: " + m.Error)
	} else if m.HasResponse() {
		scrollbarStyle := lipgloss.NewStyle().Foreground(borderColor).MarginLeft(1)
		scrollbar := scrollbarStyle.Render(RenderScrollbar(m.Viewport))
		content = lipgloss.JoinHorizontal(lipgloss.Top, m.Viewport.View(), scrollbar)
	} else {
		content = style.Unselected.Render("No response yet. Press " + utils.SendRequestShortcut() + " to send a request.")
	}

	fullContent := tabs + "\n\n" + content

	return style.SectionBox("Response", fullContent, m.Focused, width, m.height-4)
}

func (m *Model) ViewFullscreen(width, height int) string {
	fsWidth := width - 10
	fsHeight := height - 6

	m.Viewport.Width = fsWidth - 8
	m.Viewport.Height = fsHeight - 8
	m.updateViewportContent()

	tabs := m.renderTabs()

	var content string
	if m.Loading {
		content = style.Unselected.Render("Loading...")
	} else if m.Error != "" {
		content = style.Error.Render("Error: " + m.Error)
	} else if m.HasResponse() {
		scrollbarStyle := lipgloss.NewStyle().Foreground(style.ColorPurple).MarginLeft(1)
		scrollbar := scrollbarStyle.Render(RenderScrollbar(m.Viewport))
		content = lipgloss.JoinHorizontal(lipgloss.Top, m.Viewport.View(), scrollbar)
	} else {
		content = style.Unselected.Render("No response yet.")
	}

	var hintText string
	if m.currentTab == TabTree {
		hintText = "[f/esc] close  [j/k] navigate  [h/l] collapse/expand  [y] copy  [tab] switch"
	} else {
		hintText = "[f/esc] close  [j/k] scroll  [tab] switch tab  [g/G] top/bottom"
	}
	hint := style.Unselected.Render(hintText)
	fullContent := tabs + "\n\n" + content + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 2).
		Width(fsWidth).
		Height(fsHeight).
		Render(fullContent)

	return box
}
