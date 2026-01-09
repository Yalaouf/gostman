package savepopup

import (
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	visible bool
	input   textinput.Model
	err     string
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Request name"
	ti.CharLimit = 64
	ti.Width = 30

	return Model{
		input: ti,
	}
}

func (m *Model) Show() tea.Cmd {
	m.visible = true
	m.err = ""
	m.input.SetValue("")
	m.input.Focus()
	return textinput.Blink
}

func (m *Model) Hide() {
	m.visible = false
	m.input.Blur()
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m *Model) SetError(err string) {
	m.err = err
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	title := titleStyle.Render("Save Request")

	inputView := m.input.View()

	var errView string
	if m.err != "" {
		errView = "\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("Enter to save, Esc to cancel")

	content := title + "\n\n" + inputView + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Render(content)

	return box
}
