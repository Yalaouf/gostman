package requestmenu

import (
	"fmt"
	"strings"

	"github.com/Yalaouf/gostman/pkg/storage"
	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.inputMode {
		return m.handleInputMode(msg)
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	key := keyMsg.String()

	switch key {
	case "esc":
		return m.handleEscape()
	case "j", "down":
		m.moveDown()
	case "k", "up":
		m.moveUp()
	case "enter":
		return m.handleEnter()
	case "n":
		if m.viewMode == ViewCollections {
			m.startCreateCollection()
			return textinput.Blink
		}
	case "r":
		m.startRename()
		return textinput.Blink
	case "d":
		m.deleteSelected()
	case "m":
		if m.viewMode == ViewRequests && len(m.requests) > 0 {
			m.startMove()
		}
	}

	return nil
}

func (m *Model) handleInputMode(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return cmd
	}

	key := keyMsg.String()

	switch key {
	case "esc":
		m.inputMode = false
		m.inputAction = InputNone
		m.input.Blur()
		return nil
	case "enter":
		return m.confirmInput()
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m *Model) handleEscape() tea.Cmd {
	switch m.viewMode {
	case ViewCollections:
		m.Hide()
	case ViewRequests:
		m.viewMode = ViewCollections
		m.index = 0
		m.refresh()
	case ViewMoveTarget:
		m.viewMode = ViewRequests
		m.moveRequestID = ""
		m.index = 0
	}
	return nil
}

func (m *Model) handleEnter() tea.Cmd {
	switch m.viewMode {
	case ViewCollections:
		m.enterCollection()
	case ViewRequests:
		if len(m.requests) > 0 && m.index < len(m.requests) {
			req := m.requests[m.index]
			m.Hide()
			return func() tea.Msg {
				return LoadRequestMsg{Request: req}
			}
		}
	case ViewMoveTarget:
		return m.confirmMove()
	}
	return nil
}

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

func (m Model) View() string {
	if m.inputMode {
		return m.viewInput()
	}

	switch m.viewMode {
	case ViewCollections:
		return m.viewCollections()
	case ViewRequests:
		return m.viewRequests()
	case ViewMoveTarget:
		return m.viewMoveTarget()
	}

	return ""
}

func (m Model) viewCollections() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	title := titleStyle.Render("Saved Requests")

	var b strings.Builder
	b.WriteString("Collections:\n\n")

	for i, coll := range m.collections {
		count := len(m.storage.ListRequestsByCollection(coll.ID))
		line := fmt.Sprintf("%s (%d)", coll.Name, count)
		if i == m.index {
			b.WriteString(style.Selected.Render("▸ " + line))
		} else {
			b.WriteString(style.Unselected.Render("  " + line))
		}
		b.WriteString("\n")
	}

	uncatCount := len(m.storage.ListRequestsByCollection(""))
	uncatLine := fmt.Sprintf("Uncategorized (%d)", uncatCount)
	if m.index == len(m.collections) {
		b.WriteString(style.Selected.Render("▸ " + uncatLine))
	} else {
		b.WriteString(style.Unselected.Render("  " + uncatLine))
	}

	var errView string
	if m.err != "" {
		errView = "\n\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("[enter]open [n]ew [r]ename [d]elete [esc]close")

	content := title + "\n\n" + b.String() + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Width(50).
		Render(content)

	return box
}

func (m Model) viewRequests() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected
	methodStyle := lipgloss.NewStyle().Foreground(style.ColorBlue).Width(8)

	title := titleStyle.Render(m.selectedCollName)

	var b strings.Builder

	if len(m.requests) == 0 {
		b.WriteString(style.Unselected.Render("  No requests in this collection"))
	} else {
		for i, req := range m.requests {
			method := methodStyle.Render(req.Method)
			line := fmt.Sprintf("%s %s", method, req.Name)
			if i == m.index {
				b.WriteString(style.Selected.Render("▸ ") + line)
			} else {
				b.WriteString(style.Unselected.Render("  ") + line)
			}
			b.WriteString("\n")
		}
	}

	var errView string
	if m.err != "" {
		errView = "\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("[enter]load [r]ename [d]elete [m]ove [esc]back")

	content := title + "\n\n" + b.String() + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Width(50).
		Render(content)

	return box
}

func (m Model) viewMoveTarget() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	title := titleStyle.Render("Move to Collection")

	var b strings.Builder

	for i, coll := range m.collections {
		if i == m.index {
			b.WriteString(style.Selected.Render("▸ " + coll.Name))
		} else {
			b.WriteString(style.Unselected.Render("  " + coll.Name))
		}
		b.WriteString("\n")
	}

	uncatLine := "Uncategorized"
	if m.index == len(m.collections) {
		b.WriteString(style.Selected.Render("▸ " + uncatLine))
	} else {
		b.WriteString(style.Unselected.Render("  " + uncatLine))
	}

	var errView string
	if m.err != "" {
		errView = "\n\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("[enter]confirm [esc]cancel")

	content := title + "\n\n" + b.String() + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Width(50).
		Render(content)

	return box
}

func (m Model) viewInput() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(style.ColorOrange)
	hintStyle := style.Unselected

	var title string
	switch m.inputAction {
	case InputCreateCollection:
		title = titleStyle.Render("New Collection")
	case InputRenameCollection:
		title = titleStyle.Render("Rename Collection")
	case InputRenameRequest:
		title = titleStyle.Render("Rename Request")
	}

	inputView := m.input.View()

	var errView string
	if m.err != "" {
		errView = "\n" + style.Error.Render(m.err)
	}

	hint := hintStyle.Render("Enter to confirm, Esc to cancel")

	content := title + "\n\n" + inputView + errView + "\n\n" + hint

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.ColorPurple).
		Padding(1, 3).
		Render(content)

	return box
}
