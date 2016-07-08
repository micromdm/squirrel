package munkiserver

import (
	"net/http"

	"golang.org/x/net/context"
)

func decodeListPkgsinfosRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listPkgsinfosRequest{}, nil
}
