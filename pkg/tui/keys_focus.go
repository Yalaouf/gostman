package tui

import (
	"github.com/Yalaouf/gostman/pkg/tui/types"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleFocusChange(section types.FocusSection) (Model, tea.Cmd) {
	m.focusSection = section

	m.method.Blur()
	m.url.Blur()
	m.headers.Blur()
	m.body.Blur()
	m.response.Blur()

	switch section {
	case types.FocusMethod:
		m.method.Focus()
		return m, nil
	case types.FocusURL:
		m.url.Focused = true
		return m, m.url.Focus()
	case types.FocusHeaders:
		m.headers.Focus()
		return m, nil
	case types.FocusBody:
		return m, m.body.Focus()
	case types.FocusResult:
		m.response.Focus()
		return m, nil
	}

	return m, nil
}
