package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case requestMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.response = msg.response
		return m, nil
	case tea.KeyMsg:
		if !m.urlInput.Focused() {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "0":
				m.focusSection = METHOD
				m.urlInput.Blur()
				return m, nil
			case "i":
				m.focusSection = URL
				m.urlInput.Focus()
				return m, textinput.Blink
			case "2":
				m.focusSection = HEADERS
				m.urlInput.Blur()
				return m, nil
			case "3":
				m.focusSection = BODY
				m.urlInput.Blur()
				return m, nil
			case "enter":
				m.loading = true
				m.errorMsg = ""
				return m, m.sendRequest()
			}
		}
		if msg.String() == "esc" {
			m.urlInput.Blur()
			return m, nil
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if m.focusSection == URL {
		m.urlInput, cmd = m.urlInput.Update(msg)
		m.req.SetURL(m.urlInput.Value())
	}

	return m, cmd
}
