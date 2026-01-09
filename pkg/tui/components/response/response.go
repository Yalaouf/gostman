package response

import (
	"fmt"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/utils"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
)

type Model struct {
	Viewport viewport.Model
	Response request.Response
	Focused  bool
	Error    string
	width    int
	height   int
}

func New() Model {
	vp := viewport.New(80, 10)

	return Model{
		Viewport: vp,
		Focused:  false,
	}
}

func (m Model) HasResponse() bool {
	return m.Response.StatusCode != 0
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.Viewport.Width = width - 6
	m.Viewport.Height = height + 6
}

func (m *Model) SetResponse(res request.Response) {
	m.Response = res
	m.Error = ""

	body := res.Body
	if utils.IsJSON(res.Body) {
		body = utils.HighlightJSON(res.Body)
	}

	if m.Viewport.Width > 0 {
		body = wrap.String(body, m.Viewport.Width-2)
	}

	padding := "\n\n"

	content := fmt.Sprintf(
		"%s  •  %s\n\n%s%s",
		colorStatusCode(res.StatusCode),
		colorTimeTaken(res.TimeTaken),
		body,
		padding,
	)

	m.Viewport.SetContent(content)
	m.Viewport.GotoTop()
}

func (m *Model) SetError(err string) {
	m.Error = err
	m.Response = request.Response{}
	m.Viewport.SetContent("")
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

func (m Model) View(width int) string {
	borderColor := style.ColorGray
	if m.Focused {
		borderColor = style.ColorPurple
	}

	var content string
	if m.Error != "" {
		content = style.Error.Render("Error: " + m.Error)
	} else if m.HasResponse() {
		scrollbarStyle := lipgloss.NewStyle().Foreground(borderColor).MarginLeft(1)
		scrollbar := scrollbarStyle.Render(RenderScrollbar(m.Viewport))
		content = lipgloss.JoinHorizontal(lipgloss.Top, m.Viewport.View(), scrollbar)
	} else {
		content = style.Unselected.Render("No response yet. Press Alt+Enter to send a request.")
	}

	return style.SectionBox("Response", content, m.Focused, width, m.height+7)
}
