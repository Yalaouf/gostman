package request

type httpMethod string

const (
	GET     httpMethod = "GET"
	POST    httpMethod = "POST"
	PUT     httpMethod = "PUT"
	DELETE  httpMethod = "DELETE"
	PATCH   httpMethod = "PATCH"
	OPTIONS httpMethod = "OPTIONS"
	HEAD    httpMethod = "HEAD"
	TRACE   httpMethod = "TRACE"
	CONNECT httpMethod = "CONNECT"
)

const DefaultTimeout int64 = 30000

type Model struct {
	Method  httpMethod
	URL     string
	Body    string
	Header  map[string]string
	Timeout int64 // in milliseconds
}

type Response struct {
	TimeTaken  int64 // in milliseconds
	StatusCode int
	Header     map[string]string
	Body       string
}
