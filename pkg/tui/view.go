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

	headersView := m.headers.View(m.width / 2)
	bodyView := m.body.View(m.width / 2)
	middle := lipgloss.JoinHorizontal(lipgloss.Top, headersView, bodyView)

	var bottom string
	if m.response.HasResponse() {
		bottom = m.response.View()
	}

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.displayTitle(),
		top,
		middle,
		bottom,
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

	if m.errorMsg != "" {
		statusLeft = style.Error.Render("Error: ", m.errorMsg)
	}

	keyStyle := lipgloss.NewStyle().Foreground(style.ColorGreen)
	sepStyle := style.Unselected

	keybinds := keyStyle.Render("[u]") + sepStyle.Render("rl ") +
		keyStyle.Render("[m]") + sepStyle.Render("ethod ") +
		keyStyle.Render("[h]") + sepStyle.Render("eaders ") +
		keyStyle.Render("[b]") + sepStyle.Render("ody ") +
		keyStyle.Render("[r]") + sepStyle.Render("esponse ") +
		keyStyle.Render("[esc]") + sepStyle.Render(" exit mode")

	helpHint := style.Unselected.Render("? help")

	leftWidth := lipgloss.Width(statusLeft)
	centerWidth := lipgloss.Width(keybinds)
	rightWidth := lipgloss.Width(helpHint)
	totalWidth := m.width - 2

	// Center the keybinds
	leftPad := (totalWidth-centerWidth)/2 - leftWidth
	if leftPad < 1 {
		leftPad = 1
	}
	rightPad := totalWidth - leftWidth - leftPad - centerWidth - rightWidth
	if rightPad < 1 {
		rightPad = 1
	}

	return statusLeft +
		strings.Repeat(" ", leftPad) +
		keybinds +
		strings.Repeat(" ", rightPad) +
		helpHint
}
