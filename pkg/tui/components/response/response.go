package response

import (
	"fmt"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Viewport viewport.Model
	Response request.Response
	Focused  bool
}

func New() Model {
	vp := viewport.New(80, 10)
	vp.Style = style.Viewport

	return Model{
		Viewport: vp,
		Focused:  false,
	}
}

func (m Model) HasResponse() bool {
	return m.Response.StatusCode != 0
}

func (m *Model) SetSize(width, height int) {
	m.Viewport.Width = width
	m.Viewport.Height = height
}

func (m *Model) SetResponse(res request.Response) {
	m.Response = res
	content := fmt.Sprintf(
		"Status %d  •  Time: %dms\n\n%s",
		res.StatusCode,
		res.TimeTaken,
		res.Body,
	)
	m.Viewport.SetContent(content)
}

func (m *Model) Focus() {
	m.Focused = true
}

func (m *Model) Blur() {
	m.Focused = false
}

func (m *Model) ScrollDown(n int) {
	m.Viewport.ScrollDown(n)
}

func (m *Model) ScrollUp(n int) {
	m.Viewport.ScrollUp(n)
}

func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

func (m Model) View() string {
	header := style.SectionTitle.Render("Response")

	borderColor := style.ColorGray
	if m.Focused {
		borderColor = style.ColorPurple
	}

	scrollbarStyle := lipgloss.NewStyle().Foreground(borderColor).MarginLeft(1)
	scrollbar := scrollbarStyle.Render(RenderScrollbar(m.Viewport))

	body := lipgloss.JoinHorizontal(lipgloss.Top, m.Viewport.View(), scrollbar)

	return header + "\n" + body
}
