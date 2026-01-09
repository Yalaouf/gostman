package tui

import (
	"github.com/Yalaouf/gostman/pkg/storage"
	"github.com/Yalaouf/gostman/pkg/tui/components/body"
	"github.com/Yalaouf/gostman/pkg/tui/types"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == types.KeyQuestion {
		m.showHelp = !m.showHelp
		return m, nil
	}

	if m.showHelp {
		if key == types.KeyEscape {
			m.showHelp = false
			return m, nil
		}
		return m, nil
	}

	if m.savePopup.Visible() {
		return m.handleSavePopup(msg)
	}

	if m.requestMenu.Visible() {
		return m.handleRequestMenu(msg)
	}

	if key == types.KeyAltEnter {
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

		if key == types.KeyA || key == types.KeyD || key == types.KeyP || key == types.KeyEnter || key == types.KeySpace {
			return m.handleHeadersInput(msg)
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

	if key == types.KeyEnter {
		return m.handleEnter()
	}

	return m, nil
}

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

func (m Model) handleNavigation(key string) Model {
	switch m.focusSection {
	case types.FocusMethod:
		if key == types.KeyJ || key == types.KeyDown {
			m.method.Next()
		} else {
			m.method.Previous()
		}
	case types.FocusBody:
		if key == types.KeyTab {
			m.body.NextType()
			m.syncContentType()
		}
	case types.FocusResult:
		switch key {
		case types.KeyJ, types.KeyDown:
			m.response.ScrollDown(1)
		case types.KeyK, types.KeyUp:
			m.response.ScrollUp(1)
		case types.KeyTab:
			m.response.NextTab()
		}
	}
	return m
}

func (m Model) handleScroll(key string) Model {
	if m.focusSection != types.FocusResult {
		return m
	}

	if key == types.KeyG {
		m.response.GotoTop()
	} else {
		m.response.GotoBottom()
	}
	return m
}

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

func (m Model) handleURLInput(msg tea.Msg) (Model, tea.Cmd) {
	cmd := m.url.Update(msg)
	return m, cmd
}

func (m Model) handleBodyInput(msg tea.Msg) (Model, tea.Cmd) {
	cmd := m.body.Update(msg)
	return m, cmd
}

func (m Model) handleHeadersInput(msg tea.Msg) (Model, tea.Cmd) {
	cmd := m.headers.Update(msg)
	return m, cmd
}

func (m Model) handleSavePopup(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case types.KeyEscape:
		m.savePopup.Hide()
		return m, nil

	case types.KeyEnter:
		name := m.savePopup.Value()
		if name == "" {
			m.savePopup.SetError("Name is required")
			return m, nil
		}

		req := &storage.Request{
			Name:    name,
			Method:  string(m.method.Selected()),
			URL:     m.url.Value(),
			Headers: m.headers.EnabledHeaders(),
			Body:    m.body.Value(),
		}

		switch m.body.BodyType {
		case body.TypeJSON:
			req.BodyType = "json"
		case body.TypeFormData:
			req.BodyType = "form-data"
		case body.TypeURLEncoded:
			req.BodyType = "urlencoded"
		default:
			req.BodyType = "none"
		}

		if err := m.storage.SaveRequest(req); err != nil {
			m.savePopup.SetError(err.Error())
			return m, nil
		}

		m.savePopup.Hide()
		return m, nil
	}

	cmd := m.savePopup.Update(msg)
	return m, cmd
}

func (m Model) handleRequestMenu(msg tea.KeyMsg) (Model, tea.Cmd) {
	cmd := m.requestMenu.Update(msg)
	return m, cmd
}
