package datastore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (r *SimpleRepo) AllPkgsinfo() ([]string, error) {
	pkgsPath := fmt.Sprintf("%v/pkgsinfo", r.Path)
	var pkgs []string
	err := filepath.Walk(pkgsPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			rel, err := filepath.Rel(r.Path, path)
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

func (r *SimpleRepo) SavePkgsinfo(path string, file io.Reader) error {
	path = filepath.Join(r.Path, path)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return err
}
