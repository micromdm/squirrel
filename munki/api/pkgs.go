package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/groob/ape/datastore"
	"github.com/julienschmidt/httprouter"
)

func handlePkgsCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)

		// process the multipart form
		filename, file, err := processFileUpload(r)
		if err != nil {
			respondError(rw, http.StatusBadRequest, accept, err)
			return
		}

		// save to datastore
		if err = db.AddPkg(filename, file); err != nil {
			respondError(rw, errStatus(err), accept,
				fmt.Errorf("Failed to save file: %v", err))
			return
		}
	}
}

func handlePkgsDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		err := db.DeletePkg(name)
		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}

		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to delete pkgsinfo from the datastore: %v", err))
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
