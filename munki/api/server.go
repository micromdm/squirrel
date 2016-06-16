package api

import (
	"log"
	"net/http"

	"github.com/micromdm/squirrel/munki/datastore"
)

type config struct {
	db datastore.Datastore
	// all filesystem interactions use the datastore,
	// the repoPath is needed to serve static files over http
	repoPath  string
	jwtAuth   bool
	basicAuth bool
	jwtSecret string // jwt signing key
	mux       http.Handler
}

// NewServer returns an http Handler
func NewServer(options ...func(*config) error) http.Handler {
	conf := &config{}
	for _, option := range options {
		if err := option(conf); err != nil {
			log.Fatal(err)
		}
	}
	if conf.db == nil {
		log.Fatal("No datastore configured")
	}
	conf.mux = router(conf)
	return conf.mux
}

// SimpleRepo adds a file based backend
func SimpleRepo(path string) func(*config) error {
	return func(c *config) error {
		repo := &datastore.SimpleRepo{Path: path}
		go repo.WatchCatalogs()
		c.db = repo
		c.repoPath = repo.Path
		return nil
	}
}

// PushRepo adds a file based backend with git push on commit
func PushRepo(path string) func(*config) error {
	return func(c *config) error {
		db := &datastore.SimpleRepo{Path: path}
		repo := &datastore.GitRepo{
			Path:      path,
			Datastore: db,
		}
		go db.WatchCatalogs()
		c.db = repo
		c.repoPath = repo.Path
		return nil
	}
}

// JWTAuth enables JWT authentication middleware
func JWTAuth(secret string) func(*config) error {
	return func(c *config) error {
		c.jwtSecret = secret
		c.jwtAuth = true
		return nil
	}
}

// BasicAuth enables basic authentication for the API
func BasicAuth() func(*config) error {
	return func(c *config) error {
		c.basicAuth = true
		return nil
	}
}

// ServerOptions is a slice of config functions
type ServerOptions []func(*config) error
