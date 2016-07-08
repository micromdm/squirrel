package munkiserver

import (
	"github.com/micromdm/squirrel/munki/datastore"
	"github.com/micromdm/squirrel/munki/munki"
	"golang.org/x/net/context"
)

// Service describes the actions of a munki server
type Service interface {
	ListManifests(ctx context.Context) (*munki.ManifestCollection, error)
	ShowManifest(ctx context.Context, name string) (*munki.Manifest, error)
	CreateManifest(ctx context.Context, name string, manifest *munki.Manifest) (*munki.Manifest, error)
	ReplaceManifest(ctx context.Context, name string, manifest *munki.Manifest) (*munki.Manifest, error)
	DeleteManifest(ctx context.Context, name string) error
	UpdateManifest(ctx context.Context, name string, payload *munki.ManifestPayload) (*munki.Manifest, error)
}

type service struct {
	repo datastore.Datastore
}

func (svc service) ListManifests(ctx context.Context) (*munki.ManifestCollection, error) {
	return svc.repo.AllManifests()
}

func (svc service) ShowManifest(ctx context.Context, name string) (*munki.Manifest, error) {
	return svc.repo.Manifest(name)
}

func (svc service) CreateManifest(ctx context.Context, name string, manifest *munki.Manifest) (*munki.Manifest, error) {
	_, err := svc.repo.NewManifest(name)
	if err != nil {
		return nil, err
	}
	if err := svc.repo.SaveManifest(name, manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}

func (svc service) DeleteManifest(ctx context.Context, name string) error {
	return svc.repo.DeleteManifest(name)
}

func (svc service) ReplaceManifest(ctx context.Context, name string, manifest *munki.Manifest) (*munki.Manifest, error) {
	if err := svc.repo.DeleteManifest(name); err != nil {
		return nil, err
	}
	return svc.CreateManifest(ctx, name, manifest)
}

func (svc service) UpdateManifest(ctx context.Context, name string, payload *munki.ManifestPayload) (*munki.Manifest, error) {
	manifest, err := svc.repo.Manifest(name)
	if err != nil {
		return nil, err
	}
	manifest.UpdateFromPayload(payload)
	if err := svc.repo.SaveManifest(name, manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}

// NewService creates a new munki api service
func NewService(repo datastore.Datastore) (Service, error) {
	return &service{
		repo: repo,
	}, nil
}
