package cli

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/groob/finalizer/logutil"
	"github.com/micromdm/go4/version"

	"github.com/micromdm/squirrel/storage/gcs"
	"github.com/micromdm/squirrel/storage/s3"
)

func Serve() {
	serveCMD := flag.NewFlagSet("ca", flag.ExitOnError)
	status := serve(serveCMD)
	os.Exit(status)
}

const authUsername = "squirrel"

func serve(cmd *flag.FlagSet) int {
	var (
		flProvider       = cmd.String("provider", envString("SQUIRREL_MUNKI_REPO_PROVIDER", "filesystem"), "munki repo provider: GCS or filesystem")
		flGCSCredentials = cmd.String("gcs-credentials", envString("SQUIRREL_GCS_CREDENTIALS", ""), "path to google cloud storage credentials json")
		flRepo           = cmd.String("repo", envString("SQUIRREL_MUNKI_REPO_PATH", ""), "path to munki repo")
		flBasicPassword  = cmd.String("basic-auth", envString("SQUIRREL_BASIC_AUTH", ""), "http basic auth password for /repo/")
		flTLS            = cmd.Bool("tls", envBool("SQUIRREL_USE_TLS", true), "use https")
		flTLSCert        = cmd.String("tls-cert", envString("SQUIRREL_TLS_CERT", ""), "path to TLS certificate")
		flTLSKey         = cmd.String("tls-key", envString("SQUIRREL_TLS_KEY", ""), "path to TLS private key")
		flTLSAuto        = cmd.String("tls-domain", envString("SQUIRREL_AUTO_TLS_DOMAIN", ""), "Automatically fetch certs from Let's Encrypt")
		flLogFormat      = cmd.String("log-format", envString("SQUIRREL_LOG_FORMAT", "logfmt"), "Enable structured logging. Supported formats: logfmt, json")
		flSilent         = cmd.Bool("no-help", envBool("SQUIRREL_NO_HELP_TEXT", false), "Silence help text to avoid displaying Auth headers in log.")
	)
	cmd.Parse(os.Args[2:])

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
		logger = kitlog.With(logger, "component", "http")
		logger = level.Info(logger)

	}
	h := logutil.NewHTTPLogger(logger).Middleware(mux)

	srv := &http.Server{
		Addr:              ":https",
		Handler:           h,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Minute,
		MaxHeaderBytes:    1 << 18, // 0.25 MB
		TLSConfig:         tlsConfig(),
	}

	if !*flSilent {
		printMunkiHeadersHelp(*flBasicPassword)
	}
	if !*flTLS {
		log.Fatal(http.ListenAndServe(":8080", h))
		return 0
	}

	if *flTLSAuto != "" {
		serveACME(srv, *flTLSAuto)
		return 0
	}

	tlsFromFile := *flTLSAuto == "" || (*flTLSCert != "" && *flTLSKey != "")
	if tlsFromFile {
		serveTLS(srv, *flTLSCert, *flTLSKey)
		return 0
	}

	return 0
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

func serveTLS(server *http.Server, certPath, keyPath string) {
	redirectTLS()
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}

func serveACME(server *http.Server, domain string) {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache("./certificates"),
	}
	server.TLSConfig.GetCertificate = m.GetCertificate
	redirectTLS()
	log.Fatal(server.ListenAndServeTLS("", ""))
}

// redirects port 80 to port 443
func redirectTLS() {
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}
	go func() { log.Fatal(srv.ListenAndServe()) }()
}

func authMW(next http.Handler, repoPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, password, ok := r.BasicAuth()
		if !ok || password != repoPassword {
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

func tlsConfig() *tls.Config {
	cfg := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	return cfg
}

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string, def bool) bool {
	if env := os.Getenv(key); env == "true" {
		return true
	}
	return def
}

func basicAuth(password string) string {
	auth := authUsername + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
