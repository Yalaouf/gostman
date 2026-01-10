package headers

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Header struct {
	Key     textinput.Model
	Value   textinput.Model
	Enabled bool
	Auto    bool
}

type Model struct {
	Headers      []Header
	cursor       int
	fieldFocus   int
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

func (m *Model) SetHeaders(headers map[string]string) {
	m.Headers = []Header{}
	for key, value := range headers {
		h := newHeader(key, value, false)
		m.Headers = append(m.Headers, h)
	}
	m.cursor = 0
	m.updateViewportContent()
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

func (m *Model) ensureCursorVisible() {
	if m.cursor < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursor)
	} else if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(m.cursor - m.viewport.Height + 1)
	}
}
