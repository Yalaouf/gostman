package requestmenu

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) moveDown() {
	maxIndex := m.getMaxIndex()
	if m.index < maxIndex {
		m.index++
	}
}

func (m *Model) moveUp() {
	if m.index > 0 {
		m.index--
	}
}

func (m *Model) getMaxIndex() int {
	switch m.viewMode {
	case ViewCollections, ViewMoveTarget:
		return len(m.collections)
	case ViewRequests:
		return len(m.requests) - 1
	}
	return 0
}

func (m *Model) enterCollection() {
	if m.index < len(m.collections) {
		coll := m.collections[m.index]
		m.selectedCollID = coll.ID
		m.selectedCollName = coll.Name
	} else {
		m.selectedCollID = ""
		m.selectedCollName = "Uncategorized"
	}
	m.viewMode = ViewRequests
	m.index = 0
	m.refresh()
}

func (m *Model) startCreateCollection() {
	m.inputMode = true
	m.inputAction = InputCreateCollection
	m.input.Placeholder = "Collection name"
	m.input.SetValue("")
	m.input.Focus()
	m.err = ""
}

func (m *Model) startRename() {
	m.inputMode = true
	m.err = ""

	switch m.viewMode {
	case ViewCollections:
		if m.index < len(m.collections) {
			m.inputAction = InputRenameCollection
			m.input.Placeholder = "New name"
			m.input.SetValue(m.collections[m.index].Name)
			m.input.Focus()
		}
	case ViewRequests:
		if len(m.requests) > 0 && m.index < len(m.requests) {
			m.inputAction = InputRenameRequest
			m.input.Placeholder = "New name"
			m.input.SetValue(m.requests[m.index].Name)
			m.input.Focus()
		}
	}
}

func (m *Model) confirmInput() tea.Cmd {
	value := strings.TrimSpace(m.input.Value())
	if value == "" {
		m.err = "Name cannot be empty"
		return nil
	}

	var err error

	switch m.inputAction {
	case InputCreateCollection:
		_, err = m.storage.CreateCollection(value)
	case InputRenameCollection:
		if m.index < len(m.collections) {
			_, err = m.storage.UpdateCollection(m.collections[m.index].ID, value)
		}
	case InputRenameRequest:
		if m.index < len(m.requests) {
			req := m.requests[m.index]
			req.Name = value
			err = m.storage.SaveRequest(req)
		}
	}

	if err != nil {
		m.err = err.Error()
		return nil
	}

	m.inputMode = false
	m.inputAction = InputNone
	m.input.Blur()
	m.refresh()
	return nil
}

func (m *Model) deleteSelected() {
	switch m.viewMode {
	case ViewCollections:
		if m.index < len(m.collections) {
			coll := m.collections[m.index]
			err := m.storage.DeleteCollection(coll.ID, false)
			if err != nil {
				m.err = err.Error()
				return
			}
			m.refresh()
			if m.index >= len(m.collections) && m.index > 0 {
				m.index--
			}
		}
	case ViewRequests:
		if len(m.requests) > 0 && m.index < len(m.requests) {
			req := m.requests[m.index]
			err := m.storage.DeleteRequest(req.ID)
			if err != nil {
				m.err = err.Error()
				return
			}
			m.refresh()
			if m.index >= len(m.requests) && m.index > 0 {
				m.index--
			}
		}
	}
	m.err = ""
}

func (m *Model) startMove() {
	if m.index < len(m.requests) {
		m.moveRequestID = m.requests[m.index].ID
		m.viewMode = ViewMoveTarget
		m.index = 0
		m.refresh()
	}
}

func (m *Model) confirmMove() tea.Cmd {
	var targetCollID string
	if m.index < len(m.collections) {
		targetCollID = m.collections[m.index].ID
	} else {
		targetCollID = ""
	}

	err := m.storage.MoveRequest(m.moveRequestID, targetCollID)
	if err != nil {
		m.err = err.Error()
		return nil
	}

	m.viewMode = ViewRequests
	m.moveRequestID = ""
	m.index = 0
	m.refresh()
	return nil
}
