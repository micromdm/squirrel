package munkiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func decodeListManifestsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listManifestsRequest{}, nil
}

func decodeShowManifestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	return showManifestRequest{Path: path}, nil
}

func decodeCreateManifestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createManifestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeDeleteManifestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	return deleteManifestRequest{Path: path}, nil
}

func decodeReplaceManifestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request replaceManifestRequest
	if err := json.NewDecoder(r.Body).Decode(&request.Manifest); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	request.Path = path
	return request, nil
}

func decodeUpdateManifestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request updateManifestRequest
	if err := json.NewDecoder(r.Body).Decode(&request.ManifestPayload); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	request.Path = path
	return request, nil
}
