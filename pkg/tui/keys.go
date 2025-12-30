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
	KeyTab     = "tab"
	KeyMethod  = "m"
	KeyInput   = "i"
	KeyHeaders = "h"
	KeyBody    = "b"
	KeyResult  = "r"
	KeyDown    = "down"
	KeyUp      = "up"
	KeyLeft    = "left"
	KeyRight   = "right"
	KeyH       = "h"
	KeyJ       = "j"
	KeyK       = "k"
	KeyL       = "l"
	KeyG       = "g"
	KeyShiftG  = "G"
)
