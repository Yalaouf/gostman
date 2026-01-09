package response

import (
	"fmt"
	"strings"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/utils"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
)

type Model struct {
	Viewport   viewport.Model
	Response   request.Response
	Focused    bool
	Error      string
	Loading    bool
	width      int
	height     int
	currentTab Tab
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
	m.updateViewportContent()
	m.Viewport.GotoTop()
}

func (m *Model) SetError(err string) {
	m.Error = err
	m.Response = request.Response{}
	m.Viewport.SetContent("")
}

func (m *Model) SetLoading(loading bool) {
	m.Loading = loading
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

func (m *Model) updateViewportContent() {
	if !m.HasResponse() {
		return
	}

	var content string
	switch m.currentTab {
	case TabPretty:
		body := m.Response.Body
		if utils.IsJSON(body) {
			body = utils.HighlightJSON(body)
		}

		if m.Viewport.Width > 0 {
			body = wrap.String(body, m.Viewport.Width-2)
		}

		content = body
	case TabRaw:
		body := m.Response.Body

		if m.Viewport.Width > 0 {
			body = wrap.String(body, m.Viewport.Width-2)
		}

		content = body
	case TabHeaders:
		var body string

		for key, values := range m.Response.Headers {
			for _, val := range values {
				body += fmt.Sprintf("%s: %s\n\n", key, val)
			}
		}

		if m.Viewport.Width > 0 {
			body = wrap.String(body, m.Viewport.Width-2)
		}

		content = body
	}

	padding := "\n\n"
	fullContent := fmt.Sprintf(
		"%s  •  %s\n\n%s%s",
		colorStatusCode(m.Response.StatusCode),
		colorTimeTaken(m.Response.TimeTaken),
		content,
		padding,
	)
	m.Viewport.SetContent(fullContent)
}

func (m *Model) NextTab() {
	i := int(m.currentTab)
	i = (i + 1) % len(AllTabs)
	m.currentTab = AllTabs[i]
	m.updateViewportContent()
}

func (m *Model) renderTabs() string {
	var tabs []string
	for _, t := range AllTabs {
		label := t.String()
		if t == m.currentTab {
			tabs = append(tabs, style.Selected.Render("["+label+"]"))
		} else {
			tabs = append(tabs, style.Unselected.Render(" "+label+" "))
		}
	}

	return strings.Join(tabs, "")
}

func (m Model) View(width int) string {
	borderColor := style.ColorGray
	if m.Focused {
		borderColor = style.ColorPurple
	}

	tabs := m.renderTabs()

	var content string
	if m.Loading {
		content = style.Unselected.Render("Loading...")
	} else if m.Error != "" {
		content = style.Error.Render("Error: " + m.Error)
	} else if m.HasResponse() {
		scrollbarStyle := lipgloss.NewStyle().Foreground(borderColor).MarginLeft(1)
		scrollbar := scrollbarStyle.Render(RenderScrollbar(m.Viewport))
		content = lipgloss.JoinHorizontal(lipgloss.Top, m.Viewport.View(), scrollbar)
	} else {
		content = style.Unselected.Render("No response yet. Press Alt+Enter to send a request.")
	}

	fullContent := tabs + "\n\n" + content

	return style.SectionBox("Response", fullContent, m.Focused, width, m.height+7)
}
