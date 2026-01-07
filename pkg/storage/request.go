package storage

import (
	"time"

	"github.com/google/uuid"
)

func (s *Storage) findRequestIndex(id string) int {
	for i, req := range s.store.Requests {
		if req.ID == id {
			return i
		}
	}

	return -1
}

func (s *Storage) SaveRequest(req *Request) error {
	now := time.Now()

	if req.ID == "" {
		req.ID = uuid.NewString()
		req.CreatedAt = now
		req.UpdatedAt = now
		s.store.Requests = append(s.store.Requests, req)
	} else {
		i := s.findRequestIndex(req.ID)
		if i == -1 {
			return ErrRequestNotFound
		}

		req.UpdatedAt = now
		req.CreatedAt = s.store.Requests[i].CreatedAt
		s.store.Requests[i] = req
	}

	return s.save()
}

func (s *Storage) GetRequest(id string) (*Request, error) {
	for _, req := range s.store.Requests {
		if req.ID == id {
			return req, nil
		}
	}

	return nil, ErrRequestNotFound
}

func (s *Storage) ListRequests() []*Request {
	return s.store.Requests
}

func (s *Storage) DeleteRequest(id string) error {
	i := s.findRequestIndex(id)
	if i == -1 {
		return ErrRequestNotFound
	}

	s.store.Requests = append(s.store.Requests[:i], s.store.Requests[i+1:]...)
	return s.save()
}
