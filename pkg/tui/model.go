package tui

import (
	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/tui/components/body"
	"github.com/Yalaouf/gostman/pkg/tui/components/headers"
	"github.com/Yalaouf/gostman/pkg/tui/components/method"
	"github.com/Yalaouf/gostman/pkg/tui/components/response"
	"github.com/Yalaouf/gostman/pkg/tui/components/url"
	"github.com/Yalaouf/gostman/pkg/tui/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

func (m *Model) syncContentType() {
	var contentType string

	switch m.body.BodyType {
	case body.TypeJSON:
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

func (m Model) buildRequestModel() *request.Model {
	req := request.NewModel()

	req.SetURL(m.url.Value())
	req.SetMethod(m.method.Selected())
	req.SetBody(m.body.Value())
	req.SetTimeout(30000)

	switch m.body.BodyType {
	case body.TypeJSON:
		req.SetBodyType(request.BodyTypeJSON)
	case body.TypeFormData:
		req.SetBodyType(request.BodyTypeFormData)
	case body.TypeURLEncoded:
		req.SetBodyType(request.BodyTypeURLEncoded)
	default:
		req.SetBodyType(request.BodyTypeNone)
	}

	for key, value := range m.headers.EnabledHeaders() {
		req.AddHeader(key, value)
	}

	return req
}

func (m Model) sendRequest() tea.Cmd {
	return func() tea.Msg {
		req := m.buildRequestModel()

		res, err := request.SendRequest(req)
		if err != nil {
			return requestMsg{err: err}
		}

		return requestMsg{response: *res}
	}
}
