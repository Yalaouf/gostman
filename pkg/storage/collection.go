package storage

import (
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

func (s *Storage) CreateCollection(name string) (*Collection, error) {
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
	for _, c := range s.store.Collections {
		if c.ID == id {
			return c, nil
		}
	}

	return nil, ErrCollectionNotFound
}

func (s *Storage) ListCollections() []*Collection {
	return s.store.Collections
}

func (s *Storage) UpdateCollection(id, name string) (*Collection, error) {
	i := s.findCollectionIndex(id)
	if i == -1 {
		return nil, ErrCollectionNotFound
	}

	s.store.Collections[i].Name = name
	s.store.Collections[i].UpdatedAt = time.Now()

	if err := s.save(); err != nil {
		return nil, err
	}

	return s.store.Collections[i], nil
}

func (s *Storage) DeleteCollection(id string, force bool) error {
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
	} else {
		s.store.Requests = filterRequests(s.store.Requests, func(r *Request) bool {
			return r.CollectionID != id
		})
	}

	s.store.Collections = append(s.store.Collections[:i], s.store.Collections[i+1:]...)
	return s.save()
}

func (s *Storage) ListRequestsByCollection(collectionID string) []*Request {
	var result []*Request

	for _, req := range s.store.Requests {
		if req.CollectionID == collectionID {
			result = append(result, req)
		}
	}

	return result
}

func (s *Storage) MoveRequest(requestID, collectionID string) error {
	requestIndex := s.findRequestIndex(requestID)
	if requestIndex == -1 {
		return ErrRequestNotFound
	}

	if collectionID != "" && s.findCollectionIndex(collectionID) == -1 {
		return ErrCollectionNotFound
	}

	s.store.Requests[requestIndex].CollectionID = collectionID
	s.store.Requests[requestIndex].UpdatedAt = time.Now()

	return s.save()
}
