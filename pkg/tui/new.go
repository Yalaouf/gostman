package tui

import (
	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func New() Model {
	t := textinput.New()
	t.Placeholder = "http://localhost:3000"
	t.Cursor.Blink = true
	t.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))

	methods := []request.HttpMethod{
		request.GET, request.POST, request.PUT,
		request.DELETE, request.PATCH, request.HEAD,
		request.TRACE, request.CONNECT,
	}

	v := viewport.New(80, 10)
	v.Style = viewportStyle

	return Model{
		urlInput:     t,
		methods:      methods,
		methodIndex:  0,
		focusSection: URL,
		responseView: v,
	}
}
