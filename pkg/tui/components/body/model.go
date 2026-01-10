package body

import (
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
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

func (m *Model) SetValue(value string) {
	m.Editor.SetValue(value)
}

func (m *Model) SetType(t Type) {
	m.BodyType = t
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
