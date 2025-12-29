package tui

import (
	"github.com/Yalaouf/gostman/pkg/request"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) sendRequest() tea.Cmd {
	return func() tea.Msg {
		m.req.SetTimeout(30000)
		res, err := request.SendRequest(&m.req)
		if err != nil {
			return requestMsg{err: err}
		}

		return requestMsg{response: *res}
	}
}
