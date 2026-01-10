package tui

import (
	"github.com/Yalaouf/gostman/pkg/tui/types"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focusSection {
	case types.FocusMethod:
		return m.handleFocusChange(types.FocusURL)
	case types.FocusBody:
		return m, m.body.EnterEditMode()
	case types.FocusHeaders:
		return m, m.headers.EnterEditMode()
	}

	return m, nil
}

func (m Model) handleEscape() (Model, tea.Cmd) {
	switch m.focusSection {
	case types.FocusMethod:
		m.focusSection = types.FocusURL
	case types.FocusBody:
		if m.body.IsFocused() {
			m.body.ExitEditMode()
			return m, nil
		}
	case types.FocusHeaders:
		if m.headers.IsFocused() {
			m.headers.ExitEditMode()
			return m, nil
		}
	}
	m.url.Blur()
	return m, nil
}
