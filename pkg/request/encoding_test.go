package request

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeURLEncoded(t *testing.T) {
	t.Run("should return an error if JSON is not valid", func(t *testing.T) {
		body := "not a valid JSON"

		res, err := encodeURLEncoded(body)

		assert.Empty(t, res)
		assert.ErrorContains(t, err, "invalid JSON for url-encoded")
	})

	t.Run("should return an encoded JSON if valid", func(t *testing.T) {
		body := "{\"valid\":\"json\"}"
		wanted := "valid=json"

		res, err := encodeURLEncoded(body)

		assert.NoError(t, err)
		assert.Equal(t, wanted, res)
	})
}

func TestEncodeFormData(t *testing.T) {
	t.Run("should return an error on invalid JSON", func(t *testing.T) {
		body := "not a valid JSON"

		r, ct, err := encodeFormData(body)

		assert.Empty(t, r)
		assert.Empty(t, ct)
		assert.ErrorContains(t, err, "invalid JSON for form-data")
	})

	t.Run("should encode and return on success", func(t *testing.T) {
		body := "{\"key\":\"value\"}"

		r, ct, err := encodeFormData(body)

		assert.NotEmpty(t, r)
		assert.Contains(t, ct, "multipart/form-data; boundary=")
		assert.NoError(t, err)
	})
}

func TestEncodeBody(t *testing.T) {
	t.Run("should return a nil reader on BodyTypeNone", func(t *testing.T) {
		body := ""
		contentType := BodyTypeNone

		r, ct, err := EncodeBody(body, contentType)

		assert.NoError(t, err)
		assert.Nil(t, r)
		assert.Empty(t, ct)
	})

	t.Run("should return the unchanged body on BodyTypeJSON", func(t *testing.T) {
		body := "{\"key\":\"value\"}"
		contentType := BodyTypeJSON

		r, ct, err := EncodeBody(body, contentType)

		assert.NoError(t, err)
		assert.Equal(t, r, bytes.NewReader([]byte(body)))
		assert.Equal(t, ct, "application/json")
	})

	t.Run("should return the encoded body on BodyTypeURLEncoded", func(t *testing.T) {
		body := "{\"key\":\"value\"}"
		contentType := BodyTypeURLEncoded

		encoded, err := encodeURLEncoded(body)
		require.NoError(t, err)

		r, ct, err := EncodeBody(body, contentType)

		assert.NoError(t, err)
		assert.Equal(t, r, bytes.NewReader([]byte(encoded)))
		assert.Equal(t, ct, "application/x-www-form-urlencoded")
	})

	t.Run("should return error on BodyTypeURLEncoded with invalid JSON", func(t *testing.T) {
		body := "not a valid JSON"
		contentType := BodyTypeURLEncoded

		r, ct, err := EncodeBody(body, contentType)

		assert.Error(t, err)
		assert.Nil(t, r)
		assert.Empty(t, ct)
		assert.ErrorContains(t, err, "invalid JSON for url-encoded")
	})

	t.Run("should return a multipart form on BodyTypeFormData", func(t *testing.T) {
		body := "{\"key\":\"value\"}"
		contentType := BodyTypeFormData

		r, ct, err := EncodeBody(body, contentType)

		assert.NoError(t, err)
		assert.NotEmpty(t, r)
		assert.Contains(t, ct, "multipart/form-data;")
	})

	t.Run("should return error on BodyTypeFormData with invalid JSON", func(t *testing.T) {
		body := "not a valid JSON"
		contentType := BodyTypeFormData

		r, ct, err := EncodeBody(body, contentType)

		assert.Error(t, err)
		assert.Nil(t, r)
		assert.Empty(t, ct)
		assert.ErrorContains(t, err, "invalid JSON for form-data")
	})
}
