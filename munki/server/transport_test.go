package munkiserver_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/squirrel/munki/datastore"
	"github.com/micromdm/squirrel/munki/server"
	"golang.org/x/net/context"
)

func newServer(t *testing.T) (*httptest.Server, munkiserver.Service) {
	ctx := context.Background()
	l := log.NewLogfmtLogger(os.Stderr)
	logger := log.NewContext(l).With("source", "testing")
	path := "testdata/testrepo"
	repo := &datastore.SimpleRepo{Path: path}
	svc, err := munkiserver.NewService(repo)
	if err != nil {
		t.Fatal(err)
	}
	handler := munkiserver.ServiceHandler(ctx, svc, logger)
	server := httptest.NewServer(handler)
	return server, svc
}
