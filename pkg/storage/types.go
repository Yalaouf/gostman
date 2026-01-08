package storage

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCollectionNotFound = errors.New("collection not found")
	ErrCollectionNotEmpty = errors.New("collection is not empty")
	ErrRequestNotFound    = errors.New("request not found")
)

type Collection struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Request struct {
	ID           string            `json:"id"`
	CollectionID string            `json:"collection_id,omitempty"`
	Name         string            `json:"name"`
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body,omitempty"`
	BodyType     string            `json:"body_type,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type Store struct {
	Version     int           `json:"version"`
	Collections []*Collection `json:"collections"`
	Requests    []*Request    `json:"requests"`
}

type Storage struct {
	mutex sync.RWMutex
	path  string
	store *Store
}
