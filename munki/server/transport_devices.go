package munkiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

func decodeListDevicesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listDevicesRequest{}, nil
}

func decodeShowDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	return showDeviceRequest{Name: path}, nil
}

func decodeCreateDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeDeleteDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	return deleteDeviceRequest{SerialNumber: path}, nil
}

func decodeReplaceDeviceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request replaceDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request.Device); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	path, ok := vars["path"]
	if !ok {
		return nil, errBadRouting
	}
	request.SerialNumber = path
	return request, nil
}
