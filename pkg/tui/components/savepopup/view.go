package savepopup

import (
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	title := titleStyle.Render("Save Request")

	inputView := m.input.View()

	var errView string
	if m.err != "" {
		errView = "\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("Enter to save, Esc to cancel")

	content := title + "\n\n" + inputView + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Render(content)

	return box
}
