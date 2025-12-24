package request

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func SendRequest(model *Model) (*Response, error) {
	req, err := http.NewRequest(model.MethodString(), model.URL, nil)
	if err != nil {
		return nil, err
	}

	if model.Method == POST ||
		model.Method == PUT ||
		model.Method == PATCH ||
		model.Method == DELETE {
		req.Body = io.NopCloser(strings.NewReader(model.Body))
	}

	for key, value := range model.Header {
		req.Header.Add(key, value)
	}

	req.ContentLength = int64(len(model.Body))

	client := &http.Client{
		Timeout: time.Duration(model.Timeout) * time.Millisecond,
	}

	startTime := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	timeTaken := time.Since(startTime).Milliseconds()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Header:     make(map[string]string),
		TimeTaken:  timeTaken,
		Body:       string(bodyBytes),
	}

	for key, values := range resp.Header {
		if len(values) > 0 {
			response.Header[key] = values[0]
		}
	}

	return response, nil
}
