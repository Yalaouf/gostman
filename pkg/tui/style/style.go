package style

import "github.com/charmbracelet/lipgloss"

var (
	ColorBlue   = lipgloss.Color("#89b4fa")
	ColorGray   = lipgloss.Color("#6c7086")
	ColorGreen  = lipgloss.Color("#a6e3a1")
	ColorOrange = lipgloss.Color("#fab387")
	ColorPurple = lipgloss.Color("#cba6f7")
	ColorRed    = lipgloss.Color("#f38ba8")
	ColorText   = lipgloss.Color("#cdd6f4")
	ColorYellow = lipgloss.Color("#f9e2af")
)

var (
	Error        = lipgloss.NewStyle().Foreground(ColorRed)
	Title        = lipgloss.NewStyle().Bold(true).Foreground(ColorOrange).Padding(1, 0)
	Selected     = lipgloss.NewStyle().Foreground(ColorGreen)
	Unselected   = lipgloss.NewStyle().Foreground(ColorGray)
	SectionTitle = lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	TextInput    = lipgloss.NewStyle().Foreground(ColorText)
	TextArea     = lipgloss.NewStyle()
)

var (
	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple)

	Section = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorGray).
		Padding(0, 1)

	FocusedSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple).
			Padding(0, 1)

	Viewport = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGray).
			Padding(0, 1)
)

const (
	TrackChar = "▒"
	ThumbChar = "▓"
)

const ChromaStyle = "catppuccin-mocha"

func SectionBox(title, content string, focused bool, width int) string {
	style := Section
	if focused {
		style = FocusedSection
	}

	header := SectionTitle.Render(title)
	body := style.Width(width - 4).Render(content)

	return header + "\n" + body
}
