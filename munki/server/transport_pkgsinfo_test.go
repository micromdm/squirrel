package munkiserver_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/micromdm/squirrel/munki/munki"
)

func TestListPkgsinfos(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()

	testListPkgsinfosHTTP(t, server, http.StatusOK)
}

func testListPkgsinfosHTTP(t *testing.T, server *httptest.Server, expectedStatus int) *munki.PkgsInfoCollection {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/pkgsinfos"
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

func TestCreatePkgsinfo(t *testing.T) {
	server, _ := newServer(t)
	defer server.Close()
	pkgsinfos := []*munki.PkgsInfo{
		&munki.PkgsInfo{
			Filename: "foo-pkgsinfo",
		},
	}

	for _, p := range pkgsinfos {
		testCreatePkgsinfoHTTP(t, server, p.Filename, p, http.StatusOK)
		os.Remove("testdata/testrepo/pkgsinfo/" + p.Filename)
	}
}

type createPkgsinfoRequest struct {
	Filename string `plist:"filename" json:"filename"`
	*munki.PkgsInfo
}

func testCreatePkgsinfoHTTP(t *testing.T, server *httptest.Server, filename string, pkgsinfo *munki.PkgsInfo, expectedStatus int) *munki.PkgsInfo {
	client := http.DefaultClient
	theURL := server.URL + "/api/v1/pkgsinfos"
	var req = &createPkgsinfoRequest{
		Filename: filename,
		PkgsInfo: pkgsinfo,
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
