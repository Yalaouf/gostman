package requestmenu

import (
	"github.com/Yalaouf/gostman/pkg/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewMode uint

const (
	ViewCollections ViewMode = iota
	ViewRequests
	ViewMoveTarget
)

type InputAction uint

const (
	InputNone InputAction = iota
	InputCreateCollection
	InputRenameCollection
	InputRenameRequest
)

type LoadRequestMsg struct {
	Request *storage.Request
}

type Model struct {
	visible  bool
	viewMode ViewMode
	index    int

	collections []*storage.Collection
	requests    []*storage.Request

	selectedCollID   string
	selectedCollName string

	moveRequestID string

	inputMode   bool
	inputAction InputAction
	input       textinput.Model
	err         string

	storage *storage.Storage
}

func New(s *storage.Storage) Model {
	ti := textinput.New()
	ti.CharLimit = 64
	ti.Width = 30

	return Model{
		storage: s,
		input:   ti,
	}
}

func (m *Model) Show() tea.Cmd {
	m.visible = true
	m.viewMode = ViewCollections
	m.index = 0
	m.err = ""
	m.inputMode = false
	m.refresh()
	return nil
}

func (m *Model) Hide() {
	m.visible = false
	m.inputMode = false
	m.input.Blur()
}

func (m Model) Visible() bool {
	return m.visible
}

func (m *Model) refresh() {
	m.collections = m.storage.ListCollections()
	if m.viewMode == ViewRequests || m.viewMode == ViewMoveTarget {
		m.requests = m.storage.ListRequestsByCollection(m.selectedCollID)
	}
}
