package tui

import tea "github.com/charmbracelet/bubbletea"

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
