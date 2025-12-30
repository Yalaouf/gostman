package request

import "context"

func NewModel() *Model {
	return &Model{
		Header:  make(map[string]string),
		Timeout: DefaultTimeout,
	}
}

func (m *Model) SetContext(ctx context.Context) *Model {
	m.Ctx = ctx
	return m
}

func (m *Model) SetMethod(method HttpMethod) *Model {
	m.Method = method
	return m
}

func (m *Model) SetURL(url string) *Model {
	m.URL = url
	return m
}

func (m *Model) SetBody(body string) *Model {
	m.Body = body
	return m
}

func (m *Model) AddHeader(key, value string) *Model {
	if m.Header == nil {
		m.Header = make(map[string]string)
	}
	m.Header[key] = value
	return m
}

func (m *Model) MethodString() string {
	return string(m.Method)
}

func (m *Model) SetTimeout(timeout int64) *Model {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	m.Timeout = timeout
	return m
}
