package tui

import (
	"github.com/Yalaouf/gostman/pkg/request"
	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() Model {
	return Model{
		Method: request.GET,
		URL:    "",
		Body:   "",
		Header: make(map[string]string),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
