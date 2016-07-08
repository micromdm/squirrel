package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	kitlog "github.com/go-kit/kit/log"

	"github.com/micromdm/squirrel/munki/datastore"
	"github.com/micromdm/squirrel/munki/server"
)

const usage = "usage: MUNKI_REPO_PATH= SQUIRREL_HTTP_LISTEN_PORT= ape -repo MUNKI_REPO_PATH -port SQUIRREL_HTTP_LISTEN_PORT"

func main() {
	var (
		flRepo      = flag.String("repo", envString("SQUIRREL_MUNKI_REPO_PATH", ""), "path to munki repo")
		flPort      = flag.String("port", envString("SQUIRREL_HTTP_LISTEN_PORT", ""), "port to listen on")
		flBasic     = flag.Bool("basic", envBool("SQUIRREL_BASIC_AUTH"), "enable basic auth")
		flJWT       = flag.Bool("jwt", envBool("SQUIRREL_JWT_AUTH"), "enable jwt authentication for api calls")
		flJWTSecret = flag.String("jwt-signing-key", envString("SQUIRREL_JWT_SIGNING_KEY", ""), "jwt signing key")
		flTLS       = flag.Bool("tls", envBool("SQUIRREL_USE_TLS"), "use https")
		flTLSCert   = flag.String("tls-cert", envString("SQUIRREL_TLS_CERT", ""), "path to TLS certificate")
		flTLSKey    = flag.String("tls-key", envString("SQUIRREL_TLS_KEY", ""), "path to TLS private key")
	)
	*flTLS = true
	flag.Parse()
	if *flRepo == "" {
		flag.Usage()
		log.Fatal(usage)
	}

	// create the folders if they don't yet exist
	checkRepo(*flRepo)

	// validate port flag
	if *flPort == "" {
		port := defaultPort(*flTLS)
		log.Printf("no port flag specified. Using %v by default", port)
		*flPort = port
	}

	if *flTLS {
		checkTLSFlags(*flTLSKey, *flTLSCert)
	}

	// validate JWT flags
	if *flJWT {
		checkJWTFlags(*flJWTSecret)
	}
	// validate basic auth
	if *flBasic && !*flJWT {
		log.Fatal("Basic Authentication is used to issue JWT Tokens. You must enable JWT as well")
	}

	var repo datastore.Datastore
	{
		repo = &datastore.SimpleRepo{Path: *flRepo}

	}

	var err error
	var svc munkiserver.Service
	{
		svc, err = munkiserver.NewService(repo)
		if err != nil {
			log.Fatal(err)
		}
	}

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.NewContext(logger).With("ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.NewContext(logger).With("caller", kitlog.DefaultCaller)
	}

	ctx := context.Background()
	var h http.Handler
	{
		h = munkiserver.ServiceHandler(ctx, svc, logger)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", h)
	mux.Handle("/repo/", http.StripPrefix("/repo/", http.FileServer(http.Dir(*flRepo))))

	port := fmt.Sprintf(":%v", *flPort)

	if *flTLS {
		log.Fatal(http.ListenAndServeTLS(port, *flTLSCert, *flTLSKey, mux))
	} else {
		log.Fatal(http.ListenAndServe(port, mux))
	}
}

func defaultPort(tls bool) string {
	if tls {
		return "443"
	}
	return "80"
}

func checkJWTFlags(secret string) {
	if secret == "" {
		log.Fatal("You must provide a signing key to enable JWT authentication")
	}
}

func checkTLSFlags(keypath, certpath string) {
	if keypath == "" || certpath == "" {
		log.Fatal("You must provide a valid path to a TLS cert and key")
	}
}

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string) bool {
	if env := os.Getenv(key); env == "true" {
		return true
	}
	return false
}

func createDir(path string) {
	if !dirExists(path) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("%v must exits", path)
		}
	}
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func checkRepo(repoPath string) {
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo/", repoPath)
	createDir(pkgsinfoPath)

	manifestPath := fmt.Sprintf("%v/manifests/", repoPath)
	createDir(manifestPath)

	pkgsPath := fmt.Sprintf("%v/pkgs/", repoPath)
	createDir(pkgsPath)

	catalogsPath := fmt.Sprintf("%v/catalogs/", repoPath)
	createDir(catalogsPath)
}
