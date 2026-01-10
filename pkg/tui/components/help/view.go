package help

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func renderContent() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	keyStyle := lipgloss.NewStyle().Foreground(style.ColorGreen).Width(12)
	descStyle := lipgloss.NewStyle().Foreground(style.ColorText)

	var lines []string
	for i, section := range Sections {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, titleStyle.Render(section.Title))
		for _, kb := range section.Keys {
			line := "  " + keyStyle.Render(kb.Key) + descStyle.Render(kb.Desc)
			lines = append(lines, line)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func renderScrollbar(vp viewport.Model) string {
	height := vp.Height
	if height <= 0 {
		return ""
	}

	totalLines := vp.TotalLineCount()

	if totalLines <= height {
		return strings.Repeat(style.TrackChar+"\n", height)
	}

	thumbSize := max(1, height*height/totalLines)
	thumbPos := int(vp.ScrollPercent() * float64(height-thumbSize))

	var b strings.Builder
	for i := range height {
		if i >= thumbPos && i < thumbPos+thumbSize {
			b.WriteString(style.ThumbChar)
		} else {
			b.WriteString(style.TrackChar)
		}

		if i < height-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	title := titleStyle.Render("Keyboard Shortcuts")

	scrollbar := renderScrollbar(m.viewport)
	contentWithScrollbar := lipgloss.JoinHorizontal(lipgloss.Top, m.viewport.View(), " ", scrollbar)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Render(title + "\n\n" + contentWithScrollbar + "\n\n" + style.Unselected.Render("Press ? or Esc to close"))

	return box
}
