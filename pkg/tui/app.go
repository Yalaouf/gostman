package tui

import (
	"fmt"
	"os"

	"github.com/Yalaouf/gostman/pkg/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func Gostman() {
	s, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	m := New(s)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
