package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
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

	return sectionBox("Request", content, m.focusSection == URL, m.width-2)
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
	header := sectionTitleStyle.Render("Response")

	borderColor := lipgloss.Color("#6c7086")
	if m.focusSection == RESULT {
		borderColor = lipgloss.Color("#cba6f7")
	}

	scrollbarStyle := lipgloss.NewStyle().Foreground(borderColor).MarginLeft(1)

	scrollbar := scrollbarStyle.Render(renderScrollbar(m.responseView))

	body := lipgloss.JoinHorizontal(lipgloss.Top, m.responseView.View(), scrollbar)

	return header + "\n" + body
}

func renderScrollbar(viewport viewport.Model) string {
	height := viewport.Height
	if height <= 0 {
		return ""
	}

	totalLines := viewport.TotalLineCount()

	if totalLines <= height {
		return strings.Repeat(trackChar+"\n", height)
	}

	thumbSize := max(1, height*height/totalLines)
	thumbPos := int(viewport.ScrollPercent() * float64(height-thumbSize))

	var b strings.Builder
	for i := range height {
		if i >= thumbPos && i < thumbPos+thumbSize {
			b.WriteString(thumbChar)
		} else {
			b.WriteString(trackChar)
		}

		if i < height-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
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
