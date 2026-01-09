package request

import (
	"context"
	"net/http"
)

type HTTPMethod string

const (
	GET     HTTPMethod = http.MethodGet
	POST    HTTPMethod = http.MethodPost
	PUT     HTTPMethod = http.MethodPut
	DELETE  HTTPMethod = http.MethodDelete
	PATCH   HTTPMethod = http.MethodPatch
	OPTIONS HTTPMethod = http.MethodOptions
	HEAD    HTTPMethod = http.MethodHead
	TRACE   HTTPMethod = http.MethodTrace
	CONNECT HTTPMethod = http.MethodConnect
)

const DefaultTimeout int64 = 30000

type Model struct {
	Ctx      context.Context
	Method   HTTPMethod
	URL      string
	Body     string
	BodyType BodyType
	Headers  map[string]string
	Timeout  int64
	Client   *http.Client
}

type Response struct {
	TimeTaken  int64
	StatusCode int
	Headers    map[string][]string
	Body       string
}

type BodyType uint

const (
	BodyTypeNone BodyType = iota
	BodyTypeJSON
	BodyTypeFormData
	BodyTypeURLEncoded
)
