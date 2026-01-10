package method

import "github.com/Yalaouf/gostman/pkg/request"

type Model struct {
	Methods []request.HTTPMethod
	Index   int
	Focused bool
}

func New() Model {
	return Model{
		Methods: []request.HTTPMethod{
			request.GET, request.POST, request.PUT,
			request.DELETE, request.PATCH, request.HEAD,
			request.OPTIONS, request.TRACE, request.CONNECT,
		},
		Index:   0,
		Focused: false,
	}
}

func (m Model) Selected() request.HTTPMethod {
	return m.Methods[m.Index]
}

func (m *Model) Next() {
	m.Index = (m.Index + 1) % len(m.Methods)
}

func (m *Model) Previous() {
	m.Index = (m.Index - 1 + len(m.Methods)) % len(m.Methods)
}

func (m *Model) SetMethod(method request.HTTPMethod) {
	for i, meth := range m.Methods {
		if meth == method {
			m.Index = i
			return
		}
	}
}

func (m *Model) Focus() {
	m.Focused = true
}

func (m *Model) Blur() {
	m.Focused = false
}
