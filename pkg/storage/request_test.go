package storage

import (
	"os"
	"path/filepath"
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

func makeReadOnly(t *testing.T, s *Storage) {
	os.Remove(s.path)

	dir := filepath.Dir(s.path)
	err := os.Chmod(dir, 0555)
	require.NoError(t, err)

	t.Cleanup(func() {
		os.Chmod(dir, 0755)
	})
}

func TestStorageFindIndex(t *testing.T) {
	s := setupTestStorage(t)

	req := &Request{
		Name:   "Test",
		Method: "GET",
		URL:    "http://localhost",
	}

	require.NoError(t, s.SaveRequest(req))

	t.Run("should return -1 if no request was found", func(t *testing.T) {
		res := s.findRequestIndex("random ID")

		assert.Equal(t, -1, res)
	})

	t.Run("should return the request's index if found", func(t *testing.T) {
		res := s.findRequestIndex(s.store.Requests[0].ID)

		assert.Equal(t, 0, res)
	})
}

func TestStorageSaveRequest(t *testing.T) {
	t.Run("should create a new request", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			Name:   "test",
			Method: "GET",
			URL:    "http://localhost",
		}

		err := s.SaveRequest(req)

		require.NoError(t, err)
		require.Len(t, s.ListRequests(), 1)

		saved := s.ListRequests()[0]
		assert.NotEmpty(t, saved.ID)
		assert.Equal(t, "test", saved.Name)
		assert.Equal(t, "GET", saved.Method)
		assert.Equal(t, "http://localhost", saved.URL)
		assert.NotZero(t, saved.CreatedAt)
		assert.NotZero(t, saved.UpdatedAt)
	})

	t.Run("should not mutate the input request", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			Name:   "test",
			Method: "GET",
			URL:    "http://localhost",
		}

		err := s.SaveRequest(req)

		require.NoError(t, err)
		assert.Empty(t, req.ID, "input should not be mutated")
		assert.True(t, req.CreatedAt.IsZero(), "input should not be mutated")
	})

	t.Run("should modify an existing request", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			Name:   "test",
			Method: "GET",
			URL:    "http://localhost",
		}

		require.NoError(t, s.SaveRequest(req))

		saved := s.ListRequests()[0]
		firstCreatedAt := saved.CreatedAt

		updatedReq := &Request{
			ID:     saved.ID,
			Name:   "updated",
			Method: "POST",
			URL:    "http://example.com",
		}
		err := s.SaveRequest(updatedReq)

		require.NoError(t, err)

		updated := s.ListRequests()[0]
		assert.Equal(t, "updated", updated.Name)
		assert.Equal(t, "POST", updated.Method)
		assert.Equal(t, "http://example.com", updated.URL)
		assert.Equal(t, firstCreatedAt, updated.CreatedAt)
		assert.True(
			t,
			updated.UpdatedAt.After(firstCreatedAt) || updated.UpdatedAt.Equal(firstCreatedAt),
		)
	})

	t.Run("should return an error for non existing requests", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			ID:     "random-id",
			Name:   "Test",
			Method: "GET",
			URL:    "http://localhost",
		}
		err := s.SaveRequest(req)

		assert.ErrorIs(t, err, ErrRequestNotFound)
	})

	t.Run("should return ErrEmptyURL when URL is empty", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			Name:   "Test",
			Method: "GET",
			URL:    "",
		}
		err := s.SaveRequest(req)

		assert.ErrorIs(t, err, ErrEmptyURL)
	})

	t.Run("should return ErrEmptyName when Name is empty", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			Name:   "",
			Method: "GET",
			URL:    "http://localhost",
		}
		err := s.SaveRequest(req)

		assert.ErrorIs(t, err, ErrEmptyName)
	})

	t.Run("should rollback on save failure when creating", func(t *testing.T) {
		s := setupTestStorage(t)
		makeReadOnly(t, s)

		req := &Request{
			Name:   "Test",
			Method: "GET",
			URL:    "http://localhost",
		}

		err := s.SaveRequest(req)

		assert.Error(t, err)
		assert.Empty(t, s.store.Requests)
	})

	t.Run("should rollback on save failure when updating", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{
			Name:   "Original",
			Method: "GET",
			URL:    "http://localhost",
		}
		require.NoError(t, s.SaveRequest(req))

		saved := s.ListRequests()[0]
		originalID := saved.ID
		originalCreatedAt := saved.CreatedAt
		originalUpdatedAt := saved.UpdatedAt

		makeReadOnly(t, s)

		updatedReq := &Request{
			ID:     originalID,
			Name:   "Updated",
			Method: "POST",
			URL:    "http://example.com",
		}
		err := s.SaveRequest(updatedReq)

		assert.Error(t, err)
		assert.Equal(t, "Original", s.store.Requests[0].Name)
		assert.Equal(t, "GET", s.store.Requests[0].Method)
		assert.Equal(t, originalID, s.store.Requests[0].ID)
		assert.Equal(t, originalCreatedAt, s.store.Requests[0].CreatedAt)
		assert.Equal(t, originalUpdatedAt, s.store.Requests[0].UpdatedAt)
	})

	t.Run("should persist request to storage", func(t *testing.T) {
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

func TestStorageGetRequest(t *testing.T) {
	t.Run("should return request by ID", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		savedID := s.ListRequests()[0].ID
		r, err := s.GetRequest(savedID)

		require.NoError(t, err)
		assert.Equal(t, "Test", r.Name)
		assert.Equal(t, savedID, r.ID)
	})

	t.Run("should return a copy, not the internal pointer", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		savedID := s.ListRequests()[0].ID
		r1, _ := s.GetRequest(savedID)
		r2, _ := s.GetRequest(savedID)

		r1.Name = "Modified"

		assert.Equal(t, "Test", r2.Name, "modifying returned copy should not affect other copies")
		assert.Equal(
			t,
			"Test",
			s.store.Requests[0].Name,
			"modifying returned copy should not affect internal storage",
		)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		s := setupTestStorage(t)

		_, err := s.GetRequest("randomID")

		assert.ErrorIs(t, err, ErrRequestNotFound)
	})
}

func TestListRequests(t *testing.T) {
	t.Run("should return empty slice when no requests", func(t *testing.T) {
		s := setupTestStorage(t)

		requests := s.ListRequests()

		assert.Empty(t, requests)
	})

	t.Run("should return all requests", func(t *testing.T) {
		s := setupTestStorage(t)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost"}
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req1))
		require.NoError(t, s.SaveRequest(req2))

		requests := s.ListRequests()

		assert.Len(t, requests, 2)
		assert.Equal(t, "Request 1", requests[0].Name)
		assert.Equal(t, "Request 2", requests[1].Name)
	})

	t.Run("should return copies, not internal pointers", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		requests := s.ListRequests()
		requests[0].Name = "Modified"

		assert.Equal(
			t,
			"Test",
			s.store.Requests[0].Name,
			"modifying returned slice should not affect internal storage",
		)
	})
}

func TestStorageDeleteRequest(t *testing.T) {
	t.Run("should delete request", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		savedID := s.ListRequests()[0].ID
		err := s.DeleteRequest(savedID)

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

		firstID := s.ListRequests()[0].ID
		err := s.DeleteRequest(firstID)

		assert.NoError(t, err)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, "Request 2", s.store.Requests[0].Name)
		assert.Equal(t, "Request 3", s.store.Requests[1].Name)
	})

	t.Run("should delete the middle request", func(t *testing.T) {
		s := setupTestStorage(t)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req1))
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req2))
		req3 := &Request{Name: "Request 3", Method: "DELETE", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req3))

		middleID := s.ListRequests()[1].ID
		err := s.DeleteRequest(middleID)

		assert.NoError(t, err)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, "Request 1", s.store.Requests[0].Name)
		assert.Equal(t, "Request 3", s.store.Requests[1].Name)
	})

	t.Run("should delete the last request", func(t *testing.T) {
		s := setupTestStorage(t)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req1))
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req2))
		req3 := &Request{Name: "Request 3", Method: "DELETE", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req3))

		lastID := s.ListRequests()[2].ID
		err := s.DeleteRequest(lastID)

		assert.NoError(t, err)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, "Request 1", s.store.Requests[0].Name)
		assert.Equal(t, "Request 2", s.store.Requests[1].Name)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		s := setupTestStorage(t)

		err := s.DeleteRequest("randomID")

		assert.ErrorIs(t, err, ErrRequestNotFound)
	})

	t.Run("should persist deletion to storage", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))
		reqID := s.ListRequests()[0].ID

		err := s.DeleteRequest(reqID)
		require.NoError(t, err)

		s2, err := New()
		require.NoError(t, err)

		assert.Empty(t, s2.ListRequests())
		_, err = s2.GetRequest(reqID)
		assert.ErrorIs(t, err, ErrRequestNotFound)
	})

	t.Run("should rollback on save failure", func(t *testing.T) {
		s := setupTestStorage(t)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost"}
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req1))
		require.NoError(t, s.SaveRequest(req2))

		firstID := s.ListRequests()[0].ID

		makeReadOnly(t, s)

		err := s.DeleteRequest(firstID)

		assert.Error(t, err)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, "Request 1", s.store.Requests[0].Name)
		assert.Equal(t, "Request 2", s.store.Requests[1].Name)
	})
}

func TestRequestCopy(t *testing.T) {
	t.Run("should create a deep copy", func(t *testing.T) {
		original := &Request{
			ID:      "test-id",
			Name:    "Test",
			Method:  "GET",
			URL:     "http://localhost",
			Headers: map[string]string{"Authorization": "Bearer token"},
		}

		copied := original.Copy()

		copied.Name = "Modified"
		copied.Headers["Authorization"] = "Modified"

		assert.Equal(t, "Test", original.Name)
		assert.Equal(t, "Bearer token", original.Headers["Authorization"])
	})

	t.Run("should handle nil headers", func(t *testing.T) {
		original := &Request{
			ID:      "test-id",
			Name:    "Test",
			Headers: nil,
		}

		copied := original.Copy()

		assert.Nil(t, copied.Headers)
	})
}
