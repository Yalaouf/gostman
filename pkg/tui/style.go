package tui

import "github.com/charmbracelet/lipgloss"

var (
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff000000"))
	titleStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fab387"))
	selectedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))
	unselectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	focusedBorder       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#cba6f7"))
	sectionStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#6c7086")).Padding(0, 1)
	focusedSectionStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#cba6f7")).Padding(0, 1)
	sectionTitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cdd6f4"))
)
