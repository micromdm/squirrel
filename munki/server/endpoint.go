package munkiserver

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/squirrel/munki/munki"
	"golang.org/x/net/context"
)

type listManifestsRequest struct {
}

type listManifestsResponse struct {
	manifests *munki.ManifestCollection
	Err       error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r listManifestsResponse) subset() interface{} {
	return r.manifests
}

func makeListManifestsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		manifests, err := svc.ListManifests(ctx)
		return listManifestsResponse{manifests: manifests, Err: err}, nil
	}
}
