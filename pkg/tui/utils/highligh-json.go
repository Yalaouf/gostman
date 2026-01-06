package utils

import (
	"bytes"
	"encoding/json"

	"github.com/Yalaouf/gostman/pkg/tui/style"
	"github.com/alecthomas/chroma/v2/quick"
)

func HighlightJSON(body string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(body), "", "  "); err != nil {
		body = prettyJSON.String()
	}

	var buf bytes.Buffer
	err := quick.Highlight(&buf, body, "json", "terminal256", style.ChromaStyle)
	if err != nil {
		return body
	}

	return buf.String()
}
