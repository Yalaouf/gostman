package storage

import (
	"os"
	"path/filepath"
)

func getConfigDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "gostman"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "gostman"), nil
}

func filterRequests(requests []*Request, keep func(*Request) bool) []*Request {
	var result []*Request
	for _, r := range requests {
		if keep(r) {
			result = append(result, r)
		}
	}

	return result
}
