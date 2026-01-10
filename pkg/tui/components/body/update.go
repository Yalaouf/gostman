package body

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.EditMode {
		var cmd tea.Cmd
		m.Editor, cmd = m.Editor.Update(msg)
		return cmd
	}
	return nil
}
