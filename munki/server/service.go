package munkiserver

import "github.com/micromdm/squirrel/munki/datastore"

// Service describes the actions of a munki server
type Service interface {
	ManifestService
	PkgsinfoService
}

type service struct {
	repo datastore.Datastore
}

// NewService creates a new munki api service
func NewService(repo datastore.Datastore) (Service, error) {
	return &service{
		repo: repo,
	}, nil
}
