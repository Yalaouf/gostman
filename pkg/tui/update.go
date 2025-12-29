package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.urlInput.Width = msg.Width - 10
		return m, nil
	case requestMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			return m, nil
		}
		m.res = msg.response
		return m, nil
	case tea.KeyMsg:
		if !m.urlInput.Focused() {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "m":
				m.focusSection = METHOD
				m.urlInput.Blur()
				return m, nil
			case "i":
				m.focusSection = URL
				m.urlInput.Focus()
				return m, textinput.Blink
			case "h":
				m.focusSection = HEADERS
				m.urlInput.Blur()
				return m, nil
			case "b":
				m.focusSection = BODY
				m.urlInput.Blur()
				return m, nil
			case "enter":
				switch m.focusSection {
				case METHOD:
					m.req.SetMethod(m.methods[m.methodIndex])
					m.focusSection = URL
					m.urlInput.Focus()
					return m, textinput.Blink
				default:
					m.loading = true
					m.errorMsg = ""
					return m, m.sendRequest()
				}
			case "j", "down":
				if m.focusSection == METHOD {
					m.methodIndex = (m.methodIndex + 1) % len(m.methods)
					m.req.SetMethod(m.methods[m.methodIndex])
				}
				return m, nil
			case "k", "up":
				if m.focusSection == METHOD {
					m.methodIndex = (m.methodIndex - 1) % len(m.methods)
					m.req.SetMethod(m.methods[m.methodIndex])
				}
				return m, nil
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
