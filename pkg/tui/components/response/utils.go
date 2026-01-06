package response

import (
	"fmt"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/lipgloss"
)

func colorStatusCode(statusCode int) string {
	statusStyle := lipgloss.NewStyle()

	switch {
	case statusCode >= 200 && statusCode < 300:
		statusStyle = statusStyle.Foreground(style.ColorGreen)
	case statusCode >= 300 && statusCode < 400:
		statusStyle = statusStyle.Foreground(style.ColorBlue)
	case statusCode >= 400 && statusCode < 500:
		statusStyle = statusStyle.Foreground(style.ColorOrange)
	default:
		statusStyle = statusStyle.Foreground(style.ColorRed)
	}

	return statusStyle.Render(fmt.Sprintf("Status %d", statusCode))
}

func colorTimeTaken(timeTaken int64) string {
	timeTakenStyle := lipgloss.NewStyle()

	switch {
	case timeTaken <= 250:
		timeTakenStyle = timeTakenStyle.Foreground(style.ColorGreen)
	case timeTaken <= 500:
		timeTakenStyle = timeTakenStyle.Foreground(style.ColorYellow)
	default:
		timeTakenStyle = timeTakenStyle.Foreground(style.ColorRed)
	}

	return timeTakenStyle.Render(fmt.Sprintf("%dms", timeTaken))
}
