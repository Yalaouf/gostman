package url

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return cmd
}
