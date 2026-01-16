package utils

import "runtime"

func SendRequestShortcut() string {
	if runtime.GOOS == "darwin" {
		return "ctrl+g"
	}
	return "alt+enter"
}
