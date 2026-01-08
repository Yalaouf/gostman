package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

func encodeURLEncoded(body string) (string, error) {
	var data map[string]any

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return "", fmt.Errorf("invalid JSON for url-encoded: %w", err)
	}

	values := url.Values{}
	for key, val := range data {
		values.Set(key, fmt.Sprintf("%v", val))
	}

	return values.Encode(), nil
}

func encodeFormData(body string) (io.Reader, string, error) {
	var data map[string]any

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, "", fmt.Errorf("invalid JSON for form-data: %w", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for key, val := range data {
		if err := writer.WriteField(key, fmt.Sprintf("%v", val)); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &buf, writer.FormDataContentType(), nil
}

func EncodeBody(body string, bodyType BodyType) (io.Reader, string, error) {
	if body == "" || bodyType == BodyTypeNone {
		return nil, "", nil
	}

	switch bodyType {
	case BodyTypeJSON:
		return bytes.NewReader([]byte(body)), "application/json", nil
	case BodyTypeURLEncoded:
		encoded, err := encodeURLEncoded(body)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader([]byte(encoded)), "application/x-www-form-urlencoded", nil
	case BodyTypeFormData:
		return encodeFormData(body)
	default:
		return bytes.NewReader([]byte(body)), "", nil
	}
}
