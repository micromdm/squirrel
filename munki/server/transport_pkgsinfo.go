package munkiserver

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func decodeListPkgsinfosRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listPkgsinfosRequest{}, nil
}

func decodeCreatePkgsinfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createPkgsinfoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
