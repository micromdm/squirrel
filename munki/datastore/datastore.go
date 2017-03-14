package datastore

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	// ErrExists file already exists
	ErrExists = errors.New("resource already exists")

	// ErrNotFound = resource not found
	ErrNotFound = errors.New("resource not found")
)

type SimpleRepo struct {
	Path string
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
