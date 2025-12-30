package response

import (
	"strings"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/charmbracelet/bubbles/viewport"
)

func RenderScrollbar(vp viewport.Model) string {
	height := vp.Height
	if height <= 0 {
		return ""
	}

	totalLines := vp.TotalLineCount()

	if totalLines <= height {
		return strings.Repeat(style.TrackChar+"\n", height)
	}

	thumbSize := max(1, height*height/totalLines)
	thumbPos := int(vp.ScrollPercent() * float64(height-thumbSize))

	var b strings.Builder
	for i := range height {
		if i >= thumbPos && i < thumbPos+thumbSize {
			b.WriteString(style.ThumbChar)
		} else {
			b.WriteString(style.TrackChar)
		}

		if i < height-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
