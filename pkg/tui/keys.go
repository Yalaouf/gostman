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
	KeyQuit    = "q"
	KeyCtrlC   = "ctrl+c"
	KeyEscape  = "esc"
	KeyEnter   = "enter"
	KeyMethod  = "m"
	KeyInput   = "i"
	KeyHeaders = "h"
	KeyBody    = "b"
	KeyResult  = "r"
	KeyDown    = "down"
	KeyUp      = "up"
	KeyJ       = "j"
	KeyK       = "k"
	KeyG       = "g"
	KeyShiftG  = "G"
)
