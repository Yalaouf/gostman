package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	t.Run("should return a non-nil model", func(t *testing.T) {
		assert.NotNil(t, model)
	})

	t.Run("should have default timeout", func(t *testing.T) {
		assert.Equal(t, DefaultTimeout, model.Timeout)
	})

	t.Run("should have an initialized header map", func(t *testing.T) {
		assert.NotNil(t, model.Header)
	})
}

func TestModelSetters(t *testing.T) {
	model := NewModel()

	t.Run("SetMethod should set the method", func(t *testing.T) {
		model.SetMethod(POST)
		assert.Equal(t, POST, model.Method)
	})

	t.Run("MethodString should return correct string", func(t *testing.T) {
		assert.Equal(t, "POST", model.MethodString())
	})

	t.Run("SetURL should set the URL", func(t *testing.T) {
		url := "https://example.com"
		model.SetURL(url)
		assert.Equal(t, url, model.URL)
	})

	t.Run("SetBody should set the body", func(t *testing.T) {
		body := "test body"
		model.SetBody(body)
		assert.Equal(t, body, model.Body)
	})

	t.Run("AddHeader should add a header", func(t *testing.T) {
		key := "Content-Type"
		value := "application/json"
		model.AddHeader(key, value)
		assert.Equal(t, value, model.Header[key])
	})

	t.Run("SetTimeout should set a valid timeout", func(t *testing.T) {
		timeout := int64(5000)
		model.SetTimeout(timeout)
		assert.Equal(t, timeout, model.Timeout)
	})

	t.Run("SetTimeout should reset to default on invalid timeout", func(t *testing.T) {
		model.SetTimeout(-1000)
		assert.Equal(t, DefaultTimeout, model.Timeout)
	})
}
