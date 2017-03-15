package cli

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/micromdm/squirrel/version"
)

func Serve() {
	serveCMD := flag.NewFlagSet("ca", flag.ExitOnError)
	status := serve(serveCMD)
	os.Exit(status)
}

const authUsername = "squirrel"

func serve(cmd *flag.FlagSet) int {
	var (
		flRepo          = cmd.String("repo", envString("SQUIRREL_MUNKI_REPO_PATH", ""), "path to munki repo")
		flBasicPassword = cmd.String("basic-auth", envString("SQUIRREL_BASIC_AUTH", ""), "http basic auth password for /repo/")
		flTLS           = cmd.Bool("tls", envBool("SQUIRREL_USE_TLS", true), "use https")
		flTLSCert       = cmd.String("tls-cert", envString("SQUIRREL_TLS_CERT", ""), "path to TLS certificate")
		flTLSKey        = cmd.String("tls-key", envString("SQUIRREL_TLS_KEY", ""), "path to TLS private key")
		flTLSAuto       = cmd.String("tls-domain", envString("SQUIRREL_AUTO_TLS_DOMAIN", ""), "Automatically fetch certs from Let's Encrypt")
	)
	cmd.Parse(os.Args[2:])
	mux := http.NewServeMux()
	mux.Handle("/repo/", repoHandler(*flBasicPassword, *flRepo))
	mux.Handle("/version", version.Handler())
	mux.Handle("/healthz", healthz(*flRepo))

	srv := &http.Server{
		Addr:              ":https",
		Handler:           mux,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Minute,
		MaxHeaderBytes:    1 << 18, // 0.25 MB
		TLSConfig:         tlsConfig(),
	}

	printMunkiHeadersHelp(*flBasicPassword)
	if !*flTLS {
		log.Fatal(http.ListenAndServe(":80", mux))
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
