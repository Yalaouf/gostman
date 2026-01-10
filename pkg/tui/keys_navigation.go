package tui

import "github.com/Yalaouf/gostman/pkg/tui/types"

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
