package storage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageLoadFailure(t *testing.T) {
	s := setupTestStorage(t)

	t.Run("should return on error on ReadFile failure", func(t *testing.T) {
		err := s.load()

		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("should load the saved requests to the store", func(t *testing.T) {
		req := &Request{
			Name:   "Test",
			Method: "GET",
			URL:    "http://localhost",
		}

		require.NoError(t, s.SaveRequest(req))

		s2, err := New()

		assert.NoError(t, err)
		assert.Len(t, s2.store.Requests, 1)
	})
}

func TestNew(t *testing.T) {
	t.Run("should create config directory", func(t *testing.T) {
		s := setupTestStorage(t)

		assert.NotNil(t, s)
		assert.DirExists(t, filepath.Dir(s.path))
	})

	t.Run("should load existing data", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			Name:   "Test",
			Method: "GET",
			URL:    "http://localhost",
		}

		require.NoError(t, s.SaveRequest(req))

		s2, err := New()
		require.NoError(t, err)

		assert.Len(t, s2.ListRequests(), 1)
		assert.Equal(t, "Test", s2.ListRequests()[0].Name)
	})
}
