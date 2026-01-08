package request

import (
	"context"
	"net/http"
)

type HttpMethod string

const (
	GET     HttpMethod = http.MethodGet
	POST    HttpMethod = http.MethodPost
	PUT     HttpMethod = http.MethodPut
	DELETE  HttpMethod = http.MethodDelete
	PATCH   HttpMethod = http.MethodPatch
	OPTIONS HttpMethod = http.MethodOptions
	HEAD    HttpMethod = http.MethodHead
	TRACE   HttpMethod = http.MethodTrace
	CONNECT HttpMethod = http.MethodConnect
)

const DefaultTimeout int64 = 30000

type Model struct {
	Ctx      context.Context
	Method   HttpMethod
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
