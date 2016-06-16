package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/groob/ape/datastore"
	"github.com/groob/ape/models"
	"github.com/groob/plist"
)

func respond(rw http.ResponseWriter, body models.Viewer, accept string, status int) {
	setContentType(rw, accept)
	view, err := body.View(accept)
	switch err {
	case nil:
		rw.WriteHeader(status)
		rw.Write(view.Data)
		return
	case models.ErrNoData:
		respondError(rw, http.StatusNotFound, accept, errors.New("Not Found"))
		return
	case models.ErrUnsupportedMedia:
		log.Fatal(err)
		return
	default:
		log.Fatal(err)
		return
	}
}

func respondError(rw http.ResponseWriter, status int, accept string, errs ...error) {
	setContentType(rw, accept)
	resp := &models.ErrorResponse{}
	for _, err := range errs {
		resp.Errors = append(resp.Errors, err.Error())
	}
	view, err := resp.View(accept)
	if err != nil {
		rw.WriteHeader(status)
		log.Println(err)
		return
	}
	rw.WriteHeader(status)
	rw.Write(view.Data)
}

func respondOK(rw http.ResponseWriter, body models.Viewer, accept string) {
	respond(rw, body, accept, http.StatusOK)
}

func respondCreated(rw http.ResponseWriter, body models.Viewer, accept string) {
	respond(rw, body, accept, http.StatusCreated)
}

// if header is not set to json or xml, return json header
func acceptHeader(r *http.Request) string {
	accept := r.Header.Get("Accept")
	switch accept {
	case "application/xml", "application/xml; charset=utf-8":
		return "application/xml"
	default:
		return "application/json"
	}
}

func contentHeader(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/xml", "application/xml; charset=utf-8":
		return "application/xml"
	default:
		return "application/json"
	}
}

func applyPkgsinfoFilters(pkgsinfos *models.PkgsInfoCollection, values url.Values) *models.PkgsInfoCollection {
	if val, ok := values["catalogs"]; ok {
		catalogs := strings.Split(val[0], ",")
		pkgsinfos = pkgsinfos.ByCatalog(catalogs...)
	}

	if _, ok := values["name"]; ok {
		name := values.Get("name")
		pkgsinfos = pkgsinfos.ByName(name)
	}

	return pkgsinfos
}

// set the Content-Type header
func setContentType(rw http.ResponseWriter, accept string) {
	switch accept {
	case "application/xml":
		rw.Header().Set("Content-Type", "application/xml; charset=utf-8")
		return
	default:
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
}

// convert error to status code
// checks an error from the datastore layer and
// returns an appropriate statuscode
func errStatus(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	case io.EOF:
		return http.StatusBadRequest
	case datastore.ErrExists:
		return http.StatusConflict
	case datastore.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// decode an http request into a correct model type
func decodeRequest(r *http.Request, into interface{}) error {
	contentType := contentHeader(r)
	var err error
	switch contentType {
	case "application/xml", "application/xml; charset=utf-8":
		err = plist.NewDecoder(r.Body).Decode(into)
	case "application/json":
		err = json.NewDecoder(r.Body).Decode(into)
	default:
		err = fmt.Errorf("Incorrect Content-Type: %v", contentType)
	}
	return err
}

func processFileUpload(r *http.Request) (string, io.Reader, error) {
	filename := r.FormValue("filename")
	if filename == "" {
		return "", nil, errors.New("Upload form must containt a filename key")
	}
	file, _, err := r.FormFile("filedata")
	// check if file is missing
	if err != nil && err == http.ErrMissingFile {
		return "", nil, errors.New("Filedata must contain a file.")
	}
	// check remaining errors
	if err != nil {
		return "", nil, err
	}
	return filename, file, nil
}
