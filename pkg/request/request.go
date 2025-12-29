package request

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

func SendRequest(model *Model) (*Response, error) {
	ctx := model.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(model.Timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, model.MethodString(), model.URL, nil)
	if err != nil {
		return nil, err
	}

	if model.Method == POST ||
		model.Method == PUT ||
		model.Method == PATCH ||
		model.Method == DELETE {
		req.Body = io.NopCloser(strings.NewReader(model.Body))
		req.ContentLength = int64(len(model.Body))
	}

	for key, value := range model.Header {
		req.Header.Add(key, value)
	}

	client := http.DefaultClient

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
		Header:     resp.Header,
		TimeTaken:  timeTaken,
		Body:       string(bodyBytes),
	}

	return response, nil
}
