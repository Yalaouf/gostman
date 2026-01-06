package tui

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/tui/components/body"
	"github.com/Yalaouf/gostman/pkg/tui/components/headers"
	"github.com/Yalaouf/gostman/pkg/tui/components/method"
	"github.com/Yalaouf/gostman/pkg/tui/components/response"
	"github.com/Yalaouf/gostman/pkg/tui/components/url"
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/Yalaouf/gostman/pkg/tui/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type requestMsg struct {
	response request.Response
	err      error
}

type Model struct {
	width  int
	height int

	loading  bool
	errorMsg string
	showHelp bool

	focusSection types.FocusSection

	method   method.Model
	url      url.Model
	headers  headers.Model
	body     body.Model
	response response.Model

	req request.Model
}

func New() Model {
	return Model{
		focusSection: types.FocusURL,
		method:       method.New(),
		url:          url.New(),
		headers:      headers.New(),
		body:         body.New(),
		response:     response.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg), nil

	case requestMsg:
		return m.handleRequestComplete(msg), nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	if m.focusSection == types.FocusURL {
		return m.handleURLInput(msg)
	}

	if m.focusSection == types.FocusBody {
		return m.handleBodyInput(msg)
	}

	if m.focusSection == types.FocusHeaders {
		return m.handleHeadersInput(msg)
	}

	return m, nil
}

func (m Model) View() string {
	if m.showHelp {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.viewHelp())
	}

	if m.focusSection == types.FocusMethod {
		picker := m.method.View()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, picker)
	}

	top := m.url.ViewWithMethod(m.method.ViewSelected(), m.width-2)

	headersView := m.headers.View(m.width / 2)
	bodyView := m.body.View(m.width / 2)
	middle := lipgloss.JoinHorizontal(lipgloss.Top, headersView, bodyView)

	var bottom string
	if m.response.HasResponse() {
		bottom = m.response.View()
	}

	var statusLeft string
	if m.loading {
		statusLeft = "Loading..."
	}

	if m.errorMsg != "" {
		statusLeft = style.Error.Render("Error: ", m.errorMsg)
	}

	helpHint := style.Unselected.Render("Press ? for help")
	status := lipgloss.NewStyle().Width(m.width - 2).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, statusLeft, lipgloss.PlaceHorizontal(m.width-2-lipgloss.Width(statusLeft), lipgloss.Right, helpHint)),
	)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.displayTitle(),
		top,
		middle,
		bottom,
	)

	contentHeight := lipgloss.Height(mainContent)
	spacerHeight := m.height - contentHeight - 1
	spacer := ""
	if spacerHeight > 0 {
		spacer = strings.Repeat("\n", spacerHeight-1)
	}

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, spacer, status)
}

func (m Model) displayTitle() string {
	title := style.Title.Render("GOSTMAN")
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title)
}

func (m Model) viewHelp() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	keyStyle := lipgloss.NewStyle().Foreground(style.ColorGreen)
	descStyle := lipgloss.NewStyle().Foreground(style.ColorText)

	title := titleStyle.Render("Keyboard Shortcuts")

	sections := []string{
		titleStyle.Render("Navigation"),
		keyStyle.Render("  i") + descStyle.Render("        Focus URL input"),
		keyStyle.Render("  m") + descStyle.Render("        Focus method selector"),
		keyStyle.Render("  h") + descStyle.Render("        Focus headers"),
		keyStyle.Render("  b") + descStyle.Render("        Focus body"),
		keyStyle.Render("  r") + descStyle.Render("        Focus response"),
		"",
		titleStyle.Render("Actions"),
		keyStyle.Render("  Alt+Enter") + descStyle.Render(" Send request"),
		keyStyle.Render("  Enter") + descStyle.Render("     Enter edit mode"),
		keyStyle.Render("  Esc") + descStyle.Render("       Exit edit mode"),
		keyStyle.Render("  q") + descStyle.Render("         Quit"),
		"",
		titleStyle.Render("Headers"),
		keyStyle.Render("  a") + descStyle.Render("         Add new header"),
		keyStyle.Render("  d") + descStyle.Render("         Delete header"),
		keyStyle.Render("  p") + descStyle.Render("         Open presets"),
		keyStyle.Render("  Space") + descStyle.Render("     Toggle header"),
		keyStyle.Render("  j/k") + descStyle.Render("       Navigate up/down"),
		"",
		titleStyle.Render("Body"),
		keyStyle.Render("  Tab") + descStyle.Render("       Cycle body type"),
		"",
		titleStyle.Render("Response"),
		keyStyle.Render("  j/k") + descStyle.Render("       Scroll up/down"),
		keyStyle.Render("  g/G") + descStyle.Render("       Go to top/bottom"),
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Render(title + "\n\n" + content + "\n\n" + style.Unselected.Render("Press ? or Esc to close"))

	return box
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) Model {
	m.width = msg.Width
	m.height = msg.Height
	m.url.SetWidth(msg.Width - 10)
	m.headers.SetSize(msg.Width/2, msg.Height/6)
	m.body.SetSize(msg.Width/2, msg.Height/5)
	m.response.SetSize(msg.Width-4, msg.Height/2-1)
	return m
}

func (m Model) handleRequestComplete(msg requestMsg) Model {
	m.loading = false

	if msg.err != nil {
		m.errorMsg = msg.err.Error()
		return m
	}

	m.response.SetResponse(msg.response)
	return m
}

func (m Model) handleFocusChange(section types.FocusSection) (Model, tea.Cmd) {
	m.focusSection = section

	m.method.Blur()
	m.url.Blur()
	m.headers.Blur()
	m.body.Blur()
	m.response.Blur()

	switch section {
	case types.FocusMethod:
		m.method.Focus()
		return m, nil
	case types.FocusURL:
		m.url.Focused = true
		return m, m.url.Focus()
	case types.FocusHeaders:
		m.headers.Focus()
		return m, nil
	case types.FocusBody:
		return m, m.body.Focus()
	case types.FocusResult:
		m.response.Focus()
		return m, nil
	}

	return m, nil
}

func (m Model) handleNavigation(key string) Model {
	switch m.focusSection {
	case types.FocusMethod:
		if key == types.KeyJ || key == types.KeyDown {
			m.method.Next()
		} else {
			m.method.Previous()
		}
		m.req.SetMethod(m.method.Selected())
	case types.FocusBody:
		if key == types.KeyTab {
			m.body.NextType()
			m.syncContentType()
		}
	case types.FocusResult:
		if key == types.KeyJ || key == types.KeyDown {
			m.response.ScrollDown(1)
		} else {
			m.response.ScrollUp(1)
		}
	}
	return m
}

func (m Model) handleScroll(key string) Model {
	if m.focusSection != types.FocusResult {
		return m
	}

	if key == types.KeyG {
		m.response.GotoTop()
	} else {
		m.response.GotoBottom()
	}
	return m
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focusSection {
	case types.FocusMethod:
		m.req.SetMethod(m.method.Selected())
		return m.handleFocusChange(types.FocusURL)
	case types.FocusBody:
		return m, m.body.EnterEditMode()
	case types.FocusHeaders:
		return m, m.headers.EnterEditMode()
	}

	return m, nil
}

func (m Model) handleEscape() (Model, tea.Cmd) {
	switch m.focusSection {
	case types.FocusMethod:
		m.focusSection = types.FocusURL
	case types.FocusBody:
		if m.body.IsFocused() {
			m.body.ExitEditMode()
			return m, nil
		}
	case types.FocusHeaders:
		if m.headers.IsFocused() {
			m.headers.ExitEditMode()
			return m, nil
		}
	}
	m.url.Blur()
	return m, nil
}

func (m Model) handleURLInput(msg tea.Msg) (Model, tea.Cmd) {
	cmd := m.url.Update(msg)
	m.req.SetURL(m.url.Value())
	return m, cmd
}

func (m Model) handleBodyInput(msg tea.Msg) (Model, tea.Cmd) {
	cmd := m.body.Update(msg)
	m.req.SetBody(m.body.Value())
	return m, cmd
}

func (m Model) handleHeadersInput(msg tea.Msg) (Model, tea.Cmd) {
	cmd := m.headers.Update(msg)
	return m, cmd
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == types.KeyQuestion {
		m.showHelp = !m.showHelp
		return m, nil
	}

	if m.showHelp {
		if key == types.KeyEscape {
			m.showHelp = false
			return m, nil
		}
		return m, nil
	}

	if key == types.KeyAltEnter {
		m.loading = true
		m.errorMsg = ""
		return m, m.sendRequest()
	}

	if key == types.KeyEscape {
		return m.handleEscape()
	}

	if key == types.KeyCtrlC {
		return m, tea.Quit
	}

	if m.url.IsFocused() {
		return m.handleURLInput(msg)
	}

	if m.body.IsFocused() {
		return m.handleBodyInput(msg)
	}

	if m.headers.EditMode {
		return m.handleHeadersInput(msg)
	}

	if key == types.KeyQ {
		return m, tea.Quit
	}

	if m.focusSection == types.FocusBody {
		switch key {
		case types.KeyTab:
			return m.handleNavigation(key), nil
		}
	}

	if m.focusSection == types.FocusHeaders {
		switch key {
		case types.KeyJ, types.KeyK, types.KeyUp, types.KeyDown, types.KeyEnter, types.KeyTab:
			return m.handleHeadersInput(msg)
		}

		if key == types.KeyA || key == types.KeyD || key == types.KeyP || key == types.KeyEnter || key == types.KeySpace {
			return m.handleHeadersInput(msg)
		}
	}

	switch key {
	case types.KeyM:
		return m.handleFocusChange(types.FocusMethod)
	case types.KeyU:
		return m.handleFocusChange(types.FocusURL)
	case types.KeyH:
		return m.handleFocusChange(types.FocusHeaders)
	case types.KeyB:
		return m.handleFocusChange(types.FocusBody)
	case types.KeyR:
		if m.response.HasResponse() {
			return m.handleFocusChange(types.FocusResult)
		}
		return m, nil
	}

	switch key {
	case types.KeyJ, types.KeyDown, types.KeyK, types.KeyUp:
		return m.handleNavigation(key), nil
	}

	switch key {
	case types.KeyG, types.KeyShiftG:
		return m.handleScroll(key), nil
	}

	if key == types.KeyEnter {
		return m.handleEnter()
	}

	return m, nil
}

func (m *Model) syncContentType() {
	var contentType string

	switch m.body.BodyType {
	case body.TypeRaw:
		contentType = "application/json"
	case body.TypeFormData:
		contentType = "multipart/form-data"
	case body.TypeURLEncoded:
		contentType = "application/x-www-form-urlencoded"
	case body.TypeNone:
		contentType = ""
	}

	m.headers.SetContentType(contentType)
}

func (m Model) sendRequest() tea.Cmd {
	return func() tea.Msg {
		m.req.SetTimeout(30000)

		for key, value := range m.headers.EnabledHeaders() {
			m.req.AddHeader(key, value)
		}

		res, err := request.SendRequest(&m.req)
		if err != nil {
			return requestMsg{err: err}
		}

		return requestMsg{response: *res}
	}
}
