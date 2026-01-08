package storage

import (
	"maps"
	"slices"
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

func (r *Request) Copy() *Request {
	var headers map[string]string
	if r.Headers != nil {
		headers = make(map[string]string, len(r.Headers))
		maps.Copy(headers, r.Headers)
	}

	return &Request{
		ID:           r.ID,
		CollectionID: r.CollectionID,
		Name:         r.Name,
		Method:       r.Method,
		URL:          r.URL,
		Headers:      headers,
		Body:         r.Body,
		BodyType:     r.BodyType,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func (s *Storage) SaveRequest(req *Request) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()

	if req.ID == "" {
		req.ID = uuid.NewString()
		req.CreatedAt = now
		req.UpdatedAt = now
		s.store.Requests = append(s.store.Requests, req)

		if err := s.save(); err != nil {
			s.store.Requests = s.store.Requests[:len(s.store.Requests)-1]
			req.ID = ""
			req.CreatedAt = time.Time{}
			req.UpdatedAt = time.Time{}

			return err
		}

		return nil
	}

	i := s.findRequestIndex(req.ID)
	if i == -1 {
		return ErrRequestNotFound
	}

	oldRequest := s.store.Requests[i]
	req.UpdatedAt = now
	req.CreatedAt = s.store.Requests[i].CreatedAt
	s.store.Requests[i] = req

	if err := s.save(); err != nil {
		s.store.Requests[i] = oldRequest
		return err
	}

	return nil
}

func (s *Storage) GetRequest(id string) (*Request, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, req := range s.store.Requests {
		if req.ID == id {
			return req.Copy(), nil
		}
	}

	return nil, ErrRequestNotFound
}

func (s *Storage) ListRequests() []*Request {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]*Request, len(s.store.Requests))
	for i, r := range s.store.Requests {
		result[i] = r.Copy()
	}

	return result
}

func (s *Storage) DeleteRequest(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i := s.findRequestIndex(id)
	if i == -1 {
		return ErrRequestNotFound
	}

	deleted := s.store.Requests[i]
	s.store.Requests = append(s.store.Requests[:i], s.store.Requests[i+1:]...)

	if err := s.save(); err != nil {
		s.store.Requests = slices.Insert(s.store.Requests, i, deleted)
		return err
	}
	return nil
}
