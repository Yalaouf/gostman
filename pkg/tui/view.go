package tui

import (
	"fmt"
	"strings"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	focusedBorder   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#cba6f7"))
)

func displayTitle() string {
	return titleStyle.Render("Gostman\n\n")
}

func displayMethod(m Model) string {
	b := strings.Builder{}

	for i, method := range m.methods {
		if i == m.methodIndex {
			b.WriteString(selectedStyle.Render("▸ " + string(method)))
		} else {
			b.WriteString(selectedStyle.Render("  " + string(method)))
		}
		b.WriteString("\n")
	}

	content := b.String()
	if m.focusSection == METHOD {
		return focusedBorder.Render(content) + "\n"
	}

	return content + "\n"
}

func displayMethodAndUrl(m Model) string {
	method := selectedStyle.Render(string(m.methods[m.methodIndex]))
	return method + m.urlInput.View() + "\n\n"
}

func displayHeaders() string {
	return "Headers\n\n"
}

func displayBody() string {
	return "Body\n\n"
}

func displayResult(m Model) string {
	s := ""

	s += "Result\n"
	s += fmt.Sprintf("Time taken: %dms\n", m.res.TimeTaken)
	s += fmt.Sprintf("Status: %d\n\n", m.res.StatusCode)
	s += fmt.Sprintf("Body: %s\n\n", m.res.Body)

	s += "Result headers\n"
	for key, value := range m.res.Header {
		s += fmt.Sprintf("%s: %s\n", key, value)
	}

	s += "\n\n"

	return s
}

func displayError(m Model) string {
	return errorStyle.Render(fmt.Sprintf("Error: %s\n\n", m.errorMsg))
}

func (m Model) View() string {
	b := strings.Builder{}

	b.WriteString(displayTitle())
	b.WriteString(displayMethodAndUrl(m))
	b.WriteString(displayHeaders())

	if m.req.Method != request.GET {
		b.WriteString(displayBody())
	}

	if m.focusSection == METHOD {
		b.WriteString(displayMethod(m))
	}

	if m.res.StatusCode != 0 {
		b.WriteString(displayResult(m))
	}

	if m.loading {
		b.WriteString("loading.....\n")
	}

	if m.errorMsg != "" {
		b.WriteString(displayError(m))
	}

	return b.String()
}
