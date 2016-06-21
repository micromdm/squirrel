package munkiserver_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/squirrel/munki/datastore"
	"github.com/micromdm/squirrel/munki/munki"
	"github.com/micromdm/squirrel/munki/server"
)

func TestListManifests(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	testListManifestsHTTP(t, server, http.StatusOK)
}

func testListManifestsHTTP(t *testing.T, server *httptest.Server, expectedStatus int) *munki.ManifestCollection {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/manifests"
	resp, err := client.Get(theURL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}
	io.Copy(os.Stdout, resp.Body)
	return nil
}

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
