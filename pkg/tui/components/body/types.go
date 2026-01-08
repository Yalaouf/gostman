package body

type Type uint

const (
	TypeNone Type = iota
	TypeJSON
	TypeFormData
	TypeURLEncoded
)

func (t Type) String() string {
	switch t {
	case TypeNone:
		return "none"
	case TypeJSON:
		return "json"
	case TypeFormData:
		return "form-data"
	case TypeURLEncoded:
		return "x-www-form-urlencoded"
	default:
		return "none"
	}
}

func (t Type) ContentType() string {
	switch t {
	case TypeJSON:
		return "application/json"
	case TypeFormData:
		return "multipart/form-data"
	case TypeURLEncoded:
		return "application/x-www-form-urlencoded"
	default:
		return ""
	}
}

var AllTypes = []Type{TypeNone, TypeJSON, TypeFormData, TypeURLEncoded}
