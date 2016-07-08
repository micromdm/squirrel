package munkiserver

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/squirrel/munki/munki"
	"golang.org/x/net/context"
)

type listPkgsinfosRequest struct {
}

type listPkgsinfosResponse struct {
	pkgsinfos *munki.PkgsInfoCollection
	Err       error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r listPkgsinfosResponse) subset() interface{} {
	return r.pkgsinfos
}

func (r listPkgsinfosResponse) error() error { return r.Err }

func makeListPkgsinfosEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		pkgsinfos, err := svc.ListPkgsinfos(ctx)
		return listPkgsinfosResponse{pkgsinfos: pkgsinfos, Err: err}, nil
	}
}
