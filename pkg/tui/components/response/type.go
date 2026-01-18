package response

type Tab uint

const (
	TabPretty = iota
	TabRaw
	TabHeaders
	TabTree
)

var AllTabs = []Tab{TabPretty, TabRaw, TabHeaders, TabTree}

func (t Tab) String() string {
	switch t {
	case TabPretty:
		return "pretty"
	case TabRaw:
		return "raw"
	case TabHeaders:
		return "headers"
	case TabTree:
		return "tree"
	default:
		return "pretty"
	}
}
