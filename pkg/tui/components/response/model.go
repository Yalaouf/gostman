package response

import (
	"fmt"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/tui/utils"
	"github.com/charmbracelet/bubbles/viewport"
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
	fullscreen bool
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
	m.Viewport.Height = height - 8
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

func (m *Model) IsFullscreen() bool {
	return m.fullscreen
}

func (m *Model) ToggleFullscreen() {
	m.fullscreen = !m.fullscreen
}

func (m *Model) ExitFullscreen() {
	m.fullscreen = false
}

func (m *Model) NextTab() {
	i := int(m.currentTab)
	i = (i + 1) % len(AllTabs)
	m.currentTab = AllTabs[i]
	m.updateViewportContent()
}

func (m *Model) GetContent() string {
	if !m.HasResponse() {
		return ""
	}

	switch m.currentTab {
	case TabPretty, TabRaw:
		return m.Response.Body
	case TabHeaders:
		var content string
		for key, values := range m.Response.Headers {
			for _, val := range values {
				content += fmt.Sprintf("%s: %s\n", key, val)
			}
		}
		return content
	}
	return ""
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
