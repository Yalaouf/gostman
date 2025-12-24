package tui

func (m Model) View() string {
	s := "Gostman - A simple API testing tool\n\n"

	s += "Method: " + string(m.Method) + "\n"
	s += "URL: " + m.URL + "\n"
	s += "Body: " + m.Body + "\n"
	s += "Headers:\n"
	for k, v := range m.Header {
		s += "  " + k + ": " + v + "\n"
	}

	s += "\nPress q to quit.\n"

	return s
}
