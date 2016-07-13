package munkiserver

import (
	"github.com/micromdm/squirrel/munki/munki"
	"golang.org/x/net/context"
)

// PkgsinfoService describes the methods for managing Pkgsinfo files in a munki repository
type PkgsinfoService interface {
	ListPkgsinfos(ctx context.Context) (*munki.PkgsInfoCollection, error)
	CreatePkgsinfo(ctx context.Context, name string, pkgsinfo *munki.PkgsInfo) (*munki.PkgsInfo, error)
}

func (svc service) ListPkgsinfos(ctx context.Context) (*munki.PkgsInfoCollection, error) {
	return svc.repo.AllPkgsinfos()
}

func (svc service) CreatePkgsinfo(ctx context.Context, name string, pkgsinfo *munki.PkgsInfo) (*munki.PkgsInfo, error) {
	_, err := svc.repo.NewPkgsinfo(name)
	if err != nil {
		return nil, err
	}
	if err := svc.repo.SavePkgsinfo(name, pkgsinfo); err != nil {
		return nil, err
	}
	return pkgsinfo, nil
}
