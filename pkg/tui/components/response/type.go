package response

type Tab uint

const (
	TabPretty = iota
	TabRaw
	TabHeaders
)

var AllTabs = []Tab{TabPretty, TabRaw, TabHeaders}

func (t Tab) String() string {
	switch t {
	case TabPretty:
		return "pretty"
	case TabRaw:
		return "raw"
	case TabHeaders:
		return "headers"
	default:
		return "pretty"
	}
}
