package tui

import (
	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/charmbracelet/bubbles/textinput"
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

	urlInput textinput.Model

	methods     []request.HttpMethod
	methodIndex int

	focusSection uint

	req request.Model
	res request.Response
}

const (
	METHOD = iota
	URL
	HEADERS
	BODY
)
