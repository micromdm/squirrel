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

// AllManifests returns an array of manifests
func (r *SimpleRepo) AllManifests() (*munki.ManifestCollection, error) {
	manifests := &munki.ManifestCollection{}
	err := loadManifests(r.Path, manifests)
	if err != nil {
		return nil, err
	}
	r.updateManifestIndex(manifests)
	return manifests, nil
}

// Manifest returns a single manifest from repo
func (r *SimpleRepo) Manifest(name string) (*munki.Manifest, error) {
	manifests := &munki.ManifestCollection{}
	err := loadManifests(r.Path, manifests)
	if err != nil {
		return nil, err
	}
	r.updateManifestIndex(manifests)

	manifest, ok := r.indexManifests[name]
	if !ok {
		return nil, ErrNotFound
	}
	return manifest, nil
}

// NewManifest returns a single manifest from repo
func (r *SimpleRepo) NewManifest(name string) (*munki.Manifest, error) {
	manifest := &munki.Manifest{}
	manifestPath := fmt.Sprintf("%v/manifests/%v", r.Path, name)
	err := createFile(manifestPath)
	return manifest, err
}

// SaveManifest saves a manifest to the datastore
func (r *SimpleRepo) SaveManifest(manifest *munki.Manifest) error {
	if manifest.Filename == "" {
		return errors.New("filename key must be set")
	}
	manifestPath := fmt.Sprintf("%v/manifests/%v", r.Path, manifest.Filename)
	file, err := os.OpenFile(manifestPath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := plist.NewEncoder(file).Encode(manifest); err != nil {
		return err
	}
	return nil
}

// DeleteManifest removes a manifest file from the repository
func (r *SimpleRepo) DeleteManifest(name string) error {
	manifestPath := fmt.Sprintf("%v/manifests/%v", r.Path, name)
	return deleteFile(manifestPath)
}

func (r *SimpleRepo) updateManifestIndex(manifests *munki.ManifestCollection) {
	r.indexManifests = make(map[string]*munki.Manifest, len(*manifests))
	for _, manifest := range *manifests {
		r.indexManifests[manifest.Filename] = manifest
	}
}

func walkManifests(manifests *munki.ManifestCollection) filepath.WalkFunc {
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
			// Decode manifest
			manifest := &munki.Manifest{}
			err := plist.NewDecoder(file).Decode(manifest)
			if err != nil {
				log.Printf("git-repo: failed to decode %v, skipping \n", info.Name())
				return nil
			}
			// set filename to relative path + filename
			manifest.Filename = info.Name()
			// add to ManifestCollection
			*manifests = append(*manifests, manifest)
			return nil
		}
		return nil
	}
}

func loadManifests(path string, manifests *munki.ManifestCollection) error {
	manifestPath := fmt.Sprintf("%v/manifests", path)
	return filepath.Walk(manifestPath, walkManifests(manifests))
}
