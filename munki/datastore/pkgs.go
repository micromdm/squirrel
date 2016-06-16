package datastore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// AddPkg adds a file to the Datastore
func (r *SimpleRepo) AddPkg(filename string, body io.Reader) error {
	pkgPath := fmt.Sprintf("%v/pkgs/%v", r.Path, filename)
	// check if exists
	if _, err := os.Stat(pkgPath); err == nil {
		return ErrExists
	}

	// create the directory structure
	dir := filepath.Dir(pkgPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// create file
	f, err := os.Create(pkgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// finaly copy the body to db
	_, err = io.Copy(f, body)
	if err != nil {
		return err
	}
	return nil
}

// DeletePkg deletes a file from the repository
func (r *SimpleRepo) DeletePkg(filename string) error {
	pkgPath := fmt.Sprintf("%v/pkgs/%v", r.Path, filename)
	if _, err := os.Stat(pkgPath); err != nil {
		return ErrNotFound
	}
	return os.Remove(pkgPath)
}
