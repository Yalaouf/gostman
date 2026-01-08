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

	t.Run("should rollback on save failure", func(t *testing.T) {
		s := setupTestStorage(t)
		makeReadOnly(t, s)

		c, err := s.CreateCollection("My API")

		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Empty(t, s.store.Collections)
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

	t.Run("should return a copy, not the internal pointer", func(t *testing.T) {
		s := setupTestStorage(t)

		created, err := s.CreateCollection("Test")
		require.NoError(t, err)

		c1, _ := s.GetCollection(created.ID)
		c2, _ := s.GetCollection(created.ID)

		c1.Name = "Modified"

		assert.Equal(t, "Test", c2.Name, "modifying returned copy should not affect other copies")
		assert.Equal(t, "Test", s.store.Collections[0].Name, "modifying returned copy should not affect internal storage")
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

	t.Run("should return copies, not internal pointers", func(t *testing.T) {
		s := setupTestStorage(t)

		s.CreateCollection("Test")

		collections := s.ListCollections()
		collections[0].Name = "Modified"

		assert.Equal(t, "Test", s.store.Collections[0].Name, "modifying returned slice should not affect internal storage")
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

	t.Run("should return a copy, not the internal pointer", func(t *testing.T) {
		s := setupTestStorage(t)

		created, err := s.CreateCollection("Test")
		require.NoError(t, err)

		updated, err := s.UpdateCollection(created.ID, "New Name")
		require.NoError(t, err)

		updated.Name = "Modified"

		assert.Equal(t, "New Name", s.store.Collections[0].Name, "modifying returned copy should not affect internal storage")
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

	t.Run("should rollback on save failure", func(t *testing.T) {
		s := setupTestStorage(t)

		created, err := s.CreateCollection("Original Name")
		require.NoError(t, err)

		originalUpdatedAt := created.UpdatedAt

		makeReadOnly(t, s)

		_, err = s.UpdateCollection(created.ID, "New Name")

		assert.Error(t, err)
		assert.Equal(t, "Original Name", s.store.Collections[0].Name)
		assert.Equal(t, originalUpdatedAt, s.store.Collections[0].UpdatedAt)
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
		require.NoError(t, s.SaveRequest(req))

		err = s.DeleteCollection(c.ID, false)

		assert.ErrorIs(t, err, ErrCollectionNotEmpty)
		assert.Len(t, s.ListCollections(), 1)
	})

	t.Run("should force delete collection with requests", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)
		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		require.NoError(t, s.SaveRequest(req))

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
		require.NoError(t, s.SaveRequest(reqInCollection))
		require.NoError(t, s.SaveRequest(reqOutside))

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

	t.Run("should persist deletion to storage", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)
		collectionID := c.ID

		err = s.DeleteCollection(collectionID, false)
		require.NoError(t, err)

		s2, err := New()
		require.NoError(t, err)

		assert.Empty(t, s2.ListCollections())
		_, err = s2.GetCollection(collectionID)
		assert.ErrorIs(t, err, ErrCollectionNotFound)
	})

	t.Run("should persist force deletion with requests to storage", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)

		reqInCollection := &Request{Name: "In Collection", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		reqOutside := &Request{Name: "Outside", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(reqInCollection))
		require.NoError(t, s.SaveRequest(reqOutside))

		err = s.DeleteCollection(c.ID, true)
		require.NoError(t, err)

		s2, err := New()
		require.NoError(t, err)

		assert.Empty(t, s2.ListCollections())
		assert.Len(t, s2.ListRequests(), 1)
		assert.Equal(t, "Outside", s2.ListRequests()[0].Name)
	})

	t.Run("should rollback on save failure", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)

		makeReadOnly(t, s)

		err = s.DeleteCollection(c.ID, false)

		assert.Error(t, err)
		assert.Len(t, s.store.Collections, 1)
		assert.Equal(t, "Test", s.store.Collections[0].Name)
	})

	t.Run("should rollback force deletion on save failure", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("Test")
		require.NoError(t, err)

		reqInCollection := &Request{Name: "In Collection", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		reqOutside := &Request{Name: "Outside", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(reqInCollection))
		require.NoError(t, s.SaveRequest(reqOutside))

		makeReadOnly(t, s)

		err = s.DeleteCollection(c.ID, true)

		assert.Error(t, err)
		assert.Len(t, s.store.Collections, 1)
		assert.Len(t, s.store.Requests, 2)
		assert.Equal(t, "In Collection", s.store.Requests[0].Name)
		assert.Equal(t, "Outside", s.store.Requests[1].Name)
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
		require.NoError(t, s.SaveRequest(req1))
		require.NoError(t, s.SaveRequest(req2))
		require.NoError(t, s.SaveRequest(req3))

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
		require.NoError(t, s.SaveRequest(reqInCollection))
		require.NoError(t, s.SaveRequest(reqNoCollection))

		requests := s.ListRequestsByCollection("")

		assert.Len(t, requests, 1)
		assert.Equal(t, "No Collection", requests[0].Name)
	})

	t.Run("should return copies, not internal pointers", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		require.NoError(t, s.SaveRequest(req))

		requests := s.ListRequestsByCollection(c.ID)
		requests[0].Name = "Modified"

		assert.Equal(t, "Test", s.store.Requests[0].Name, "modifying returned slice should not affect internal storage")
	})
}

func TestMoveRequest(t *testing.T) {
	t.Run("should move request to collection", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		reqID := s.ListRequests()[0].ID
		err = s.MoveRequest(reqID, c.ID)

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
		require.NoError(t, s.SaveRequest(req))

		reqID := s.ListRequests()[0].ID
		err = s.MoveRequest(reqID, c2.ID)

		assert.NoError(t, err)
		assert.Empty(t, s.ListRequestsByCollection(c1.ID))
		assert.Len(t, s.ListRequestsByCollection(c2.ID), 1)
	})

	t.Run("should remove request from collection when empty string is passed", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost", CollectionID: c.ID}
		require.NoError(t, s.SaveRequest(req))

		reqID := s.ListRequests()[0].ID
		err = s.MoveRequest(reqID, "")

		assert.NoError(t, err)
		assert.Empty(t, s.ListRequestsByCollection(c.ID))
		assert.Equal(t, "", s.store.Requests[0].CollectionID)
	})

	t.Run("should update the request's UpdatedAt timestamp", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		originalUpdatedAt := s.store.Requests[0].UpdatedAt
		reqID := s.ListRequests()[0].ID

		err = s.MoveRequest(reqID, c.ID)

		require.NoError(t, err)
		assert.True(t, s.store.Requests[0].UpdatedAt.After(originalUpdatedAt) || s.store.Requests[0].UpdatedAt.Equal(originalUpdatedAt))
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
		require.NoError(t, s.SaveRequest(req))

		reqID := s.ListRequests()[0].ID
		err := s.MoveRequest(reqID, "random-collection-id")

		assert.ErrorIs(t, err, ErrCollectionNotFound)
	})

	t.Run("should persist move to storage", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		reqID := s.ListRequests()[0].ID
		err = s.MoveRequest(reqID, c.ID)
		require.NoError(t, err)

		s2, err := New()
		require.NoError(t, err)

		assert.Len(t, s2.ListRequestsByCollection(c.ID), 1)
		assert.Equal(t, c.ID, s2.ListRequests()[0].CollectionID)
	})

	t.Run("should rollback on save failure", func(t *testing.T) {
		s := setupTestStorage(t)

		c, err := s.CreateCollection("API")
		require.NoError(t, err)

		req := &Request{Name: "Test", Method: "GET", URL: "http://localhost"}
		require.NoError(t, s.SaveRequest(req))

		originalUpdatedAt := s.store.Requests[0].UpdatedAt
		reqID := s.ListRequests()[0].ID

		makeReadOnly(t, s)

		err = s.MoveRequest(reqID, c.ID)

		assert.Error(t, err)
		assert.Equal(t, "", s.store.Requests[0].CollectionID)
		assert.Equal(t, originalUpdatedAt, s.store.Requests[0].UpdatedAt)
	})
}

func TestCollectionCopy(t *testing.T) {
	t.Run("should create a copy", func(t *testing.T) {
		original := &Collection{
			ID:   "test-id",
			Name: "Test",
		}

		copied := original.Copy()

		copied.Name = "Modified"

		assert.Equal(t, "Test", original.Name)
	})
}
