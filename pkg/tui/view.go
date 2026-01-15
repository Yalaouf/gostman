package tui

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/types"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.showHelp {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.help.View())
	}

	if m.savePopup.Visible() {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.savePopup.View(),
		)
	}

	if m.requestMenu.Visible() {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.requestMenu.View(),
		)
	}

	if m.response.IsFullscreen() {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.response.ViewFullscreen(m.width, m.height),
		)
	}

	if m.focusSection == types.FocusMethod {
		picker := m.method.View()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, picker)
	}

	top := m.url.ViewWithMethod(m.method.ViewSelected(), m.width-2)

	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth - 2

	headersView := m.headers.View(leftWidth)
	bodyView := m.body.View(leftWidth)
	leftPanel := lipgloss.JoinVertical(lipgloss.Left, headersView, bodyView)

	responseView := m.response.View(rightWidth)

	middle := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, responseView)

	statusBar := m.statusBar()

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.displayTitle(),
		top,
		middle,
	)

	lines := strings.Split(mainContent, "\n")
	maxLines := m.height - 1
	if maxLines > 0 && len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	mainContent = strings.Join(lines, "\n")

	return mainContent + "\n\n" + statusBar
}

func (m Model) displayTitle() string {
	title := style.Title.Render("GOSTMAN")
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title)
}

func (m Model) statusBar() string {
	keyStyle := lipgloss.NewStyle().Foreground(style.ColorGreen)
	sepStyle := style.Unselected

	keybinds := keyStyle.Render("[u]") + sepStyle.Render("rl ") +
		keyStyle.Render("[m]") + sepStyle.Render("ethod ") +
		keyStyle.Render("[h]") + sepStyle.Render("eaders ") +
		keyStyle.Render("[b]") + sepStyle.Render("ody ") +
		keyStyle.Render("[r]") + sepStyle.Render("esponse ") +
		keyStyle.Render("[s]") + sepStyle.Render("ave ") +
		keyStyle.Render("[l]") + sepStyle.Render("oad ") +
		keyStyle.Render("[alt-enter]") + sepStyle.Render("send ") +
		keyStyle.Render("[q]") + sepStyle.Render("uit")

	helpHint := style.Unselected.Render("? help")

	centerWidth := lipgloss.Width(keybinds)
	rightWidth := lipgloss.Width(helpHint)
	totalWidth := m.width - 2

	leftPad := max((totalWidth-centerWidth)/2, 1)
	rightPad := max(totalWidth-leftPad-centerWidth-rightWidth, 1)

	return strings.Repeat(" ", leftPad) +
		keybinds +
		strings.Repeat(" ", rightPad) +
		helpHint
}
