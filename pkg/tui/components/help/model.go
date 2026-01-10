package help

import "github.com/charmbracelet/bubbles/viewport"

type Model struct {
	viewport viewport.Model
	ready    bool
}

func New() Model {
	return Model{}
}

func (m *Model) SetSize(width, height int) {
	contentWidth := 40
	contentHeight := height - 10

	if !m.ready {
		m.viewport = viewport.New(contentWidth, contentHeight)
		m.ready = true
	} else {
		m.viewport.Width = contentWidth
		m.viewport.Height = contentHeight
	}

	m.viewport.SetContent(renderContent())
}

func (m *Model) ScrollDown() {
	m.viewport.ScrollDown(1)
}

func (m *Model) ScrollUp() {
	m.viewport.ScrollUp(1)
}
