package munkiserver

import (
	"github.com/micromdm/squirrel/munki/munki"
	"golang.org/x/net/context"
)

// PkgsinfoService describes the methods for managing Pkgsinfo files in a munki repository
type PkgsinfoService interface {
	ListPkgsinfos(ctx context.Context) (*munki.PkgsInfoCollection, error)
}

func (svc service) ListPkgsinfos(ctx context.Context) (*munki.PkgsInfoCollection, error) {
	return svc.repo.AllPkgsinfos()
}
