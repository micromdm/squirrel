package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
	"github.com/groob/plist"
	"github.com/micromdm/squirrel/munki/datastore"
)

func Handler(repoPath string) http.Handler {
	repo := &datastore.SimpleRepo{Path: repoPath}
	svc := &service{repo: repo}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/munki/pkgs", svc.ListPkgs).Methods("GET")
	r.HandleFunc("/api/v1/munki/pkgs/{pkgpath:.*}", svc.UploadPkg).Methods("POST")
	r.HandleFunc("/api/v1/munki/pkgsinfo", svc.ListPkgsinfo).Methods("GET")
	r.HandleFunc("/api/v1/munki/pkgsinfo/{pkgpath:.*}", svc.SavePkgsinfo).Methods("PUT")
	r.HandleFunc("/api/v1/munki/manifests", svc.ListManifests).Methods("GET")
	r.HandleFunc("/api/v1/munki/manifests/{pkgpath:.*}", svc.GetManifest).Methods("GET")
	r.HandleFunc("/api/v1/munki/catalogs", svc.ListCatalogs).Methods("GET")
	r.HandleFunc("/api/v1/munki/catalogs/{pkgpath:.*}", svc.GetCatalog).Methods("GET")
	return r
}

type service struct {
	repo *datastore.SimpleRepo
}

func (svc *service) GetCatalog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path, ok := vars["pkgpath"]
	if !ok {
		http.Error(w, "no pkgpath provided", http.StatusBadRequest)
		return
	}
	cata, err := svc.repo.GetFile("catalogs", path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(cata)
}

func (svc *service) GetManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path, ok := vars["pkgpath"]
	if !ok {
		http.Error(w, "no pkgpath provided", http.StatusBadRequest)
		return
	}
	cata, err := svc.repo.GetFile("manifests", path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(cata)
}

func (svc *service) ListCatalogs(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(dump))

	manifests, err := svc.repo.AllCatalogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type list struct {
		Filename string `plist:"filename"`
	}
	var resp []list

	for _, i := range manifests {
		resp = append(resp, list{i})
	}

	if err := plist.NewEncoder(w).Encode(&resp); err != nil {
		log.Println(err)
	}
}

func (svc *service) ListManifests(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(dump))

	manifests, err := svc.repo.AllManifests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type list struct {
		Filename string `plist:"filename"`
	}
	var resp []list

	for _, i := range manifests {
		resp = append(resp, list{i})
	}

	if err := plist.NewEncoder(w).Encode(&resp); err != nil {
		log.Println(err)
	}
}

func (svc *service) ListPkgsinfo(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(dump))

	pkgsinfo, err := svc.repo.AllPkgsinfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type list struct {
		Filename string `plist:"filename"`
	}
	var resp []list

	for _, i := range pkgsinfo {
		resp = append(resp, list{i})
	}

	if err := plist.NewEncoder(w).Encode(&resp); err != nil {
		log.Println(err)
	}
}

func (svc *service) SavePkgsinfo(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(dump))

	vars := mux.Vars(r)
	path, ok := vars["pkgpath"]
	if !ok {
		http.Error(w, "no pkgpath provided", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	path = "pkgsinfo/" + path

	if err := svc.repo.SavePkgsinfo(path, r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (svc *service) UploadPkg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path, ok := vars["pkgpath"]
	if !ok {
		http.Error(w, "no pkgpath provided", http.StatusBadRequest)
		return
	}

	path = "pkgs/" + path

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("filedata")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := svc.repo.SavePkg(path, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (svc *service) ListPkgs(w http.ResponseWriter, r *http.Request) {
	pkgs, err := svc.repo.AllPkgs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp = struct {
		Filename []string `plist:"filename"`
	}{
		Filename: pkgs,
	}

	if err := plist.NewEncoder(w).Encode(&resp); err != nil {
		log.Println(err)
	}
}
