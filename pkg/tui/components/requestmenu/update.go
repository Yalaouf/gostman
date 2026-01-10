package requestmenu

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.inputMode {
		return m.handleInputMode(msg)
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	key := keyMsg.String()

	switch key {
	case "esc":
		return m.handleEscape()
	case "j", "down":
		m.moveDown()
	case "k", "up":
		m.moveUp()
	case "enter":
		return m.handleEnter()
	case "n":
		if m.viewMode == ViewCollections {
			m.startCreateCollection()
			return textinput.Blink
		}
	case "r":
		m.startRename()
		return textinput.Blink
	case "d":
		m.deleteSelected()
	case "m":
		if m.viewMode == ViewRequests && len(m.requests) > 0 {
			m.startMove()
		}
	}

	return nil
}

func (m *Model) handleInputMode(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return cmd
	}

	key := keyMsg.String()

	switch key {
	case "esc":
		m.inputMode = false
		m.inputAction = InputNone
		m.input.Blur()
		return nil
	case "enter":
		return m.confirmInput()
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m *Model) handleEscape() tea.Cmd {
	switch m.viewMode {
	case ViewCollections:
		m.Hide()
	case ViewRequests:
		m.viewMode = ViewCollections
		m.index = 0
		m.refresh()
	case ViewMoveTarget:
		m.viewMode = ViewRequests
		m.moveRequestID = ""
		m.index = 0
	}
	return nil
}

func (m *Model) handleEnter() tea.Cmd {
	switch m.viewMode {
	case ViewCollections:
		m.enterCollection()
	case ViewRequests:
		if len(m.requests) > 0 && m.index < len(m.requests) {
			req := m.requests[m.index]
			m.Hide()
			return func() tea.Msg {
				return LoadRequestMsg{Request: req}
			}
		}
	case ViewMoveTarget:
		return m.confirmMove()
	}
	return nil
}
