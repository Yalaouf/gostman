package tui

import (
	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss/list"
)

type requestMsg struct {
	response request.Response
	err      error
}

type Model struct {
	urlInput     textinput.Model
	methodsList  []list.Items
	focusSection uint
	req          request.Model
	loading      bool
	response     request.Response
	errorMsg     string
}

const (
	METHOD = iota
	URL
	HEADERS
	BODY
)
