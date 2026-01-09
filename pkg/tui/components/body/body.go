package body

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/utils"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Editor   textarea.Model
	BodyType Type
	Focused  bool
	EditMode bool
	height   int
}

func New() Model {
	ta := textarea.New()
	ta.Placeholder = `{"key": "value"}`
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetWidth(40)
	ta.SetHeight(4)
	ta.FocusedStyle.CursorLine = style.TextArea
	ta.BlurredStyle.CursorLine = style.TextArea

	return Model{
		Editor:   ta,
		BodyType: TypeNone,
		Focused:  false,
		EditMode: false,
	}
}

func (m Model) Value() string {
	if m.BodyType == TypeNone {
		return ""
	}
	return m.Editor.Value()
}

func (m Model) Type() Type {
	return m.BodyType
}

func (m *Model) SetSize(width, height int) {
	m.height = height
	m.Editor.SetWidth(width - 6)
	m.Editor.SetHeight(height - 6)
}

func (m *Model) Focus() tea.Cmd {
	m.Focused = true
	m.EditMode = false
	return nil
}

func (m *Model) Blur() {
	m.Focused = false
	m.EditMode = false
	m.Editor.Blur()
}

func (m Model) IsFocused() bool {
	return m.EditMode && m.Editor.Focused()
}

func (m *Model) NextType() {
	idx := int(m.BodyType)
	idx = (idx + 1) % len(AllTypes)
	m.BodyType = AllTypes[idx]
}

func (m *Model) PrevType() {
	idx := int(m.BodyType)
	idx = (idx - 1 + len(AllTypes)) % len(AllTypes)
	m.BodyType = AllTypes[idx]
}

func (m *Model) EnterEditMode() tea.Cmd {
	if m.BodyType == TypeNone {
		return nil
	}
	m.EditMode = true
	return m.Editor.Focus()
}

func (m *Model) ExitEditMode() {
	m.EditMode = false
	m.Editor.Blur()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.EditMode {
		var cmd tea.Cmd
		m.Editor, cmd = m.Editor.Update(msg)
		return cmd
	}
	return nil
}

func (m Model) View(width int) string {
	tabs := m.renderTabs()

	var content string
	editorHeight := m.height - 6
	if m.BodyType == TypeNone {
		content = lipgloss.NewStyle().Height(editorHeight).Render(style.Unselected.Render("No body"))
	} else if m.EditMode {
		content = m.Editor.View()
	} else {
		raw := m.Editor.Value()
		if utils.IsJSON(raw) {
			content = utils.HighlightJSON(raw)
		} else {
			content = m.Editor.View()
		}
	}

	topContent := tabs + "\n" + content
	footer := style.Unselected.Render("[tab]switch type [enter]edit mode [esc]exit edit")

	innerHeight := m.height - 4
	body := lipgloss.Place(width-6, innerHeight, lipgloss.Left, lipgloss.Bottom, footer, lipgloss.WithWhitespaceChars(" "), lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}))
	body = topContent + "\n" + body

	return style.SectionBox("Body", body, m.Focused, width, m.height-4)
}

func (m Model) renderTabs() string {
	var tabs []string

	for _, t := range AllTypes {
		label := t.String()
		if t == m.BodyType {
			tabs = append(tabs, style.Selected.Render("["+label+"]"))
		} else {
			tabs = append(tabs, style.Unselected.Render(" "+label+" "))
		}
	}

	return strings.Join(tabs, " ")
}
