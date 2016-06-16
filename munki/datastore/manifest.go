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

// AllManifests returns an array of manifests
func (r *SimpleRepo) AllManifests() (*models.ManifestCollection, error) {
	manifests := &models.ManifestCollection{}
	err := loadManifests(r.Path, manifests)
	if err != nil {
		return nil, err
	}
	r.updateManifestIndex(manifests)
	return manifests, nil
}

// Manifest returns a single manifest from repo
func (r *SimpleRepo) Manifest(name string) (*models.Manifest, error) {
	manifests := &models.ManifestCollection{}
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
func (r *SimpleRepo) NewManifest(name string) (*models.Manifest, error) {
	manifest := &models.Manifest{}
	manifestPath := fmt.Sprintf("%v/manifests/%v", r.Path, name)
	err := createFile(manifestPath)
	return manifest, err
}

// SaveManifest saves a manifest to the datastore
func (r *SimpleRepo) SaveManifest(manifest *models.Manifest) error {
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

func (r *SimpleRepo) updateManifestIndex(manifests *models.ManifestCollection) {
	r.indexManifests = make(map[string]*models.Manifest, len(*manifests))
	for _, manifest := range *manifests {
		r.indexManifests[manifest.Filename] = manifest
	}
}

func walkManifests(manifests *models.ManifestCollection) filepath.WalkFunc {
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
			manifest := &models.Manifest{}
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

func loadManifests(path string, manifests *models.ManifestCollection) error {
	manifestPath := fmt.Sprintf("%v/manifests", path)
	return filepath.Walk(manifestPath, walkManifests(manifests))
}
