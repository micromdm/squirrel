package datastore

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/groob/plist"
	"github.com/micromdm/squirrel/munki/munki"
)

// AllPkgsinfos returns a pkgsinfo collection
func (r *SimpleRepo) AllPkgsinfos() (*munki.PkgsInfoCollection, error) {
	pkgsinfos := &munki.PkgsInfoCollection{}
	if err := loadPkgsinfos(r.Path, pkgsinfos); err != nil {
		return nil, err
	}
	r.updatePkgsinfoIndex(pkgsinfos)
	return pkgsinfos, nil
}

// Pkgsinfo returns a single pkgsinfo from repo
func (r *SimpleRepo) Pkgsinfo(name string) (*munki.PkgsInfo, error) {
	pkgsinfos := &munki.PkgsInfoCollection{}
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
func (r *SimpleRepo) NewPkgsinfo(name string) (*munki.PkgsInfo, error) {
	pkgsinfo := &munki.PkgsInfo{}
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo/%v", r.Path, name)
	err := createFile(pkgsinfoPath)
	return pkgsinfo, err
}

// SavePkgsinfo saves a pkgsinfo file to the datastore
func (r *SimpleRepo) SavePkgsinfo(path string, pkgsinfo *munki.PkgsInfo) error {
	if path == "" {
		return errors.New("must specify a pkgsinfo path")
	}
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo/%v", r.Path, path)
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

func (r *SimpleRepo) updatePkgsinfoIndex(pkgsinfos *munki.PkgsInfoCollection) {
	r.indexPkgsinfo = make(map[string]*munki.PkgsInfo, len(*pkgsinfos))
	for _, pkgsinfo := range *pkgsinfos {
		r.indexPkgsinfo[pkgsinfo.Filename] = pkgsinfo
	}
}

func walkPkgsinfo(pkgsinfos *munki.PkgsInfoCollection, pkgsinfoPath string) filepath.WalkFunc {
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
			pkgsinfo := &munki.PkgsInfo{}
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
func loadPkgsinfos(path string, pkgsinfos *munki.PkgsInfoCollection) error {
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo", path)
	return filepath.Walk(pkgsinfoPath, walkPkgsinfo(pkgsinfos, pkgsinfoPath))
}
