package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func SendRequest(model *Model) (*Response, error) {
	ctx := model.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(model.Timeout)*time.Millisecond)
	defer cancel()

	bodyReader, contentType, err := EncodeBody(model.Body, model.BodyType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, model.MethodString(), model.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		if _, exists := model.Headers["Content-Type"]; !exists {
			req.Header.Set("Content-Type", contentType)
		}
	}

	for key, value := range model.Headers {
		req.Header.Add(key, value)
	}

	client := model.Client
	if client == nil {
		client = http.DefaultClient
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
		Headers:    resp.Header,
		TimeTaken:  timeTaken,
		Body:       string(bodyBytes),
	}

	return response, nil
}
