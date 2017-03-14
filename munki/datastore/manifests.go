package datastore

import (
	"fmt"
	"os"
	"path/filepath"
)

func (r *SimpleRepo) AllManifests() ([]string, error) {
	pkgsPath := fmt.Sprintf("%v/manifests", r.Path)
	var pkgs []string
	err := filepath.Walk(pkgsPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			rel, err := filepath.Rel(pkgsPath, path)
			if err != nil {
				return err
			}
			pkgs = append(pkgs, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}
