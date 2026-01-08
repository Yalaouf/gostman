package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterRequests(t *testing.T) {
	t.Run("should return empty slice when no requests match", func(t *testing.T) {
		requests := []*Request{
			{ID: "1", Name: "Request 1", CollectionID: "col-1"},
			{ID: "2", Name: "Request 2", CollectionID: "col-1"},
		}

		result := filterRequests(requests, func(r *Request) bool {
			return r.CollectionID == "col-2"
		})

		assert.Empty(t, result)
	})

	t.Run("should return matching requests", func(t *testing.T) {
		requests := []*Request{
			{ID: "1", Name: "Request 1", CollectionID: "col-1"},
			{ID: "2", Name: "Request 2", CollectionID: "col-2"},
			{ID: "3", Name: "Request 3", CollectionID: "col-1"},
		}

		result := filterRequests(requests, func(r *Request) bool {
			return r.CollectionID == "col-1"
		})

		assert.Len(t, result, 2)
		assert.Equal(t, "Request 1", result[0].Name)
		assert.Equal(t, "Request 3", result[1].Name)
	})

	t.Run("should return all requests when all match", func(t *testing.T) {
		requests := []*Request{
			{ID: "1", Name: "Request 1"},
			{ID: "2", Name: "Request 2"},
		}

		result := filterRequests(requests, func(r *Request) bool {
			return true
		})

		assert.Len(t, result, 2)
	})

	t.Run("should handle empty input slice", func(t *testing.T) {
		var requests []*Request

		result := filterRequests(requests, func(r *Request) bool {
			return true
		})

		assert.Empty(t, result)
	})
}

func TestInsertAt(t *testing.T) {
	t.Run("should insert at beginning", func(t *testing.T) {
		slice := []int{2, 3, 4}

		result := insertAt(slice, 0, 1)

		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("should insert at middle", func(t *testing.T) {
		slice := []int{1, 2, 4, 5}

		result := insertAt(slice, 2, 3)

		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("should insert at end", func(t *testing.T) {
		slice := []int{1, 2, 3}

		result := insertAt(slice, 3, 4)

		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("should work with empty slice", func(t *testing.T) {
		var slice []int

		result := insertAt(slice, 0, 1)

		assert.Equal(t, []int{1}, result)
	})

	t.Run("should work with pointers", func(t *testing.T) {
		r1 := &Request{Name: "Request 1"}
		r2 := &Request{Name: "Request 2"}
		r3 := &Request{Name: "Request 3"}
		slice := []*Request{r1, r3}

		result := insertAt(slice, 1, r2)

		assert.Len(t, result, 3)
		assert.Equal(t, "Request 1", result[0].Name)
		assert.Equal(t, "Request 2", result[1].Name)
		assert.Equal(t, "Request 3", result[2].Name)
	})
}

func TestGetConfigDir(t *testing.T) {
	t.Run("should get the config dir from XDG", func(t *testing.T) {
		tmpDir := os.TempDir()
		t.Setenv("XDG_CONFIG_HOME", tmpDir)

		dir, err := getConfigDir()

		assert.NoError(t, err)
		assert.Equal(t, dir, filepath.Join(tmpDir, "gostman"))
	})

	t.Run("should get the config dir using the home dir if XDG is not set", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		homedir, _ := os.UserHomeDir()

		dir, err := getConfigDir()

		assert.NoError(t, err)
		assert.Equal(t, dir, homedir+"/.config/gostman")
	})

	t.Run("should return an error if UserHomeDir fails", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", "")

		dir, err := getConfigDir()

		assert.EqualError(t, err, "$HOME is not defined")
		assert.Empty(t, dir)
	})
}
