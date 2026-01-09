package tui

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/components/help"
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/types"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.showHelp {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, help.View())
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

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.displayTitle(),
		top,
		middle,
	)

	contentHeight := lipgloss.Height(mainContent)
	spacerHeight := m.height - contentHeight - 1
	spacer := ""
	if spacerHeight > 0 {
		spacer = strings.Repeat("\n", spacerHeight-1)
	}

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, spacer, m.statusBar())
}

func (m Model) displayTitle() string {
	title := style.Title.Render("GOSTMAN")
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title)
}

func (m Model) statusBar() string {
	var statusLeft string
	if m.loading {
		statusLeft = "Loading..."
	}

	keyStyle := lipgloss.NewStyle().Foreground(style.ColorGreen)
	sepStyle := style.Unselected

	keybinds := keyStyle.Render("[u]") + sepStyle.Render("rl ") +
		keyStyle.Render("[m]") + sepStyle.Render("ethod ") +
		keyStyle.Render("[h]") + sepStyle.Render("eaders ") +
		keyStyle.Render("[b]") + sepStyle.Render("ody ") +
		keyStyle.Render("[r]") + sepStyle.Render("esponse ") +
		keyStyle.Render("[esc]") + sepStyle.Render("exit mode ") +
		keyStyle.Render("[alt-enter]") + sepStyle.Render("send request ") +
		keyStyle.Render("[q]") + sepStyle.Render("uit")

	helpHint := style.Unselected.Render("? help")

	leftWidth := lipgloss.Width(statusLeft)
	centerWidth := lipgloss.Width(keybinds)
	rightWidth := lipgloss.Width(helpHint)
	totalWidth := m.width - 2

	leftPad := max((totalWidth-centerWidth)/2-leftWidth, 1)
	rightPad := max(totalWidth-leftWidth-leftPad-centerWidth-rightWidth, 1)

	return statusLeft +
		strings.Repeat(" ", leftPad) +
		keybinds +
		strings.Repeat(" ", rightPad) +
		helpHint
}
