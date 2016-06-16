package api

import (
	"ape/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/groob/plist"
)

func TestManifestsShow(t *testing.T) {
	f := func(accept string) {
		req, err := client.NewRequest("manifests", "site_default", accept, "GET")
		if err != nil {
			t.Fatal(err)
		}
		setHeader(req, accept)
		resp, err := client.Do(req, nil)
		if err != nil {
			t.Fatal(err)
		}

		// check that xml and json are returned when an accept header is passed
		contentType := resp.Header.Get("Content-Type")
		if contentType != accept {
			t.Error("Expected", accept, "got", contentType)
		}

		// check that the API returns 200 OK
		if resp.StatusCode != 200 {
			t.Fatal("Expected", 200, "got", resp.StatusCode)
		}

		// check decoded response
		manifest := &models.Manifest{}
		switch accept {
		case xmlMedia:
			err = plist.NewDecoder(resp.Body).Decode(manifest)
		case jsonMedia:
			err = json.NewDecoder(resp.Body).Decode(manifest)
		}
		if err != nil {
			t.Fatal(err)
		}
		if manifest.Catalogs[0] != "production" {
			t.Error("Expected", "production", "got", manifest.Catalogs[0])
		}

	}

	f(xmlMedia)
	f(jsonMedia)
}

// Test POST request to /api/manifests
func TestCreateManifest(t *testing.T) {

	f := func(mediaType string, body []byte, status int) {
		req, err := client.NewRequest("manifests", "", mediaType, "POST")
		if err != nil {
			t.Fatal(err)
		}
		setHeader(req, mediaType)
		req.Body = &nopCloser{bytes.NewBuffer(body)}
		resp, err := client.Do(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != status {
			t.Error("Expected", status, "got", resp.StatusCode)
			io.Copy(os.Stdout, resp.Body)
		}
	}

	f(xmlMedia, manifestPLISTNew, http.StatusCreated)
	f(jsonMedia, manifestJSONNew, http.StatusCreated)
	// should return 409 because the resource already exists
	f(jsonMedia, manifestJSONNew, http.StatusConflict)
	// empty body, should return bad request
	f(jsonMedia, nil, http.StatusBadRequest)
	// no filename, should return Bad request
	f(jsonMedia, manifestJSONUpdate, http.StatusBadRequest)
	// malformed json
	f(jsonMedia, manifestJSONMalformed, http.StatusInternalServerError)

}

// Test PATCH, PUT and DELETE requests to /api/manifests/:manifest
func TestUpdateDeleteManifest(t *testing.T) {

	f := func(name, method, mediaType string, body []byte, status int) {
		req, err := client.NewRequest("manifests", name, mediaType, method)
		if err != nil {
			t.Fatal(err)
		}
		setHeader(req, mediaType)
		req.Body = &nopCloser{bytes.NewBuffer(body)}
		resp, err := client.Do(req, nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != status {
			t.Fatal("Expected", status, "got", resp.StatusCode, "Method", method)
		}
	}

	// PATCH
	f("does-not-exist", "PATCH", xmlMedia, manifestPLISTUpdate, http.StatusNotFound)
	f("CSXFLNGL9UW0", "PATCH", xmlMedia, manifestPLISTUpdate, http.StatusOK)
	f("does-not-exist", "PATCH", jsonMedia, manifestJSONUpdate, http.StatusNotFound)
	f("CSXFLNGL9UW0", "PATCH", jsonMedia, manifestJSONUpdate, http.StatusOK)
	f("CSXFLNGL9UW0", "PATCH", jsonMedia, manifestJSONMalformed, http.StatusInternalServerError)
	f("CSXFLNGL9UW0", "PATCH", xmlMedia, manifestJSONUpdate, http.StatusBadRequest)
	// PUT
	f("does-not-exist", "PUT", xmlMedia, manifestPLISTUpdate, http.StatusNotFound)
	f("CSXFLNGL9UW0", "PUT", xmlMedia, manifestPLISTUpdate, http.StatusOK)
	f("does-not-exist", "PUT", jsonMedia, manifestJSONUpdate, http.StatusNotFound)
	f("CSXFLNGL9UW0", "PUT", jsonMedia, manifestJSONUpdate, http.StatusOK)
	f("CSXFLNGL9UW0", "PUT", xmlMedia, manifestJSONUpdate, http.StatusBadRequest)
	// DELETE
	f("does-not-exist", "DELETE", jsonMedia, nil, http.StatusNotFound)
	f("CSXFLNGL9UW0", "DELETE", jsonMedia, nil, http.StatusNoContent)
	f("CSXFLNGL9UW0", "DELETE", xmlMedia, nil, http.StatusNotFound)
	// GET a manifest
	f("site_default", "GET", xmlMedia, nil, http.StatusOK)
	f("does-not-exist", "GET", jsonMedia, nil, http.StatusNotFound)
	// GET all manifests
	f("", "GET", xmlMedia, nil, http.StatusOK)
	// method not allowed error
	f("CSXFLNGL9UW0", "POST", xmlMedia, manifestJSONUpdate, http.StatusMethodNotAllowed)

}

// check that a json array of errors is returned
func TestErrorArray(t *testing.T) {
	accept := "application/json; charset=utf-8"
	req, err := client.NewRequest("manifests", "does-not-exist", accept, "GET")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	errJsn := &models.ErrorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(errJsn); err != nil {
		t.Fatal(err)
	}
	if len(errJsn.Errors) < 1 {
		t.Error("Expected", ">= 1 errors", "got", len(errJsn.Errors))
	}
}

var manifestPLISTNew = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>filename</key>
	<string>CSXFLNGL9UW0</string>
	<key>catalogs</key>
	<array>
		<string>testing</string>
	</array>
	<key>managed_installs</key>
	<array>
		<string>GoogleChrome</string>
		<string>munkitools</string>
		<string>munkitools_admin</string>
		<string>munkitools_core</string>
		<string>munkitools_launchd</string>
		<string>sal_scripts</string>
	</array>
	<key>optional_installs</key>
	<array>
		<string>AdobePhotoshop</string>
	</array>
</dict>
</plist>`)

var manifestPLISTUpdate = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>catalogs</key>
	<array>
		<string>testing</string>
		<string>production</string>
	</array>
</dict>
</plist>`)

var manifestJSONNew = []byte(`{
	"filename": "CSXFLNGL9UW1",
	 "catalogs": [
	  "testing"
	 ],
	 "optional_installs": [
	  "AdobePhotoshop"
	 ],
	 "managed_installs": [
	  "GoogleChrome",
	  "munkitools",
	  "munkitools_admin",
	  "munkitools_core",
	  "munkitools_launchd",
	  "sal_scripts"
	 ]
	}`)

var manifestJSONUpdate = []byte(`{
	 "catalogs": [
	  "testing"
	 ],
	 "optional_installs": [
	  "AdobePhotoshop"
	 ],
	 "managed_installs": [
	  "munkitools",
	  "munkitools_core",
	  "munkitools_launchd",
	  "sal_scripts"
	 ]
	}`)

var manifestJSONMalformed = []byte(`{
	 "catalogs: [
	  "testing"
	 ],
	 "optional_installs": [
	  "AdobePhotoshop"
	 ],
	 "managed_installs": [
	  "munkitools",
	  "munkitools_core",
	  "munkitools_launchd",
	  "sal_scripts"
	 ]
	}`)
