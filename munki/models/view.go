package models

import (
	"encoding/json"
	"errors"

	"github.com/groob/plist"
)

var (
	// ErrNoData means a resource was empty
	// Equivalent to HTTP 404 response
	ErrNoData = errors.New("No Data")

	// ErrUnsupportedMedia is used when a wrong accept header or media type is specified.
	ErrUnsupportedMedia = errors.New("Unsupported accept encoding type")
)

// Viewer defines a view interface
type Viewer interface {
	View(accept string) (*Response, error)
}

// Response is the API response. All models should implement the Viewer interface and return a response
type Response struct {
	Data []byte
}

// marshals uses a specified accept heading and returns a response
func marshal(view interface{}, accept string) (*Response, error) {
	var data []byte
	var err error

	switch accept {
	case "application/json":
		data, err = json.MarshalIndent(view, "", " ")
		if err != nil {
			return nil, err
		}
	case "application/xml":
		data, err = plist.MarshalIndent(view, "  ")
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedMedia
	}

	resp := &Response{
		Data: data,
	}

	return resp, nil
}
