package help

import (
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/lipgloss"
)

type Section struct {
	Title string
	Keys  []KeyBinding
}

type KeyBinding struct {
	Key  string
	Desc string
}

var Sections = []Section{
	{
		Title: "Navigation",
		Keys: []KeyBinding{
			{Key: "u", Desc: "Focus URL input"},
			{Key: "m", Desc: "Focus method selector"},
			{Key: "h", Desc: "Focus headers"},
			{Key: "b", Desc: "Focus body"},
			{Key: "r", Desc: "Focus response"},
		},
	},
	{
		Title: "Actions",
		Keys: []KeyBinding{
			{Key: "Alt+Enter", Desc: "Send request"},
			{Key: "Enter", Desc: "Enter edit mode"},
			{Key: "Esc", Desc: "Exit edit mode"},
			{Key: "q", Desc: "Quit"},
		},
	},
	{
		Title: "Headers",
		Keys: []KeyBinding{
			{Key: "a", Desc: "Add new header"},
			{Key: "d", Desc: "Delete header"},
			{Key: "p", Desc: "Open presets"},
			{Key: "Space", Desc: "Toggle header"},
			{Key: "j/k", Desc: "Navigate up/down"},
		},
	},
	{
		Title: "Body",
		Keys: []KeyBinding{
			{Key: "Tab", Desc: "Cycle body type"},
		},
	},
	{
		Title: "Response",
		Keys: []KeyBinding{
			{Key: "j/k", Desc: "Scroll up/down"},
			{Key: "g/G", Desc: "Go to top/bottom"},
		},
	},
}

func View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	keyStyle := lipgloss.NewStyle().Foreground(style.ColorGreen).Width(12)
	descStyle := lipgloss.NewStyle().Foreground(style.ColorText)

	title := titleStyle.Render("Keyboard Shortcuts")

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

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Render(title + "\n\n" + content + "\n\n" + style.Unselected.Render("Press ? or Esc to close"))

	return box
}
