package request

import "context"

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	DELETE  HttpMethod = "DELETE"
	PATCH   HttpMethod = "PATCH"
	OPTIONS HttpMethod = "OPTIONS"
	HEAD    HttpMethod = "HEAD"
	TRACE   HttpMethod = "TRACE"
	CONNECT HttpMethod = "CONNECT"
)

const DefaultTimeout int64 = 30000

type Model struct {
	Ctx     context.Context
	Method  HttpMethod
	URL     string
	Body    string
	Headers map[string]string
	Timeout int64
}

type Response struct {
	TimeTaken  int64
	StatusCode int
	Headers    map[string][]string
	Body       string
}
