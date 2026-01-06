package headers

type Preset struct {
	Key   string
	Value string
}

var CommonPresets = []Preset{
	{"Authorization", "Bearer "},
	{"Accept", "application/json"},
	{"Cache-Control", "no-cache"},
	{"X-Request-ID", ""},
	{"X-API-Key", ""},
	{"Origin", ""},
	{"User-Agent", "gostman/1.0"},
}
