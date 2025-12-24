package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	t.Parallel()

	t.Run("should send a GET request successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer server.Close()

		req := NewModel().SetMethod(GET).
			SetURL(server.URL)

		resp, err := SendRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Empty(t, resp.Body)
	})

	t.Run("should send a POST request with body successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer server.Close()

		req := NewModel().SetMethod(POST).
			SetURL(server.URL).
			SetBody("test body")

		resp, err := SendRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("should read response body correctly", func(t *testing.T) {
		expectedBody := "Hello, World!"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedBody))
		}))
		defer server.Close()

		req := NewModel().SetMethod(GET).
			SetURL(server.URL)

		res, err := SendRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedBody, res.Body)
	})

	t.Run("should send a GET request with headers successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer server.Close()

		req := NewModel().SetMethod(GET).
			SetURL(server.URL).
			AddHeader("X-Test-Header", "HeaderValue")

		resp, err := SendRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("should handle http.NewRequest error", func(t *testing.T) {
		req := NewModel().SetMethod(GET).
			SetURL("http://%41:8080/")

		res, err := SendRequest(req)

		assert.Nil(t, res)
		assert.ErrorContains(t, err, "invalid URL escape")
	})

	t.Run("should timeout for a slow server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}))
		defer server.Close()

		req := NewModel().SetMethod(GET).
			SetURL(server.URL).
			SetTimeout(100)

		res, err := SendRequest(req)

		assert.Nil(t, res)
		assert.ErrorContains(t, err, "context deadline exceeded")
	})

	t.Run("should handle io.ReadAll error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hj, _ := w.(http.Hijacker)
			conn, buf, _ := hj.Hijack()

			buf.WriteString("HTTP/1.1 200 OK\r\n")
			buf.WriteString("Content-Length: 1000\r\n")
			buf.WriteString("\r\n")
			buf.WriteString("partial")
			buf.Flush()
			conn.Close()
		}))
		defer server.Close()

		req := NewModel().SetMethod(GET).
			SetURL(server.URL)

		res, err := SendRequest(req)

		assert.Nil(t, res)
		assert.Error(t, err)
	})
}
