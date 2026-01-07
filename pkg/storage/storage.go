package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (s *Storage) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s.store)
}

func (s *Storage) save() error {
	data, err := json.MarshalIndent(s.store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func New() (*Storage, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	s := &Storage{
		path: filepath.Join(configDir, "requests.json"),
		store: &Store{
			Version:     1,
			Collections: []*Collection{},
			Requests:    []*Request{},
		},
	}

	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return s, nil
}
