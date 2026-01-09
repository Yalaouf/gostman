package headers

import (
	"fmt"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/types"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	Key     textinput.Model
	Value   textinput.Model
	Enabled bool
	Auto    bool // Managed or not for Content-Type (maybe overkill ?)
}

type Model struct {
	Headers      []Header
	cursor       int // selected row
	fieldFocus   int // 0 = key, 1 = value
	Focused      bool
	EditMode     bool
	showPresets  bool
	presetCursor int
	width        int
	height       int
	viewport     viewport.Model
}

func newTextInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder

	return ti
}

func newHeader(key, value string, auto bool) Header {
	k := newTextInput("Key")
	k.SetValue(key)

	v := newTextInput("Value")
	v.SetValue(value)

	return Header{Key: k, Value: v, Enabled: true, Auto: auto}
}

func New() Model {
	vp := viewport.New(40, 4)

	m := Model{
		Headers:    []Header{},
		cursor:     0,
		fieldFocus: 0,
		Focused:    false,
		EditMode:   false,
		viewport:   vp,
	}

	m.updateViewportContent()
	return m
}

func (m *Model) Focus() tea.Cmd {
	m.Focused = true
	return nil
}

func (m *Model) Blur() {
	m.Focused = false
	m.EditMode = false
	m.showPresets = false
}

func (m *Model) EnterEditMode() tea.Cmd {
	if len(m.Headers) == 0 {
		return nil
	}

	m.EditMode = true
	h := &m.Headers[m.cursor]
	if m.fieldFocus == 0 {
		return h.Key.Focus()
	}

	return h.Value.Focus()
}

func (m *Model) ExitEditMode() {
	m.EditMode = false
	m.showPresets = false
	if len(m.Headers) > 0 {
		m.Headers[m.cursor].Key.Blur()
		m.Headers[m.cursor].Value.Blur()
	}
}

func (m *Model) IsFocused() bool {
	return m.EditMode || m.showPresets
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 4
	m.viewport.Height = height - 6
}

func (m *Model) SetContentType(contentType string) {
	for i := range m.Headers {
		if m.Headers[i].Key.Value() == "Content-Type" && m.Headers[i].Auto {
			if contentType == "" {
				m.Headers = append(m.Headers[:i], m.Headers[i+1:]...)

				if m.cursor >= len(m.Headers) && m.cursor > 0 {
					m.cursor--
				}
			} else {
				m.Headers[i].Value.SetValue(contentType)
			}

			m.updateViewportContent()
			return
		}
	}

	if contentType != "" {
		h := newHeader("Content-Type", contentType, true)
		m.Headers = append([]Header{h}, m.Headers...)
		m.updateViewportContent()
	}
}

func (m Model) EnabledHeaders() map[string]string {
	result := make(map[string]string)

	for _, h := range m.Headers {
		if h.Enabled && h.Key.Value() != "" {
			result[h.Key.Value()] = h.Value.Value()
		}
	}

	return result
}

func (m *Model) addHeader(key, value string) {
	h := newHeader(key, value, false)
	m.Headers = append(m.Headers, h)
	m.cursor = len(m.Headers) - 1
}

func (m *Model) deleteHeader() {
	if len(m.Headers) == 0 {
		return
	}

	if m.Headers[m.cursor].Auto {
		return
	}

	m.Headers = append(m.Headers[:m.cursor], m.Headers[m.cursor+1:]...)

	if m.cursor >= len(m.Headers) && m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) toggleHeader() {
	if len(m.Headers) > 0 {
		m.Headers[m.cursor].Enabled = !m.Headers[m.cursor].Enabled
	}
}

func (m *Model) updateNav(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch keyMsg.String() {
	case types.KeyJ, types.KeyDown:
		if m.cursor < len(m.Headers)-1 {
			m.cursor++
			m.ensureCursorVisible()
		}
	case types.KeyK, types.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m.ensureCursorVisible()
		}
	case types.KeyEnter:
		return m.EnterEditMode()
	case types.KeyA:
		m.addHeader("", "")
		m.updateViewportContent()
		m.ensureCursorVisible()
		return m.EnterEditMode()
	case types.KeyD:
		m.deleteHeader()
		m.updateViewportContent()
		m.ensureCursorVisible()
	case types.KeyP:
		m.showPresets = true
		m.presetCursor = 0
	case types.KeySpace:
		m.toggleHeader()
		m.updateViewportContent()
	}

	m.updateViewportContent()
	return nil
}

func (m *Model) updateEdit(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if ok {

		switch keyMsg.String() {
		case types.KeyEscape:
			m.ExitEditMode()
			m.updateViewportContent()
			return nil
		case types.KeyTab:
			m.Headers[m.cursor].Key.Blur()
			m.Headers[m.cursor].Value.Blur()
			m.fieldFocus = (m.fieldFocus + 1) % 2
			if m.fieldFocus == 0 {
				return m.Headers[m.cursor].Key.Focus()
			}
			return m.Headers[m.cursor].Value.Focus()
		case types.KeyEnter:
			m.ExitEditMode()
			m.updateViewportContent()
			return nil
		}
	}

	var cmd tea.Cmd
	if m.fieldFocus == 0 {
		m.Headers[m.cursor].Key, cmd = m.Headers[m.cursor].Key.Update(msg)
	} else {
		m.Headers[m.cursor].Value, cmd = m.Headers[m.cursor].Value.Update(msg)
	}

	m.updateViewportContent()
	return cmd
}

func (m *Model) updatePresets(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch keyMsg.String() {
	case types.KeyJ, types.KeyDown:
		if m.presetCursor < len(CommonPresets)-1 {
			m.presetCursor++
		}
	case types.KeyK, types.KeyUp:
		if m.presetCursor > 0 {
			m.presetCursor--
		}
	case types.KeyEnter:
		p := CommonPresets[m.presetCursor]
		m.addHeader(p.Key, p.Value)
		m.showPresets = false
		m.updateViewportContent()
		m.ensureCursorVisible()
	case types.KeyEscape, types.KeyP:
		m.showPresets = false
	}

	return nil
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.showPresets {
		return m.updatePresets(msg)
	}

	if m.EditMode {
		return m.updateEdit(msg)
	}

	return m.updateNav(msg)
}

func (m Model) viewPresets(width int) string {
	content := "Select a preset:\n\n"

	for i, p := range CommonPresets {
		line := fmt.Sprintf("%s: %s", p.Key, p.Value)
		if i == m.presetCursor {
			line = lipgloss.NewStyle().Background(style.ColorSurface).Foreground(style.ColorText).Render("> " + line)
		} else {
			line = " " + line
		}

		content += line + "\n"
	}

	content += "\n" + style.Unselected.Render("[enter] select		[esc] cancel")
	return style.SectionBox("Headers - Presets", content, m.Focused, width)
}

func (m Model) renderHeaderLine(index int, h Header) string {
	isCursor := index == m.cursor && m.Focused

	check := "[ ]"
	if h.Enabled {
		check = "[X]"
	}

	var key, value string
	if m.EditMode && isCursor && m.fieldFocus == 0 {
		key = h.Key.View()
	} else {
		key = h.Key.Value()
	}

	if m.EditMode && isCursor && m.fieldFocus == 1 {
		value = h.Value.View()
	} else {
		value = h.Value.Value()
	}

	line := fmt.Sprintf("%s %s: %s", check, key, value)

	if h.Auto {
		line += style.Unselected.Render(" (auto)")
	}

	if !h.Enabled {
		line = style.Unselected.Render(line)
	} else if isCursor {
		line = lipgloss.NewStyle().Background(style.ColorSurface).Foreground(style.ColorText).Render(line)
	}

	return line
}

func (m *Model) updateViewportContent() {
	var content string
	if len(m.Headers) == 0 {
		content = style.Unselected.Render("No headers (press 'a' to add)")
	} else {
		for i, h := range m.Headers {
			line := m.renderHeaderLine(i, h)
			content += line + "\n"
		}
	}
	m.viewport.SetContent(content)
}

func (m *Model) ensureCursorVisible() {
	if m.cursor < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursor)
	} else if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(m.cursor - m.viewport.Height + 1)
	}
}

func (m Model) View(width int) string {
	if m.showPresets {
		return m.viewPresets(width)
	}

	topContent := m.viewport.View()
	footer := style.Unselected.Render("[a]dd [d]el [p]resets [space]toggle [tab]key<>value [esc/enter]validate")

	innerHeight := m.height - 4
	content := lipgloss.Place(width-6, innerHeight, lipgloss.Left, lipgloss.Bottom, footer, lipgloss.WithWhitespaceChars(" "), lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}))
	content = topContent + "\n" + content

	return style.SectionBox("Headers", content, m.Focused, width, m.height-4)
}
