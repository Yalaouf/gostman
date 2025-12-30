package body

type Type uint

const (
	TypeNone Type = iota
	TypeRaw
	TypeFormData
	TypeURLEncoded
)

func (t Type) String() string {
	switch t {
	case TypeNone:
		return "none"
	case TypeRaw:
		return "raw"
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
	case TypeRaw:
		return "application/json"
	case TypeFormData:
		return "multipart/form-data"
	case TypeURLEncoded:
		return "application/x-www-form-urlencoded"
	default:
		return ""
	}
}

var AllTypes = []Type{TypeNone, TypeRaw, TypeFormData, TypeURLEncoded}
