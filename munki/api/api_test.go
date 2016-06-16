package api

import (
	"ape/datastore"
	"ape/models"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var (
	// check if MockRepo implements datastore
	_ datastore.Datastore = (*MockRepo)(nil)
	// test client
	client = NewTestClient()

	xmlMedia  = "application/xml; charset=utf-8"
	jsonMedia = "application/json; charset=utf-8"
)

type MockRepo struct {
	Path           string
	indexManifests map[string]*models.Manifest
	indexPkgsinfo  map[string]*models.PkgsInfo
}

func newMockRepo(path string) func(*config) error {
	return func(c *config) error {
		repo := &MockRepo{
			Path:           path,
			indexManifests: mockManifests(),
		}
		c.db = repo
		c.repoPath = repo.Path
		return nil
	}
}

func newMockServer(options ...func(*config) error) http.Handler {
	conf := &config{}
	for _, option := range options {
		if err := option(conf); err != nil {
			log.Fatal(err)
		}
	}
	if conf.db == nil {
		log.Fatal("No datastore configured")
	}
	// api routes
	conf.mux = apiRouter(conf)

	return conf.mux
}

func mockManifests() map[string]*models.Manifest {
	var manifests = map[string]*models.Manifest{}
	siteDefault := &models.Manifest{
		Filename: "site_default",
		Catalogs: []string{"production"},
	}
	siteDefault.ManagedInstalls = []string{"Firefox", "AdobeFlashPlayer"}
	manifests["site_default"] = siteDefault
	return manifests
}

func (m *MockRepo) AllManifests() (*models.ManifestCollection, error) {
	manifests := models.ManifestCollection{}
	for _, manifest := range m.indexManifests {
		manifests = append(manifests, manifest)
	}
	return &manifests, nil
}

func (m *MockRepo) Manifest(name string) (*models.Manifest, error) {
	if manifest, ok := m.indexManifests[name]; ok {
		return manifest, nil
	}
	return nil, datastore.ErrNotFound
}

func (m *MockRepo) NewManifest(name string) (*models.Manifest, error) {
	// check if manifest already exists, and return a new one if it doesn't
	if _, ok := m.indexManifests[name]; ok {
		return nil, datastore.ErrExists
	}
	return &models.Manifest{}, nil
}
func (m *MockRepo) SaveManifest(manifest *models.Manifest) error {
	m.indexManifests[manifest.Filename] = manifest
	return nil
}
func (m *MockRepo) DeleteManifest(name string) error {
	if _, ok := m.indexManifests[name]; !ok {
		return datastore.ErrNotFound
	}
	delete(m.indexManifests, name)
	return nil
}
func (m *MockRepo) AllPkgsinfos() (*models.PkgsInfoCollection, error) { return nil, nil }
func (m *MockRepo) Pkgsinfo(name string) (*models.PkgsInfo, error)    { return nil, nil }
func (m *MockRepo) NewPkgsinfo(name string) (*models.PkgsInfo, error) { return nil, nil }
func (m *MockRepo) SavePkgsinfo(manifest *models.PkgsInfo) error      { return nil }
func (m *MockRepo) DeletePkgsinfo(name string) error                  { return nil }
func (m *MockRepo) AddPkg(filename string, body io.Reader) error      { return nil }
func (m *MockRepo) DeletePkg(name string) error                       { return nil }

func newTestServer() *httptest.Server {
	apiHandler := newMockServer(newMockRepo("mockrepo"))
	server := httptest.NewServer(apiHandler)
	return server
}

/* HTTP test code */

type TestClient struct {
	client *http.Client
	server *httptest.Server

	// Base URL for API requests.
	BaseURL *url.URL
}

func NewTestClient() *TestClient {
	client := &TestClient{client: http.DefaultClient}
	client.server = newTestServer()
	client.BaseURL, _ = url.Parse(client.server.URL)
	client.BaseURL.Path = "api/"
	return client
}

// create testclient request
func (c *TestClient) NewRequest(endpoint, resource, mediaType, method string) (*http.Request, error) {
	var urlStr string
	if resource != "" {
		urlStr = endpoint + "/" + resource
	} else {
		urlStr = endpoint
	}
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil

}

// run the request
func (c *TestClient) Do(req *http.Request, into interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *TestClient) teardown() {
	c.server.Close()
}

// a face io.ReadCloser for constructing request Body
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func setHeader(req *http.Request, contentType string) {
	switch contentType {
	case jsonMedia:
		req.Header.Set("Content-Type", jsonMedia)
		req.Header.Set("Accept", jsonMedia)
	case xmlMedia:
		req.Header.Set("Content-Type", xmlMedia)
		req.Header.Set("Accept", xmlMedia)
	}
}

/*	API Test Code */

// test main
func TestMain(m *testing.M) {
	retCode := m.Run()
	// call the tearedown function to close the server
	client.teardown()
	// exit
	os.Exit(retCode)
}
