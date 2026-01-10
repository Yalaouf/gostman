package savepopup

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}
