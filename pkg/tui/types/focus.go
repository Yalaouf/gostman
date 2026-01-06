package types

type FocusSection uint

const (
	FocusMethod FocusSection = iota
	FocusURL
	FocusHeaders
	FocusBody
	FocusResult
)
