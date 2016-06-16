package datastore

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/groob/ape/models"
	"github.com/groob/plist"
)

// AllPkgsinfos returns a pkgsinfo collection
func (r *SimpleRepo) AllPkgsinfos() (*models.PkgsInfoCollection, error) {
	pkgsinfos := &models.PkgsInfoCollection{}
	if err := loadPkgsinfos(r.Path, pkgsinfos); err != nil {
		return nil, err
	}
	r.updatePkgsinfoIndex(pkgsinfos)
	return pkgsinfos, nil
}

// Pkgsinfo returns a single pkgsinfo from repo
func (r *SimpleRepo) Pkgsinfo(name string) (*models.PkgsInfo, error) {
	pkgsinfos := &models.PkgsInfoCollection{}
	if err := loadPkgsinfos(r.Path, pkgsinfos); err != nil {
		return nil, err
	}
	r.updatePkgsinfoIndex(pkgsinfos)
	pkgsinfo, ok := r.indexPkgsinfo[name]
	if !ok {
		return nil, ErrNotFound
	}
	return pkgsinfo, nil
}

// NewPkgsinfo returns a single manifest from repo
func (r *SimpleRepo) NewPkgsinfo(name string) (*models.PkgsInfo, error) {
	pkgsinfo := &models.PkgsInfo{}
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo/%v", r.Path, name)
	err := createFile(pkgsinfoPath)
	return pkgsinfo, err
}

// SavePkgsinfo saves a pkgsinfo file to the datastore
func (r *SimpleRepo) SavePkgsinfo(pkgsinfo *models.PkgsInfo) error {
	if pkgsinfo.Filename == "" {
		return errors.New("filename key must be set")
	}
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo/%v", r.Path, pkgsinfo.Filename)
	file, err := os.OpenFile(pkgsinfoPath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := plist.NewEncoder(file).Encode(pkgsinfo); err != nil {
		return err
	}
	go func() {
		makecatalogs <- true
	}()
	return nil
}

// DeletePkgsinfo deletes a pkgsinfo file from the datastore and triggers makecatalogs if succesful
func (r *SimpleRepo) DeletePkgsinfo(name string) error {
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo/%v", r.Path, name)
	if err := deleteFile(pkgsinfoPath); err != nil {
		return err
	}
	go func() {
		makecatalogs <- true
	}()
	return nil
}

func (r *SimpleRepo) updatePkgsinfoIndex(pkgsinfos *models.PkgsInfoCollection) {
	r.indexPkgsinfo = make(map[string]*models.PkgsInfo, len(*pkgsinfos))
	for _, pkgsinfo := range *pkgsinfos {
		r.indexPkgsinfo[pkgsinfo.Filename] = pkgsinfo
	}
}

func walkPkgsinfo(pkgsinfos *models.PkgsInfoCollection, pkgsinfoPath string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if !info.IsDir() {
			// Decode pkgsinfo
			pkgsinfo := &models.PkgsInfo{}
			err := plist.NewDecoder(file).Decode(pkgsinfo)
			if err != nil {
				log.Printf("simple-repo: failed to decode %v, skipping \n", info.Name())
				return nil
			}
			// set filename to relative path
			relpath, err := filepath.Rel(pkgsinfoPath, path)
			if err != nil {
				log.Printf("simple-repo: failed to get relative path %v, skipping \n", info.Name())
				return err
			}
			// use the relative path as the filename
			pkgsinfo.Filename = relpath
			// add to ManifestCollection
			*pkgsinfos = append(*pkgsinfos, pkgsinfo)
			return nil
		}
		return nil
	}
}

// load the pkgsinfos
func loadPkgsinfos(path string, pkgsinfos *models.PkgsInfoCollection) error {
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo", path)
	return filepath.Walk(pkgsinfoPath, walkPkgsinfo(pkgsinfos, pkgsinfoPath))
}
