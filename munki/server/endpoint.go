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

func (r listManifestsResponse) error() error { return r.Err }

func makeListManifestsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		manifests, err := svc.ListManifests(ctx)
		return listManifestsResponse{manifests: manifests, Err: err}, nil
	}
}

type showManifestRequest struct {
	Path string
}

type showManifestResponse struct {
	*munki.Manifest
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r showManifestResponse) error() error { return r.Err }

func makeShowManifestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(showManifestRequest)
		manifest, err := svc.ShowManifest(ctx, req.Path)
		return showManifestResponse{Manifest: manifest, Err: err}, nil
	}
}

type createManifestRequest struct {
	Filename string `plist:"filename" json:"filename"`
	*munki.Manifest
}

type createManifestResponse struct {
	*munki.Manifest
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r createManifestResponse) error() error { return r.Err }

func makeCreateManifestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createManifestRequest)
		manifest, err := svc.CreateManifest(ctx, req.Filename, req.Manifest)
		return showManifestResponse{Manifest: manifest, Err: err}, nil
	}
}

type deleteManifestRequest struct {
	Path string `plist:"filename" json:"filename"`
}

type deleteManifestResponse struct {
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r deleteManifestResponse) error() error { return r.Err }

func makeDeleteManifestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteManifestRequest)
		err := svc.DeleteManifest(ctx, req.Path)
		return deleteManifestResponse{Err: err}, nil
	}
}

type replaceManifestRequest struct {
	Path string `plist:"filename" json:"filename"`
	*munki.Manifest
}

type replaceManifestResponse struct {
	*munki.Manifest
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r replaceManifestResponse) error() error { return r.Err }

func makeReplaceManifestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(replaceManifestRequest)
		manifest, err := svc.ReplaceManifest(ctx, req.Path, req.Manifest)
		return replaceManifestResponse{Manifest: manifest, Err: err}, nil
	}
}

type updateManifestRequest struct {
	Path string `plist:"filename" json:"filename"`
	*munki.ManifestPayload
}

type updateManifestResponse struct {
	*munki.Manifest
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r updateManifestResponse) error() error { return r.Err }

func makeUpdateManifestEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateManifestRequest)
		manifest, err := svc.UpdateManifest(ctx, req.Path, req.ManifestPayload)
		return updateManifestResponse{Manifest: manifest, Err: err}, nil
	}
}
