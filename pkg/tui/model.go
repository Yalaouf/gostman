package tui

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/storage"
	"github.com/Yalaouf/gostman/pkg/tui/components/body"
	"github.com/Yalaouf/gostman/pkg/tui/components/headers"
	"github.com/Yalaouf/gostman/pkg/tui/components/method"
	"github.com/Yalaouf/gostman/pkg/tui/components/requestmenu"
	"github.com/Yalaouf/gostman/pkg/tui/components/response"
	"github.com/Yalaouf/gostman/pkg/tui/components/savepopup"
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
	showHelp bool

	focusSection types.FocusSection

	method   method.Model
	url      url.Model
	headers  headers.Model
	body     body.Model
	response response.Model

	storage     *storage.Storage
	savePopup   savepopup.Model
	requestMenu requestmenu.Model
}

func New(s *storage.Storage) Model {
	return Model{
		focusSection: types.FocusURL,
		method:       method.New(),
		url:          url.New(),
		headers:      headers.New(),
		body:         body.New(),
		response:     response.New(),
		storage:      s,
		savePopup:    savepopup.New(),
		requestMenu:  requestmenu.New(s),
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

	case requestmenu.LoadRequestMsg:
		return m.handleLoadRequest(msg), nil

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

	leftWidth := msg.Width / 2
	rightWidth := msg.Width - leftWidth - 4

	panelHeight := msg.Height - 8
	sectionHeight := panelHeight/3 + 1

	m.headers.SetSize(leftWidth, sectionHeight)
	m.body.SetSize(leftWidth, sectionHeight)
	m.response.SetSize(rightWidth, sectionHeight*2-3)
	return m
}

func (m Model) handleRequestComplete(msg requestMsg) Model {
	m.response.SetLoading(false)
	if msg.err != nil {
		m.response.SetError(msg.err.Error())
		return m
	}

	m.response.SetResponse(msg.response)
	return m
}

func (m Model) handleLoadRequest(msg requestmenu.LoadRequestMsg) Model {
	req := msg.Request
	if req == nil {
		return m
	}

	m.method.SetMethod(request.HTTPMethod(req.Method))
	m.url.SetValue(req.URL)
	m.headers.SetHeaders(req.Headers)
	m.body.SetValue(req.Body)

	switch req.BodyType {
	case "json":
		m.body.SetType(body.TypeJSON)
	case "form-data":
		m.body.SetType(body.TypeFormData)
	case "urlencoded":
		m.body.SetType(body.TypeURLEncoded)
	default:
		m.body.SetType(body.TypeNone)
	}

	m.syncContentType()
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

	req.SetURL(strings.TrimSpace(m.url.Value()))
	req.SetMethod(m.method.Selected())
	req.SetBody(m.body.Value())
	req.SetTimeout(request.DefaultTimeout)

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
		req.AddHeader(strings.TrimSpace(key), strings.TrimSpace(value))
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
