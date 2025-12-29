package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
			b.WriteString(unselectedStyle.Render("  " + string(method)))
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
	content := method + " " + m.urlInput.View()

	return sectionBox("Request", content, m.focusSection == URL, m.width)
}

func displayHeaders(m Model) string {
	content := "Content-Type: application/json\n"
	return sectionBox("Headers", content, m.focusSection == HEADERS, m.width/2)
}

func displayBody(m Model) string {
	content := "{ }\n"
	return sectionBox("Body", content, m.focusSection == BODY, m.width/2)
}

func displayResult(m Model) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Status %d  •  Time: %dms\n\n", m.res.StatusCode, m.res.TimeTaken))
	b.WriteString(m.res.Body)
	return sectionBox("Response", b.String(), false, m.width)
}

func sectionBox(title, content string, focused bool, width int) string {
	style := sectionStyle
	if focused {
		style = focusedSectionStyle
	}

	header := sectionTitleStyle.Render(title)
	body := style.Width(width - 4).Render(content)

	return header + "\n" + body
}

func (m Model) View() string {
	top := displayMethodAndUrl(m)

	headers := displayHeaders(m)
	body := displayBody(m)
	middle := lipgloss.JoinHorizontal(lipgloss.Top, headers, body)

	var bottom string
	if m.res.StatusCode != 0 {
		bottom = displayResult(m)
	}

	var status string
	if m.loading {
		status = "Loading..."
	}

	if m.errorMsg != "" {
		status = errorStyle.Render("Error: ", m.errorMsg)
	}

	var methodPicker string
	if m.focusSection == METHOD {
		methodPicker = displayMethod(m)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		displayTitle(),
		methodPicker,
		top,
		middle,
		bottom,
		status,
	)
}
