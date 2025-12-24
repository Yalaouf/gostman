package tui

import "github.com/Yalaouf/gostman/pkg/request"

type Model struct {
	Method request.HttpMethod
	URL    string
	Body   string
	Header map[string]string
}
