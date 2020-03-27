package main

import (
	"crypto/subtle"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/groob/finalizer/logutil"
	"github.com/micromdm/go4/env"
	"github.com/micromdm/go4/httputil"
	"github.com/micromdm/go4/version"
	"github.com/pkg/errors"

	"github.com/micromdm/squirrel/storage/gcs"
	"github.com/micromdm/squirrel/storage/s3"
)

func runServe(args []string) error {
	flagset := flag.NewFlagSet("squirrel", flag.ExitOnError)
	var (
		flConfigPath     = flagset.String("config-path", env.String("SQUIRREL_AUTOCERT_CACHE_PATH", "/var/micromdm/squirrel"), "path to autocert cache directory")
		flProvider       = flagset.String("provider", env.String("SQUIRREL_MUNKI_REPO_PROVIDER", "filesystem"), "munki repo provider: GCS or filesystem")
		flGCSCredentials = flagset.String("gcs-credentials", env.String("SQUIRREL_GCS_CREDENTIALS", ""), "path to google cloud storage credentials json")
		flRepo           = flagset.String("repo", env.String("SQUIRREL_MUNKI_REPO_PATH", ""), "path to munki repo")
		flBasicPassword  = flagset.String("basic-auth", env.String("SQUIRREL_BASIC_AUTH", ""), "http basic auth password for /repo/")
		flTLS            = flagset.Bool("tls", env.Bool("SQUIRREL_USE_TLS", true), "use https")
		flTLSCert        = flagset.String("tls-cert", env.String("SQUIRREL_TLS_CERT", ""), "path to TLS certificate")
		flTLSKey         = flagset.String("tls-key", env.String("SQUIRREL_TLS_KEY", ""), "path to TLS private key")
		flTLSDomain      = flagset.String("tls-domain", env.String("SQUIRREL_AUTO_TLS_DOMAIN", ""), "Automatically fetch certs from Let's Encrypt")
		flLogFormat      = flagset.String("log-format", env.String("SQUIRREL_LOG_FORMAT", "logfmt"), "Enable structured logging. Supported formats: logfmt, json")
		flSilent         = flagset.Bool("no-help", env.Bool("SQUIRREL_NO_HELP_TEXT", false), "Silence help text to avoid displaying Auth headers in log.")
		flHTTPDebug      = flagset.Bool("http-debug", false, "enable debug for http(dumps full request)")
		flHTTPAddr       = flagset.String("http-addr", ":https", "http(s) listen address of http server. defaults to :8080 if tls is false")
	)

	flagset.Usage = usageFor(flagset, "squirrel serve [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	var repoFileHandler, healthzHandler http.Handler
	switch strings.ToLower(*flProvider) {
	case "filesystem":
		repoFileHandler = http.StripPrefix("/repo/", http.FileServer(http.Dir(*flRepo)))
		healthzHandler = healthz(*flRepo)
	case "s3":
		s3, err := s3.New(*flRepo)
		if err != nil {
			log.Fatal(err)
		}
		repoFileHandler = http.StripPrefix("/repo/", http.HandlerFunc(s3.FileHandler))
		healthzHandler = s3.HealthzHandler()
	case "gcs":
		if *flGCSCredentials == "" {
			helpText := `
To use squirrel with Google Cloud Storage, you must provide a service account file with 
bucket access. 

First, go to the GCS IAM for project: 
	https://console.cloud.google.com/iam-admin/serviceaccounts

Create a service accont with access to Google Cloud Storage. 
The service account must have at least Read Access to the storage bucket. 
Make sure you download a JSON key for the service account and pass it to 
squirrel with -gcs-credentials=/path/to/service-account.json 
`
			fmt.Println(helpText)
			os.Exit(1)
		}
		gcs, err := gcs.New(*flRepo, *flGCSCredentials)
		if err != nil {
			log.Fatal(err)
		}
		repoFileHandler = http.StripPrefix("/repo/", http.HandlerFunc(gcs.FileHandler))
		healthzHandler = gcs.HealthzHandler()
	default:
		log.Fatalf("unknown gcs provider %s\n", *flProvider)
	}

	if !*flSilent {
		printMunkiHeadersHelp(*flBasicPassword)
	}

	mux := http.NewServeMux()
	mux.Handle("/repo/", authMW(repoFileHandler, *flBasicPassword))
	mux.Handle("/version", version.Handler())
	mux.Handle("/healthz", healthzHandler)

	var logger kitlog.Logger
	{
		w := kitlog.NewSyncWriter(os.Stderr)
		switch *flLogFormat {
		case "json":
			logger = kitlog.NewJSONLogger(w)
		default:
			logger = kitlog.NewLogfmtLogger(w)
		}
		log.SetOutput(kitlog.NewStdlibAdapter(logger))
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger := kitlog.With(logger, "transport", "http")
		logger = level.Info(logger)
	}

	var handler http.Handler
	if *flHTTPDebug {
		handler = httputil.HTTPDebugMiddleware(os.Stdout, true, logger.Log)(mux)
	} else {
		handler = mux
	}
	handler = logutil.NewHTTPLogger(logger).Middleware(handler)

	serveOpts := httputil.Simple(
		*flConfigPath,
		handler,
		*flHTTPAddr,
		*flTLSCert,
		*flTLSKey,
		*flTLS,
		logger,
		*flTLSDomain,
	)

	err := httputil.ListenAndServe(serveOpts...)
	return errors.Wrap(err, "calling ListenAndServe")
}

func printMunkiHeadersHelp(password string) {
	const help = `
	To connect your clients to the server, you will need to set the authentication headers.
	See https://github.com/munki/munki/wiki/Using-Basic-Authentication#configuring-the-clients-to-use-a-password
	for additional details:

	The headers header you should use is:
	
	%s

	To configure manually, use:

	sudo defaults write /Library/Preferences/ManagedInstalls AdditionalHttpHeaders -array "%s"
	`
	auth := basicAuth(password)
	header := fmt.Sprintf("Authorization: Basic %s", auth)
	fmt.Println(fmt.Sprintf(help, header, header))
}

func authMW(next http.Handler, repoPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(password), []byte(repoPassword)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="munki"`)
			http.Error(w, "you need to log in", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func repoHandler(repoPassword string, path string) http.HandlerFunc {
	repo := http.StripPrefix("/repo/", http.FileServer(http.Dir(path)))
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok || password != repoPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="munki"`)
			http.Error(w, "you need to log in", http.StatusUnauthorized)
			return
		}
		repo.ServeHTTP(w, r)
	}
}

func healthz(path string) http.HandlerFunc {
	var healthy bool
	if _, err := os.Stat(path); err == nil {
		healthy = true
	} else {
		log.Printf("healthcheck failed with %s\n", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if !healthy {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

const authUsername = "squirrel"

func basicAuth(password string) string {
	auth := authUsername + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
