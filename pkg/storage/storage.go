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

	tmpFile := s.path + ".tmp"

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		os.Remove(tmpFile)
		return err
	}

	err = f.Sync()
	if err != nil {
		os.Remove(tmpFile)
		return err
	}

	err = f.Close()
	if err != nil {
		os.Remove(tmpFile)
		return err
	}

	return os.Rename(tmpFile, s.path)
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
			Collections: []*Collection{},
			Requests:    []*Request{},
		},
	}

	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return s, nil
}
