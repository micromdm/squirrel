package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func router(conf *config) http.Handler {
	mux := http.NewServeMux()
	repo := http.StripPrefix("/repo/", http.FileServer(http.Dir(conf.repoPath)))
	mux.Handle("/repo/", repo)
	addAuthOptions(mux, conf)
	return mux
}

func addAuthOptions(mux *http.ServeMux, conf *config) {
	// add login handler if basicAuth is enabled
	if conf.basicAuth {
		mux.Handle("/api/login", handleBasicAuth())
		if !conf.jwtAuth || conf.jwtSecret == "" {
			log.Fatal("basic auth needs jwt enabled")
		}
	}
	// add token endpoints and authmiddleware if
	// jwtAuth is enabled
	if conf.jwtAuth {
		mux.Handle("/api/token", handleTokenRefresh())
		mux.Handle("/api/", authMiddleware(apiRouter(conf)))
	} else {
		mux.Handle("/api/", apiRouter(conf))
	}
}

func apiRouter(conf *config) *httprouter.Router {
	router := httprouter.New()
	// handle manifest actions
	router.GET("/api/manifests", handleManifestsList(conf.db))
	router.POST("/api/manifests", handleManifestsCreate(conf.db))
	router.GET("/api/manifests/*name", handleManifestsShow(conf.db))
	router.PUT("/api/manifests/*name", handleManifestsReplace(conf.db))
	router.PATCH("/api/manifests/*name", handleManifestsUpdate(conf.db))
	router.DELETE("/api/manifests/*name", handleManifestsDelete(conf.db))
	// handle pkgsinfo actions
	router.GET("/api/pkgsinfo", handlePkgsinfoList(conf.db))
	router.POST("/api/pkgsinfo", handlePkgsinfoCreate(conf.db))
	router.GET("/api/pkgsinfo/*name", handlePkgsinfoShow(conf.db))
	router.PUT("/api/pkgsinfo/*name", handlePkgsinfoReplace(conf.db))
	// TODO: Create PATCH method for pkgsinfo
	// router.PATCH("/api/pkgsinfo/*name", handlePkgsInfoUpdate(conf.db))
	router.DELETE("/api/pkgsinfo/*name", handlePkgsinfoDelete(conf.db))
	// handle pkgs actions
	router.POST("/api/pkgs", handlePkgsCreate(conf.db))
	router.DELETE("/api/pkgs/*name", handlePkgsDelete(conf.db))
	return router
}
