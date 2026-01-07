package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStorage(t *testing.T) *Storage {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	s, err := New()
	require.NoError(t, err)

	return s
}

func TestStorageFindIndex(t *testing.T) {
	s := setupTestStorage(t)

	req := &Request{
		Name:   "Test",
		Method: "GET",
		URL:    "http://localhost",
	}

	s.SaveRequest(req)

	t.Run("should return -1 if no request was found", func(t *testing.T) {
		res := s.findRequestIndex("random ID")

		assert.Equal(t, res, -1)
	})

	t.Run("should return the request's index if found", func(t *testing.T) {
		res := s.findRequestIndex(s.store.Requests[0].ID)

		assert.Equal(t, res, 0)
	})
}

func TestStorageSaveRequest(t *testing.T) {
	t.Run("should create a new request", func(t *testing.T) {
		req := &Request{
			Name:   "test",
			Method: "GET",
			URL:    "http://localhost",
		}

		s := setupTestStorage(t)

		err := s.SaveRequest(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, req.ID)
		assert.NotZero(t, req.CreatedAt)
		assert.NotZero(t, req.UpdatedAt)
	})

	t.Run("should modify an existing request", func(t *testing.T) {
		req := &Request{
			Name:   "test",
			Method: "GET",
			URL:    "http://localhost",
		}

		s := setupTestStorage(t)

		require.NoError(t, s.SaveRequest(req))

		firstCreatedAt := req.CreatedAt
		req.Name = "updated"
		err := s.SaveRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, "updated", req.Name)
		assert.Equal(t, firstCreatedAt, req.CreatedAt)
	})

	t.Run("should return an error for non existing requests", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{ID: "random ID", Name: "Test"}
		err := s.SaveRequest(req)

		assert.ErrorIs(t, err, ErrRequestNotFound)
	})
}

func TestStorageGetRequest(t *testing.T) {
	t.Run("should return request by ID", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		r, err := s.GetRequest(req.ID)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, r.Name)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		s := setupTestStorage(t)

		_, err := s.GetRequest("randomID")

		assert.ErrorIs(t, err, ErrRequestNotFound)
	})
}

func TestStorageDeleteRequest(t *testing.T) {
	t.Run("should delete request", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		err := s.DeleteRequest(req.ID)

		assert.NoError(t, err)
		assert.Empty(t, s.ListRequests())
	})

	t.Run("should delete the first request", func(t *testing.T) {
		s := setupTestStorage(t)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req1))
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req2))
		req3 := &Request{Name: "Request 3", Method: "DELETE", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req3))

		err := s.DeleteRequest(req1.ID)

		assert.NoError(t, err)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, s.store.Requests[0].Name, req2.Name)
		assert.Equal(t, s.store.Requests[1].Name, req3.Name)
	})

	t.Run("should delete the middle request", func(t *testing.T) {
		s := setupTestStorage(t)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req1))
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req2))
		req3 := &Request{Name: "Request 3", Method: "DELETE", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req3))

		err := s.DeleteRequest(req2.ID)

		assert.NoError(t, err)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, s.store.Requests[0].Name, req1.Name)
		assert.Equal(t, s.store.Requests[1].Name, req3.Name)
	})

	t.Run("should delete the last request", func(t *testing.T) {
		s := setupTestStorage(t)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req1))
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req2))
		req3 := &Request{Name: "Request 3", Method: "DELETE", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req3))

		err := s.DeleteRequest(req3.ID)

		assert.NoError(t, err)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, s.store.Requests[0].Name, req1.Name)
		assert.Equal(t, s.store.Requests[1].Name, req2.Name)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		s := setupTestStorage(t)

		err := s.DeleteRequest("randomID")

		assert.ErrorIs(t, err, ErrRequestNotFound)
	})
}
