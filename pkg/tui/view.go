package tui

import (
	"fmt"
	"strings"

	"github.com/Yalaouf/gostman/pkg/request"
)

func displayTitle() string {
	return "Gostman\n\n"
}

func displayMethodAndUrl(m Model) string {
	s := ""

	s += string(m.req.Method)
	s += m.urlInput.View()
	s += "\n\n"

	return s
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
	s += fmt.Sprintf("Time taken: %dms\n", m.response.TimeTaken)
	s += fmt.Sprintf("Status: %d\n\n", m.response.StatusCode)
	s += fmt.Sprintf("Body: %s\n\n", m.response.Body)

	s += "Result headers\n"
	for key, value := range m.response.Header {
		s += fmt.Sprintf("%s: %s\n", key, value)
	}

	s += "\n\n"

	return s
}

func displayHelper() string {
	s := ""

	s += "i: URL (esc to leave it) | 0: Method | enter: Send | q: Quit"

	return s
}

func displayError(m Model) string {
	return fmt.Sprintf("Error: %s\n", m.errorMsg)
}

func (m Model) View() string {
	b := strings.Builder{}

	b.WriteString(displayTitle())
	b.WriteString(displayMethodAndUrl(m))
	b.WriteString(displayHeaders())

	if m.req.Method != request.GET {
		b.WriteString(displayBody())
	}

	if m.response.StatusCode != 0 {
		b.WriteString(displayResult(m))
	}

	if m.loading {
		b.WriteString("loading.....\n")
	}

	if m.errorMsg != "" {
		b.WriteString(displayError(m))
	}

	b.WriteString(displayHelper())

	return b.String()
}
