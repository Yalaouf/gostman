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

	focusSection FocusSection

	method   method.Model
	url      url.Model
	headers  headers.Model
	body     body.Model
	response response.Model

	req request.Model
}

func New() Model {
	return Model{
		focusSection: FocusURL,
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

	if m.focusSection == FocusURL {
		return m.handleURLInput(msg)
	}

	if m.focusSection == FocusBody {
		return m.handleBodyInput(msg)
	}

	return m, nil
}

func (m Model) View() string {
	// Fullscreen method picker
	if m.focusSection == FocusMethod {
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

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) Model {
	m.width = msg.Width
	m.height = msg.Height
	m.url.SetWidth(msg.Width - 10)
	m.body.SetSize(msg.Width/2, msg.Height/5)
	m.response.SetSize(msg.Width-4, msg.Height/2)
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

func (m Model) handleFocusChange(section FocusSection) (Model, tea.Cmd) {
	m.focusSection = section

	m.method.Blur()
	m.url.Blur()
	m.headers.Blur()
	m.body.Blur()
	m.response.Blur()

	switch section {
	case FocusMethod:
		m.method.Focus()
		return m, nil
	case FocusURL:
		m.url.Focused = true
		return m, m.url.Focus()
	case FocusHeaders:
		m.headers.Focus()
		return m, nil
	case FocusBody:
		return m, m.body.Focus()
	case FocusResult:
		m.response.Focus()
		return m, nil
	}

	return m, nil
}

func (m Model) handleNavigation(key string) Model {
	switch m.focusSection {
	case FocusMethod:
		if key == KeyJ || key == KeyDown {
			m.method.Next()
		} else {
			m.method.Previous()
		}
		m.req.SetMethod(m.method.Selected())
	case FocusBody:
		if key == KeyTab {
			m.body.NextType()
		}
	case FocusResult:
		if key == KeyJ || key == KeyDown {
			m.response.ScrollDown(1)
		} else {
			m.response.ScrollUp(1)
		}
	}
	return m
}

func (m Model) handleScroll(key string) Model {
	if m.focusSection != FocusResult {
		return m
	}

	if key == KeyG {
		m.response.GotoTop()
	} else {
		m.response.GotoBottom()
	}
	return m
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focusSection {
	case FocusMethod:
		m.req.SetMethod(m.method.Selected())
		return m.handleFocusChange(FocusURL)
	case FocusBody:
		return m, m.body.EnterEditMode()
	}

	return m, nil
}

func (m Model) handleEscape() (Model, tea.Cmd) {
	switch m.focusSection {
	case FocusMethod:
		m.focusSection = FocusURL
	case FocusBody:
		if m.body.IsFocused() {
			m.body.ExitEditMode()
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

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == KeyAltEnter {
		m.loading = true
		m.errorMsg = ""
		return m, m.sendRequest()
	}

	if key == KeyEscape {
		return m.handleEscape()
	}

	if key == KeyCtrlC {
		return m, tea.Quit
	}

	if m.url.IsFocused() {
		return m.handleURLInput(msg)
	}

	if m.body.IsFocused() {
		return m.handleBodyInput(msg)
	}

	if key == KeyQuit {
		return m, tea.Quit
	}

	if m.focusSection == FocusBody {
		switch key {
		case KeyH, KeyL, KeyLeft, KeyRight, KeyTab:
			return m.handleNavigation(key), nil
		}
	}

	switch key {
	case KeyMethod:
		return m.handleFocusChange(FocusMethod)
	case KeyInput:
		return m.handleFocusChange(FocusURL)
	case KeyHeaders:
		return m.handleFocusChange(FocusHeaders)
	case KeyBody:
		return m.handleFocusChange(FocusBody)
	case KeyResult:
		if m.response.HasResponse() {
			return m.handleFocusChange(FocusResult)
		}
		return m, nil
	}

	switch key {
	case KeyJ, KeyDown, KeyK, KeyUp:
		return m.handleNavigation(key), nil
	}

	switch key {
	case KeyG, KeyShiftG:
		return m.handleScroll(key), nil
	}

	if key == KeyEnter {
		return m.handleEnter()
	}

	return m, nil
}

func (m Model) sendRequest() tea.Cmd {
	return func() tea.Msg {
		m.req.SetTimeout(30000)

		if contentType := m.body.Type().ContentType(); contentType != "" {
			m.req.AddHeader("Content-Type", contentType)
		}

		res, err := request.SendRequest(&m.req)
		if err != nil {
			return requestMsg{err: err}
		}

		return requestMsg{response: *res}
	}
}
