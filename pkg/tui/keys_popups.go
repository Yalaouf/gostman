package tui

import (
	"github.com/Yalaouf/gostman/pkg/storage"
	"github.com/Yalaouf/gostman/pkg/tui/components/body"
	"github.com/Yalaouf/gostman/pkg/tui/types"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

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

func (m Model) handleResponseFullscreen(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case types.KeyEscape, types.KeyF:
		m.response.ExitFullscreen()
		return m, nil
	case types.KeyTab:
		m.response.NextTab()
		return m, nil
	case types.KeyY:
		content := m.response.GetContent()
		clipboard.WriteAll(content)
		return m, nil
	}

	if m.response.IsTreeTab() && m.response.HasTree() {
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
		}
	} else {
		switch key {
		case types.KeyJ, types.KeyDown:
			m.response.ScrollDown(1)
		case types.KeyK, types.KeyUp:
			m.response.ScrollUp(1)
		case types.KeyG:
			m.response.GotoTop()
		case types.KeyShiftG:
			m.response.GotoBottom()
		}
	}

	return m, nil
}
