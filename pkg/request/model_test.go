package request

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	t.Parallel()

	t.Run("should return a non-nil model", func(t *testing.T) {
		model := NewModel()

		assert.NotNil(t, model)
	})

	t.Run("should have default timeout", func(t *testing.T) {
		model := NewModel()

		assert.Equal(t, DefaultTimeout, model.Timeout)
	})

	t.Run("should have an initialized header map", func(t *testing.T) {
		model := NewModel()

		assert.NotNil(t, model.Headers)
	})
}

func TestModelSetters(t *testing.T) {
	t.Parallel()

	t.Run("SetContext should set the context", func(t *testing.T) {
		model := NewModel()
		ctx := context.Background()
		model.SetContext(ctx)

		assert.NotNil(t, model.Ctx)
		assert.Equal(t, ctx, model.Ctx)
	})

	t.Run("SetMethod should set the method", func(t *testing.T) {
		model := NewModel()
		model.SetMethod(POST)

		assert.Equal(t, POST, model.Method)
	})

	t.Run("MethodString should return correct string", func(t *testing.T) {
		model := NewModel()
		model.SetMethod(POST)

		assert.Equal(t, "POST", model.MethodString())
	})

	t.Run("SetURL should set the URL", func(t *testing.T) {
		model := NewModel()

		url := "https://example.com"
		model.SetURL(url)
		assert.Equal(t, url, model.URL)
	})

	t.Run("SetBody should set the body", func(t *testing.T) {
		model := NewModel()

		body := "test body"
		model.SetBody(body)
		assert.Equal(t, body, model.Body)
	})

	t.Run("AddHeader should add a header", func(t *testing.T) {
		model := NewModel()

		key := "Content-Type"
		value := "application/json"
		model.AddHeader(key, value)
		assert.Equal(t, value, model.Headers[key])
	})

	t.Run("SetTimeout should set a valid timeout", func(t *testing.T) {
		model := NewModel()

		timeout := int64(5000)
		model.SetTimeout(timeout)
		assert.Equal(t, timeout, model.Timeout)
	})

	t.Run("SetTimeout should reset to default on invalid timeout", func(t *testing.T) {
		model := NewModel()

		model.SetTimeout(-1000)
		assert.Equal(t, DefaultTimeout, model.Timeout)
	})

	t.Run("ClearHeaders should remove all headers", func(t *testing.T) {
		model := NewModel()
		model.AddHeader("Content-Type", "application/json")
		model.AddHeader("Authorization", "Bearer token")

		model.ClearHeaders()

		assert.Empty(t, model.Headers)
		assert.NotNil(t, model.Headers)
	})

	t.Run("SetClient should set the HTTP client", func(t *testing.T) {
		model := NewModel()
		client := &http.Client{}

		model.SetClient(client)

		assert.Equal(t, client, model.Client)
	})

	t.Run("SetBodyType should set the body type", func(t *testing.T) {
		model := NewModel()

		model.SetBodyType(BodyTypeJSON)

		assert.Equal(t, BodyTypeJSON, model.BodyType)
	})
}
