package datastore

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/micromdm/squirrel/munki/models"
	"github.com/micromdm/squirrel/munki/munki"
)

var (
	// ErrExists file already exists
	ErrExists = errors.New("resource already exists")

	// ErrNotFound = resource not found
	ErrNotFound = errors.New("resource not found")
)

// Datastore is an interface around munki storage
type Datastore interface {
	models.PkgsinfoStore
	models.PkgStore
	munki.ManifestStore
}

// SimpleRepo is a filesystem based backend
type SimpleRepo struct {
	Path           string
	indexManifests map[string]*munki.Manifest
	indexPkgsinfo  map[string]*models.PkgsInfo
}

func deleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return ErrNotFound
	}
	return nil
}

func createFile(path string) error {
	// check if exists
	if _, err := os.Stat(path); err == nil {
		return ErrExists
	}
	// create the relative directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
