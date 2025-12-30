package method

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/request"
	"github.com/Yalaouf/gostman/pkg/tui/style"
)

type Model struct {
	Methods []request.HttpMethod
	Index   int
	Focused bool
}

func New() Model {
	return Model{
		Methods: []request.HttpMethod{
			request.GET, request.POST, request.PUT,
			request.DELETE, request.PATCH, request.HEAD,
			request.TRACE, request.CONNECT,
		},
		Index:   0,
		Focused: false,
	}
}

func (m Model) Selected() request.HttpMethod {
	return m.Methods[m.Index]
}

func (m *Model) Next() {
	m.Index = (m.Index + 1) % len(m.Methods)
}

func (m *Model) Previous() {
	m.Index = (m.Index - 1 + len(m.Methods)) % len(m.Methods)
}

func (m *Model) Focus() {
	m.Focused = true
}

func (m *Model) Blur() {
	m.Focused = false
}

func (m Model) View() string {
	var b strings.Builder

	for i, method := range m.Methods {
		if i == m.Index {
			b.WriteString(style.Selected.Render("▸ " + string(method)))
		} else {
			b.WriteString(style.Unselected.Render("  " + string(method)))
		}
		b.WriteString("\n")
	}

	content := b.String()
	if m.Focused {
		return style.FocusedBorder.Render(content) + "\n"
	}

	return content + "\n"
}

func (m Model) ViewSelected() string {
	return style.Selected.Render(string(m.Selected()))
}
