package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (r *SimpleRepo) AllCatalogs() ([]string, error) {
	pkgsPath := fmt.Sprintf("%v/catalogs", r.Path)
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

func (r *SimpleRepo) GetFile(kind, name string) ([]byte, error) {
	pkgsPath := fmt.Sprintf("%s/%s", r.Path, kind)
	var file []byte
	err := filepath.Walk(pkgsPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if info.Name() == name {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				file = data
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return file, nil
}
