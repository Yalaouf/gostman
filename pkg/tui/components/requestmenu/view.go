package requestmenu

import (
	"fmt"
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.inputMode {
		return m.viewInput()
	}

	switch m.viewMode {
	case ViewCollections:
		return m.viewCollections()
	case ViewRequests:
		return m.viewRequests()
	case ViewMoveTarget:
		return m.viewMoveTarget()
	}

	return ""
}

func (m Model) viewCollections() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	title := titleStyle.Render("Saved Requests")

	var b strings.Builder
	b.WriteString("Collections:\n\n")

	for i, coll := range m.collections {
		count := len(m.storage.ListRequestsByCollection(coll.ID))
		line := fmt.Sprintf("%s (%d)", coll.Name, count)
		if i == m.index {
			b.WriteString(style.Selected.Render("▸ " + line))
		} else {
			b.WriteString(style.Unselected.Render("  " + line))
		}
		b.WriteString("\n")
	}

	uncatCount := len(m.storage.ListRequestsByCollection(""))
	uncatLine := fmt.Sprintf("Uncategorized (%d)", uncatCount)
	if m.index == len(m.collections) {
		b.WriteString(style.Selected.Render("▸ " + uncatLine))
	} else {
		b.WriteString(style.Unselected.Render("  " + uncatLine))
	}

	var errView string
	if m.err != "" {
		errView = "\n\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("[enter]open [n]ew [r]ename [d]elete [esc]close")

	content := title + "\n\n" + b.String() + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Width(50).
		Render(content)

	return box
}

func (m Model) viewRequests() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected
	methodStyle := lipgloss.NewStyle().Foreground(style.ColorBlue).Width(8)

	title := titleStyle.Render(m.selectedCollName)

	var b strings.Builder

	if len(m.requests) == 0 {
		b.WriteString(style.Unselected.Render("  No requests in this collection"))
	} else {
		for i, req := range m.requests {
			method := methodStyle.Render(req.Method)
			line := fmt.Sprintf("%s %s", method, req.Name)
			if i == m.index {
				b.WriteString(style.Selected.Render("▸ ") + line)
			} else {
				b.WriteString(style.Unselected.Render("  ") + line)
			}
			b.WriteString("\n")
		}
	}

	var errView string
	if m.err != "" {
		errView = "\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("[enter]load [r]ename [d]elete [m]ove [esc]back")

	content := title + "\n\n" + b.String() + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Width(50).
		Render(content)

	return box
}

func (m Model) viewMoveTarget() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	title := titleStyle.Render("Move to Collection")

	var b strings.Builder

	for i, coll := range m.collections {
		if i == m.index {
			b.WriteString(style.Selected.Render("▸ " + coll.Name))
		} else {
			b.WriteString(style.Unselected.Render("  " + coll.Name))
		}
		b.WriteString("\n")
	}

	uncatLine := "Uncategorized"
	if m.index == len(m.collections) {
		b.WriteString(style.Selected.Render("▸ " + uncatLine))
	} else {
		b.WriteString(style.Unselected.Render("  " + uncatLine))
	}

	var errView string
	if m.err != "" {
		errView = "\n\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("[enter]confirm [esc]cancel")

	content := title + "\n\n" + b.String() + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Width(50).
		Render(content)

	return box
}

func (m Model) viewInput() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	var title string
	switch m.inputAction {
	case InputCreateCollection:
		title = titleStyle.Render("New Collection")
	case InputRenameCollection:
		title = titleStyle.Render("Rename Collection")
	case InputRenameRequest:
		title = titleStyle.Render("Rename Request")
	}

	inputView := m.input.View()

	var errView string
	if m.err != "" {
		errView = "\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("Enter to confirm, Esc to cancel")

	content := title + "\n\n" + inputView + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Render(content)

	return box
}
