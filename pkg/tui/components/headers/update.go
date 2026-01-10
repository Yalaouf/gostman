package headers

import (
	"github.com/Yalaouf/gostman/pkg/tui/types"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.showPresets {
		return m.updatePresets(msg)
	}

	if m.EditMode {
		return m.updateEdit(msg)
	}

	return m.updateNav(msg)
}

func (m *Model) updateNav(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch keyMsg.String() {
	case types.KeyJ, types.KeyDown:
		if m.cursor < len(m.Headers)-1 {
			m.cursor++
			m.ensureCursorVisible()
		}
	case types.KeyK, types.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m.ensureCursorVisible()
		}
	case types.KeyEnter:
		return m.EnterEditMode()
	case types.KeyA:
		m.addHeader("", "")
		m.updateViewportContent()
		m.ensureCursorVisible()
		return m.EnterEditMode()
	case types.KeyD:
		m.deleteHeader()
		m.updateViewportContent()
		m.ensureCursorVisible()
	case types.KeyP:
		m.showPresets = true
		m.presetCursor = 0
	case types.KeySpace:
		m.toggleHeader()
		m.updateViewportContent()
	}

	m.updateViewportContent()
	return nil
}

func (m *Model) updateEdit(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if ok {
		switch keyMsg.String() {
		case types.KeyEscape:
			m.ExitEditMode()
			m.updateViewportContent()
			return nil
		case types.KeyTab:
			m.Headers[m.cursor].Key.Blur()
			m.Headers[m.cursor].Value.Blur()
			m.fieldFocus = (m.fieldFocus + 1) % 2
			if m.fieldFocus == 0 {
				return m.Headers[m.cursor].Key.Focus()
			}
			return m.Headers[m.cursor].Value.Focus()
		case types.KeyEnter:
			m.ExitEditMode()
			m.updateViewportContent()
			return nil
		}
	}

	var cmd tea.Cmd
	if m.fieldFocus == 0 {
		m.Headers[m.cursor].Key, cmd = m.Headers[m.cursor].Key.Update(msg)
	} else {
		m.Headers[m.cursor].Value, cmd = m.Headers[m.cursor].Value.Update(msg)
	}

	m.updateViewportContent()
	return cmd
}

func (m *Model) updatePresets(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch keyMsg.String() {
	case types.KeyJ, types.KeyDown:
		if m.presetCursor < len(CommonPresets)-1 {
			m.presetCursor++
		}
	case types.KeyK, types.KeyUp:
		if m.presetCursor > 0 {
			m.presetCursor--
		}
	case types.KeyEnter:
		p := CommonPresets[m.presetCursor]
		m.addHeader(p.Key, p.Value)
		m.showPresets = false
		m.updateViewportContent()
		m.ensureCursorVisible()
	case types.KeyEscape, types.KeyP:
		m.showPresets = false
	}

	return nil
}
