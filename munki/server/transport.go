package munkiserver

import (
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/groob/plist"
	"github.com/micromdm/squirrel/munki/datastore"

	"golang.org/x/net/context"
)

var (
	// ErrEmptyRequest is returned if the request body is empty
	errEmptyRequest = errors.New("request must contain all required fields")
	errBadRouting   = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func decodeListManifestsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listManifestsRequest{}, nil
}

func decodeShowManifestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	return showManifestRequest{Path: path}, nil
}

// ServiceHandler creates an HTTP handler for the munki Service
func ServiceHandler(ctx context.Context, svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(updateContext),
	}
	listManifestsHandler := kithttp.NewServer(
		ctx,
		makeListManifestsEndpoint(svc),
		decodeListManifestsRequest,
		encodeResponse,
		opts...,
	)
	showManifestHandler := kithttp.NewServer(
		ctx,
		makeShowManifestEndpoint(svc),
		decodeShowManifestRequest,
		encodeResponse,
		opts...,
	)
	r := mux.NewRouter()
	// manifests
	r.Handle("/api/v1/manifests/{path}", showManifestHandler).Methods("GET")
	r.Handle("/api/v1/manifests", listManifestsHandler).Methods("GET")
	return r
}

func updateContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, "mediaType", acceptHeader(r))
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

type errorer interface {
	error() error
}

type statuser interface {
	status() int
}

type subsetEncoder interface {
	subset() interface{}
}

func encodeJSON(w http.ResponseWriter, from interface{}) error {
	data, err := json.MarshalIndent(from, "", "  ")
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}

func encodePLIST(w http.ResponseWriter, from interface{}) error {
	enc := plist.NewEncoder(w)
	enc.Indent("  ")
	return enc.Encode(from)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	mediaType := ctx.Value("mediaType").(string)
	setContentType(w, mediaType)
	// for success responses
	if e, ok := response.(statuser); ok {
		w.WriteHeader(e.status())
		if e.status() == http.StatusNoContent {
			return nil
		}
	}

	// check if this is a collection
	if e, ok := response.(subsetEncoder); ok {
		response = e.subset()
	}
	if mediaType == "application/xml" {
		return encodePLIST(w, response)
	}
	return encodeJSON(w, response)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	mediaType := ctx.Value("mediaType").(string)
	setContentType(w, mediaType)
	w.WriteHeader(codeFrom(err))
	errData := map[string]interface{}{
		"error": err.Error(),
	}
	if mediaType == "application/xml" {
		encodePLIST(w, errData)
		return
	}
	encodeJSON(w, errData)
}

func codeFrom(err error) int {
	switch err {
	case datastore.ErrNotFound:
		return http.StatusNotFound
	default:
		if e, ok := err.(httptransport.Error); ok {
			switch e.Domain {
			case httptransport.DomainDecode:
				return http.StatusBadRequest
			case httptransport.DomainDo:
				return http.StatusServiceUnavailable
			default:
				return http.StatusInternalServerError
			}
		}
		return http.StatusInternalServerError
	}
}
