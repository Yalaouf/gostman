package storage

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

func (s *Storage) findCollectionIndex(id string) int {
	for i, c := range s.store.Collections {
		if c.ID == id {
			return i
		}
	}

	return -1
}

func (c *Collection) Copy() *Collection {
	return &Collection{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (s *Storage) CreateCollection(name string) (*Collection, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()

	collection := &Collection{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.store.Collections = append(s.store.Collections, collection)

	if err := s.save(); err != nil {
		s.store.Collections = s.store.Collections[:len(s.store.Collections)-1]
		return nil, err
	}

	return collection, nil
}

func (s *Storage) GetCollection(id string) (*Collection, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, c := range s.store.Collections {
		if c.ID == id {
			return c.Copy(), nil
		}
	}

	return nil, ErrCollectionNotFound
}

func (s *Storage) ListCollections() []*Collection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]*Collection, len(s.store.Collections))
	for i, c := range s.store.Collections {
		result[i] = c.Copy()
	}

	return result
}

func (s *Storage) UpdateCollection(id, name string) (*Collection, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i := s.findCollectionIndex(id)
	if i == -1 {
		return nil, ErrCollectionNotFound
	}

	oldName := s.store.Collections[i].Name
	oldUpdatedAt := s.store.Collections[i].UpdatedAt

	s.store.Collections[i].Name = name
	s.store.Collections[i].UpdatedAt = time.Now()

	if err := s.save(); err != nil {
		s.store.Collections[i].Name = oldName
		s.store.Collections[i].UpdatedAt = oldUpdatedAt
		return nil, err
	}

	return s.store.Collections[i].Copy(), nil
}

func (s *Storage) DeleteCollection(id string, force bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i := s.findCollectionIndex(id)
	if i == -1 {
		return ErrCollectionNotFound
	}

	if !force {
		for _, req := range s.store.Requests {
			if req.CollectionID == id {
				return ErrCollectionNotEmpty
			}
		}
	}

	deletedCollection := s.store.Collections[i]
	var originalRequests []*Request

	if force {
		originalRequests = s.store.Requests
		s.store.Requests = filterRequests(s.store.Requests, func(r *Request) bool {
			return r.CollectionID != id
		})
	}

	s.store.Collections = append(s.store.Collections[:i], s.store.Collections[i+1:]...)

	if err := s.save(); err != nil {
		s.store.Collections = slices.Insert(s.store.Collections, i, deletedCollection)

		if force {
			s.store.Requests = originalRequests
		}

		return err
	}

	return nil
}

func (s *Storage) ListRequestsByCollection(collectionID string) []*Request {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*Request

	for _, req := range s.store.Requests {
		if req.CollectionID == collectionID {
			result = append(result, req.Copy())
		}
	}

	return result
}

func (s *Storage) MoveRequest(requestID, collectionID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	requestIndex := s.findRequestIndex(requestID)
	if requestIndex == -1 {
		return ErrRequestNotFound
	}

	if collectionID != "" && s.findCollectionIndex(collectionID) == -1 {
		return ErrCollectionNotFound
	}

	oldCollectionID := s.store.Requests[requestIndex].CollectionID
	oldUpdatedAt := s.store.Requests[requestIndex].UpdatedAt

	s.store.Requests[requestIndex].CollectionID = collectionID
	s.store.Requests[requestIndex].UpdatedAt = time.Now()

	if err := s.save(); err != nil {
		s.store.Requests[requestIndex].CollectionID = oldCollectionID
		s.store.Requests[requestIndex].UpdatedAt = oldUpdatedAt
		return err
	}

	return nil
}
