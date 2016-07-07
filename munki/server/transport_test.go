package munkiserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func TestShowManifests(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	testShowManifestHTTP(t, server, "site_default", http.StatusOK)
	testShowManifestHTTP(t, server, "site_none", http.StatusNotFound)
}

func TestReplaceManifest(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	manifests := []*munki.Manifest{
		&munki.Manifest{
			Filename: "replace-manifest",
			Catalogs: []string{"production", "testing"},
		},
	}

	for _, m := range manifests {
		os.Remove("testdata/testrepo/manifests/" + m.Filename)
		testCreateManifestHTTP(t, server, m.Filename, m, http.StatusOK)
		testReplaceManifestHTTP(t, server, m.Filename, m, http.StatusOK)
		os.Remove("testdata/testrepo/manifests/" + m.Filename)
	}
}

func TestDeleteManifest(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	manifests := []*munki.Manifest{
		&munki.Manifest{
			Filename: "del-manifest",
			Catalogs: []string{"production", "testing"},
		},
	}
	for _, m := range manifests {
		os.Remove("testdata/testrepo/manifests/" + m.Filename)
		testCreateManifestHTTP(t, server, m.Filename, m, http.StatusOK)
		testDeleteManifestHTTP(t, server, m.Filename, http.StatusOK)
	}
}

func TestCreateManifest(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	manifests := []*munki.Manifest{
		&munki.Manifest{
			Filename: "foo-manifest",
			Catalogs: []string{"production", "testing"},
		},
	}

	for _, m := range manifests {
		testCreateManifestHTTP(t, server, m.Filename, m, http.StatusOK)
		os.Remove("testdata/testrepo/manifests/" + m.Filename)
	}
}

type createManifestRequest struct {
	Filename string `plist:"filename" json:"filename"`
	*munki.Manifest
}

func testCreateManifestHTTP(t *testing.T, server *httptest.Server, filename string, manifest *munki.Manifest, expectedStatus int) *munki.Manifest {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/manifests"
	var req = &createManifestRequest{
		Filename: filename,
		Manifest: manifest,
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	body := ioutil.NopCloser(bytes.NewBuffer(data))
	resp, err := client.Post(theURL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatus {
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}

	return nil
}

func testReplaceManifestHTTP(t *testing.T, server *httptest.Server, path string, m *munki.Manifest, expectedStatus int) {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/manifests/" + path
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	body := ioutil.NopCloser(bytes.NewBuffer(data))
	req, err := http.NewRequest("PUT", theURL, body)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatus {
		fmt.Println(theURL)
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}
}

func testDeleteManifestHTTP(t *testing.T, server *httptest.Server, path string, expectedStatus int) {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/manifests/" + path
	req, err := http.NewRequest("DELETE", theURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatus {
		fmt.Println(theURL)
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}
}

func testShowManifestHTTP(t *testing.T, server *httptest.Server, path string, expectedStatus int) *munki.Manifest {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/manifests/" + path
	resp, err := client.Get(theURL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatus {
		fmt.Println(theURL)
		io.Copy(os.Stdout, resp.Body)
		t.Fatal("expected", expectedStatus, "got", resp.StatusCode)
	}
	return nil
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
