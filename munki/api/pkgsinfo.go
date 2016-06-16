package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/groob/ape/datastore"
	"github.com/groob/ape/models"
	"github.com/julienschmidt/httprouter"
)

func handlePkgsinfoList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		pkgsinfos, err := db.AllPkgsinfos()
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to fetch pkgsinfo list from the datastore: %v", err))
			return
		}
		// apply any filters
		pkgsinfos = applyPkgsinfoFilters(pkgsinfos, r.URL.Query())
		respondOK(rw, pkgsinfos, accept)
	}
}

func handlePkgsinfoShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		pkgsinfo, err := db.Pkgsinfo(name)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to fetch pkgsinfo from the datastore: %v", err))
			return
		}
		respondOK(rw, pkgsinfo, accept)
	}
}

func handlePkgsinfoCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		var pkgsinfo *models.PkgsInfo
		var payload struct {
			Filename string `plist:"filename" json:"filename"`
			*models.PkgsInfo
		}
		err := decodeRequest(r, &payload)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to decode request payload: %v", err))
			return
		}

		// filename is required in the payload
		if payload.Filename == "" {
			respondError(rw, http.StatusBadRequest, accept,
				errors.New("The filename field is required to create a pkgsinfo"))
			return
		}

		// If the body contains a valid pkgsinfo, create it
		pkgsinfo = payload.PkgsInfo
		pkgsinfo.Filename = payload.Filename

		// create pkgsinfo in datastore
		_, err = db.NewPkgsinfo(payload.Filename)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to create new pkgsinfo: %v", err))
			return
		}

		// save the pkgsinfo
		if err := db.SavePkgsinfo(pkgsinfo); err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to save pkgsinfo: %v", err))
			return
		}

		respondCreated(rw, pkgsinfo, accept)
	}
}

func handlePkgsinfoDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		err := db.DeletePkgsinfo(name)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to delete pkgsinfo from the datastore: %v", err))
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func handlePkgsinfoReplace(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		pkgsinfo, err := db.Pkgsinfo(name)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to upddate pkgsinfo: %v", err))
			return
		}

		payload := &models.PkgsInfo{}
		err = decodeRequest(r, payload)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to decode request payload: %v", err))
			return
		}

		if err := db.SavePkgsinfo(pkgsinfo); err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to save pkginfo: %v", err))
			return
		}

		pkgsinfo = payload
		pkgsinfo.Filename = name

		// manifest updated ok, respond
		respondOK(rw, pkgsinfo, accept)
	}
}
