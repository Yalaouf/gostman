package tui

import (
	"github.com/Yalaouf/gostman/pkg/tui/types"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == types.KeyQuestion {
		m.showHelp = !m.showHelp
		return m, nil
	}

	if m.showHelp {
		switch key {
		case types.KeyEscape:
			m.showHelp = false
		case types.KeyJ, types.KeyDown:
			m.help.ScrollDown()
		case types.KeyK, types.KeyUp:
			m.help.ScrollUp()
		}
		return m, nil
	}

	if m.savePopup.Visible() {
		return m.handleSavePopup(msg)
	}

	if m.requestMenu.Visible() {
		return m.handleRequestMenu(msg)
	}

	if m.response.IsFullscreen() {
		return m.handleResponseFullscreen(msg)
	}

	if key == types.KeyAltEnter || key == types.KeyCtrlG {
		m.response.SetLoading(true)
		m.response.Error = ""
		return m, m.sendRequest()
	}

	if key == types.KeyEscape {
		return m.handleEscape()
	}

	if key == types.KeyCtrlC {
		return m, tea.Quit
	}

	if m.url.IsFocused() {
		return m.handleURLInput(msg)
	}

	if m.body.IsFocused() {
		return m.handleBodyInput(msg)
	}

	if m.headers.EditMode {
		return m.handleHeadersInput(msg)
	}

	if key == types.KeyQ {
		return m, tea.Quit
	}

	if m.focusSection == types.FocusBody {
		switch key {
		case types.KeyTab:
			return m.handleNavigation(key), nil
		}
	}

	if m.focusSection == types.FocusHeaders {
		switch key {
		case types.KeyJ, types.KeyK, types.KeyUp, types.KeyDown, types.KeyEnter, types.KeyTab:
			return m.handleHeadersInput(msg)
		}

		if key == types.KeyA || key == types.KeyD || key == types.KeyP || key == types.KeyEnter ||
			key == types.KeySpace {
			return m.handleHeadersInput(msg)
		}
	}

	if m.focusSection == types.FocusResult && m.response.IsTreeTab() && m.response.HasTree() {
		switch key {
		case types.KeyJ, types.KeyDown:
			m.response.TreeDown()
			return m, nil
		case types.KeyK, types.KeyUp:
			m.response.TreeUp()
			return m, nil
		case types.KeyH, types.KeyLeft:
			m.response.TreeCollapse()
			return m, nil
		case types.KeyL, types.KeyRight:
			m.response.TreeExpand()
			return m, nil
		case types.KeyEnter, types.KeySpace:
			m.response.TreeToggle()
			return m, nil
		case types.KeyTab:
			m.response.NextTab()
			return m, nil
		case types.KeyF:
			m.response.ToggleFullscreen()
			return m, nil
		case types.KeyY:
			content := m.response.GetContent()
			clipboard.WriteAll(content)
			return m, nil
		}
	}

	switch key {
	case types.KeyM:
		return m.handleFocusChange(types.FocusMethod)
	case types.KeyU:
		return m.handleFocusChange(types.FocusURL)
	case types.KeyH:
		return m.handleFocusChange(types.FocusHeaders)
	case types.KeyB:
		return m.handleFocusChange(types.FocusBody)
	case types.KeyR:
		if m.response.HasResponse() {
			return m.handleFocusChange(types.FocusResult)
		}
		return m, nil
	case types.KeyS:
		return m, m.savePopup.Show()
	case types.KeyL:
		return m, m.requestMenu.Show()
	}

	switch key {
	case types.KeyJ, types.KeyDown, types.KeyK, types.KeyUp, types.KeyTab:
		return m.handleNavigation(key), nil
	}

	switch key {
	case types.KeyG, types.KeyShiftG:
		return m.handleScroll(key), nil
	}

	if key == types.KeyF && m.focusSection == types.FocusResult && m.response.HasResponse() {
		m.response.ToggleFullscreen()
		return m, nil
	}

	if key == types.KeyY && m.focusSection == types.FocusResult && m.response.HasResponse() {
		content := m.response.GetContent()
		clipboard.WriteAll(content)
		return m, nil
	}

	if key == types.KeyEnter {
		return m.handleEnter()
	}

	return m, nil
}
