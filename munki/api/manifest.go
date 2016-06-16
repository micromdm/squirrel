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

func handleManifestsList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		manifests, err := db.AllManifests()
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to fetch manifest list from the datastore: %v", err))
			return
		}
		respondOK(rw, manifests, accept)
	}
}

func handleManifestsShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		manifest, err := db.Manifest(name)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to fetch manifest from the datastore: %v", err))
			return
		}
		respondOK(rw, manifest, accept)
	}
}

func handleManifestsCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		var manifest *models.Manifest
		var payload struct {
			Filename string `plist:"filename" json:"filename"`
			*models.Manifest
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
				errors.New("the name field is required to create a manifest"))
			return
		}

		// If the body contains a valid manifest, create it
		manifest = payload.Manifest
		manifest.Filename = payload.Filename

		// create manifest in datastore
		_, err = db.NewManifest(payload.Filename)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to create new manifest: %v", err))
			return
		}

		// save the manifest
		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}

		respondCreated(rw, manifest, accept)
	}
}

func handleManifestsDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		err := db.DeleteManifest(name)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to delete manifest from the datastore: %v", err))
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func handleManifestsUpdate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		manifest, err := db.Manifest(name)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to delete manifest from the datastore: %v", err))
			return
		}

		// decode payload from request body.
		payload := &models.ManifestPayload{}
		err = decodeRequest(r, payload)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to decode request payload: %v", err))
			return
		}

		// update manifest from payload fields.
		manifest.UpdateFromPayload(payload)

		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}

		// manifest updated ok, respond
		respondOK(rw, manifest, accept)
	}
}

func handleManifestsReplace(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		manifest, err := db.Manifest(name)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to upddate manifest: %v", err))
			return
		}

		payload := &models.Manifest{}
		err = decodeRequest(r, payload)
		if err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to decode request payload: %v", err))
			return
		}

		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}

		manifest = payload
		manifest.Filename = name

		// manifest updated ok, respond
		respondOK(rw, manifest, accept)
	}
}
