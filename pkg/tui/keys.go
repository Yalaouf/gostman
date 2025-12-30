package tui

type FocusSection uint

const (
	FocusMethod FocusSection = iota
	FocusURL
	FocusHeaders
	FocusBody
	FocusResult
)

const (
	KeyBody     = "b"
	KeyCtrlC    = "ctrl+c"
	KeyAltEnter = "alt+enter"
	KeyDown     = "down"
	KeyEnter    = "enter"
	KeyEscape   = "esc"
	KeyG        = "g"
	KeyHeaders  = "h"
	KeyH        = "h"
	KeyInput    = "i"
	KeyJ        = "j"
	KeyK        = "k"
	KeyLeft     = "left"
	KeyL        = "l"
	KeyMethod   = "m"
	KeyQuit     = "q"
	KeyResult   = "r"
	KeyRight    = "right"
	KeyShiftG   = "G"
	KeyTab      = "tab"
	KeyUp       = "up"
)
