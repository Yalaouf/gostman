package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindCollectionIndex(t *testing.T) {
	s := setupTestStorage(t)

	c, err := s.CreateCollection("Test")
	require.NoError(t, err)

	t.Run("should return -1 if no collection was found", func(t *testing.T) {
		res := s.findCollectionIndex("random ID")

		assert.Equal(t, -1, res)
	})

	t.Run("should return the collection's index if found", func(t *testing.T) {
		res := s.findCollectionIndex(c.ID)

		assert.Equal(t, 0, res)
	})
}

func TestCreateCollection(t *testing.T) {
	t.Run("should create a new collection", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("My API")

		assert.NoError(t, err)
		assert.NotEmpty(t, c.ID)
		assert.Equal(t, "My API", c.Name)
		assert.NotZero(t, c.CreatedAt)
		assert.NotZero(t, c.UpdatedAt)
	})

	t.Run("should persist collection to storage", func(t *testing.T) {
		s := setupTestStorage(t)

		c, _ := s.CreateCollection("My API")

		s2, err := New()
		require.NoError(t, err)

		assert.Len(t, s2.ListCollections(), 1)
		assert.Equal(t, c.ID, s2.ListCollections()[0].ID)
	})
}

func TestGetCollection(t *testing.T) {
	t.Run("should return collection by ID", func(t *testing.T) {
		s := setupTestStorage(t)

		created, err := s.CreateCollection("Test")
		require.NoError(t, err)

		c, err := s.GetCollection(created.ID)

		assert.NoError(t, err)
		assert.Equal(t, created.Name, c.Name)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		s := setupTestStorage(t)

		_, err := s.GetCollection("random-id")

		assert.ErrorIs(t, err, ErrCollectionNotFound)
	})
}

func TestListCollections(t *testing.T) {
	t.Run("should return empty slice when no collections", func(t *testing.T) {
		s := setupTestStorage(t)

		collections := s.ListCollections()

		assert.Empty(t, collections)
	})

	t.Run("should return all collections", func(t *testing.T) {
		s := setupTestStorage(t)

		s.CreateCollection("Collection 1")
		s.CreateCollection("Collection 2")
		s.CreateCollection("Collection 3")

		collections := s.ListCollections()

		assert.Len(t, collections, 3)
		assert.Equal(t, collections[0].Name, "Collection 1")
		assert.Equal(t, collections[1].Name, "Collection 2")
		assert.Equal(t, collections[2].Name, "Collection 3")
	})
}

func TestUpdateCollection(t *testing.T) {
	t.Run("should update collection name", func(t *testing.T) {
		s := setupTestStorage(t)

		created, err := s.CreateCollection("Old Name")
		require.NoError(t, err)

		updated, err := s.UpdateCollection(created.ID, "New Name")

		assert.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, created.ID, updated.ID)
		assert.True(t, updated.UpdatedAt.After(created.CreatedAt))
	})

	t.Run("should persist updated collection", func(t *testing.T) {
		s := setupTestStorage(t)

		created, err := s.CreateCollection("Old Name")
		require.NoError(t, err)

		s.UpdateCollection(created.ID, "New Name")

		s2, _ := New()

		assert.Equal(t, "New Name", s2.ListCollections()[0].Name)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		s := setupTestStorage(t)

		_, err := s.UpdateCollection("random-id", "Name")

		assert.ErrorIs(t, err, ErrCollectionNotFound)
	})
}

func TestDeleteCollection(t *testing.T) {
	t.Run("should delete empty collection", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)

		err = s.DeleteCollection(c.ID, false)

		assert.NoError(t, err)
		assert.Empty(t, s.ListCollections())
	})

	t.Run("should return error when collection has requests and force is false", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		s.SaveRequest(req)

		err = s.DeleteCollection(c.ID, false)

		assert.ErrorIs(t, err, ErrCollectionNotEmpty)
		assert.Len(t, s.ListCollections(), 1)
	})

	t.Run("should force delete collection with requests", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)
		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		s.SaveRequest(req)

		err = s.DeleteCollection(c.ID, true)

		assert.NoError(t, err)
		assert.Empty(t, s.ListCollections())
		assert.Empty(t, s.ListRequests())
	})

	t.Run("should only delete requests in the collection when force deleting", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)

		reqInCollection := &Request{Name: "In Collection", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		reqOutside := &Request{Name: "Outside", Method: "GET", URL: "http://localhost"}
		s.SaveRequest(reqInCollection)
		s.SaveRequest(reqOutside)

		err = s.DeleteCollection(c.ID, true)

		assert.NoError(t, err)
		assert.Len(t, s.ListRequests(), 1)
		assert.Equal(t, "Outside", s.ListRequests()[0].Name)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		s := setupTestStorage(t)

		err := s.DeleteCollection("random-id", false)

		assert.ErrorIs(t, err, ErrCollectionNotFound)
	})
}

func TestListRequestsByCollection(t *testing.T) {
	t.Run("should return empty slice when no requests in collection", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)

		requests := s.ListRequestsByCollection(c.ID)

		assert.Empty(t, requests)
	})

	t.Run("should return only requests in the specified collection", func(t *testing.T) {
		s := setupTestStorage(t)

		c1, err := s.CreateCollection("API 1")
		require.NoError(t, err)
		c2, err := s.CreateCollection("API 2")
		require.NoError(t, err)

		req1 := &Request{Name: "Request 1", Method: "GET", URL: "http://localhost", CollectionID: c1.ID}
		req2 := &Request{Name: "Request 2", Method: "POST", URL: "http://localhost", CollectionID: c1.ID}
		req3 := &Request{Name: "Request 3", Method: "DELETE", URL: "http://localhost", CollectionID: c2.ID}
		s.SaveRequest(req1)
		s.SaveRequest(req2)
		s.SaveRequest(req3)

		c1Requests := s.ListRequestsByCollection(c1.ID)
		c2Requests := s.ListRequestsByCollection(c2.ID)

		assert.Len(t, c1Requests, 2)
		assert.Len(t, c2Requests, 1)
		assert.Equal(t, "Request 3", c2Requests[0].Name)
	})

	t.Run("should return requests without collection when empty string is passed", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		reqInCollection := &Request{Name: "In Collection", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		reqNoCollection := &Request{Name: "No Collection", Method: "GET", URL: "http://localhost"}
		s.SaveRequest(reqInCollection)
		s.SaveRequest(reqNoCollection)

		requests := s.ListRequestsByCollection("")

		assert.Len(t, requests, 1)
		assert.Equal(t, "No Collection", requests[0].Name)
	})
}

func TestMoveRequest(t *testing.T) {
	t.Run("should move request to collection", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		s.SaveRequest(req)

		err = s.MoveRequest(req.ID, c.ID)

		assert.NoError(t, err)
		assert.Len(t, s.ListRequestsByCollection(c.ID), 1)
		assert.Equal(t, c.ID, s.store.Requests[0].CollectionID)
	})

	t.Run("should move request between collections", func(t *testing.T) {
		s := setupTestStorage(t)

		c1, err := s.CreateCollection("API 1")
		require.NoError(t, err)
		c2, err := s.CreateCollection("API 2")
		require.NoError(t, err)
		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost", CollectionID: c1.ID}
		s.SaveRequest(req)

		err = s.MoveRequest(req.ID, c2.ID)

		assert.NoError(t, err)
		assert.Empty(t, s.ListRequestsByCollection(c1.ID))
		assert.Len(t, s.ListRequestsByCollection(c2.ID), 1)
	})

	t.Run("should remove request from collection when empty string is passed", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		s.SaveRequest(req)

		err = s.MoveRequest(req.ID, "")

		assert.NoError(t, err)
		assert.Empty(t, s.ListRequestsByCollection(c.ID))
		assert.Equal(t, "", s.store.Requests[0].CollectionID)
	})

	t.Run("should update the request's UpdatedAt timestamp", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		s.SaveRequest(req)
		originalUpdatedAt := req.UpdatedAt

		s.MoveRequest(req.ID, c.ID)

		assert.True(t, s.store.Requests[0].UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should return error for non-existent request", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		err = s.MoveRequest("random-id", c.ID)

		assert.ErrorIs(t, err, ErrRequestNotFound)
	})

	t.Run("should return error for non-existent collection", func(t *testing.T) {
		s := setupTestStorage(t)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		s.SaveRequest(req)

		err := s.MoveRequest(req.ID, "random-collection-id")

		assert.ErrorIs(t, err, ErrCollectionNotFound)
	})
}
