package help

import "github.com/Yalaouf/gostman/pkg/tui/utils"

type Section struct {
	Title string
	Keys  []KeyBinding
}

type KeyBinding struct {
	Key  string
	Desc string
}

func GetSections() []Section {
	return []Section{
		{
			Title: "Navigation",
			Keys: []KeyBinding{
				{Key: "u", Desc: "Focus URL input"},
				{Key: "m", Desc: "Focus method selector"},
				{Key: "h", Desc: "Focus headers"},
				{Key: "b", Desc: "Focus body"},
				{Key: "r", Desc: "Focus response"},
				{Key: "j/k", Desc: "Navigate up/down"},
			},
		},
		{
			Title: "Actions",
			Keys: []KeyBinding{
				{Key: utils.SendRequestShortcut(), Desc: "Send request"},
				{Key: "Enter", Desc: "Enter edit mode"},
				{Key: "Esc", Desc: "Exit edit mode"},
				{Key: "s", Desc: "Save request"},
				{Key: "l", Desc: "Load request menu"},
				{Key: "?", Desc: "Toggle help"},
				{Key: "q/Ctrl+C", Desc: "Quit"},
			},
		},
		{
			Title: "Headers",
			Keys: []KeyBinding{
				{Key: "a", Desc: "Add new header"},
				{Key: "d", Desc: "Delete header"},
				{Key: "p", Desc: "Open presets"},
				{Key: "Space", Desc: "Toggle header"},
				{Key: "Tab", Desc: "Switch key/value"},
				{Key: "j/k", Desc: "Navigate up/down"},
			},
		},
		{
			Title: "Body",
			Keys: []KeyBinding{
				{Key: "Tab", Desc: "Cycle body type"},
			},
		},
		{
			Title: "Response",
			Keys: []KeyBinding{
				{Key: "j/k", Desc: "Scroll up/down"},
				{Key: "g/G", Desc: "Go to top/bottom"},
				{Key: "Tab", Desc: "Cycle view mode"},
				{Key: "f", Desc: "Toggle fullscreen"},
			},
		},
		{
			Title: "Requests Menu",
			Keys: []KeyBinding{
				{Key: "Enter", Desc: "Open/load"},
				{Key: "n", Desc: "New collection"},
				{Key: "r", Desc: "Rename"},
				{Key: "d", Desc: "Delete"},
				{Key: "m", Desc: "Move request"},
				{Key: "Esc", Desc: "Back/close"},
			},
		},
	}
}
