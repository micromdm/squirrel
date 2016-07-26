package munkiserver

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/micromdm/squirrel/munki/munki"
	"golang.org/x/net/context"
)

type listDevicesRequest struct{}

type listDevicesResponse struct {
	devices []munki.Device
	Err     error `json:"error,omitempty" plist:"error,omitempty"`
}

// encode a list view which includes the Filename parameter for each manifest
func (r listDevicesResponse) subset() interface{} {
	type devicesView struct {
		SerialNumber string `plist:"serial_number,omitempty" json:"serial_number,omitempty"`
		munki.Device
	}
	var view []devicesView
	for _, item := range r.devices {
		viewItem := devicesView{
			item.SerialNumber,
			item,
		}
		view = append(view, viewItem)
	}

	return view
}

func (r listDevicesResponse) error() error { return r.Err }

func makeListDevicesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		devices, err := svc.ListDevices(ctx)
		return listDevicesResponse{devices: devices, Err: err}, nil
	}
}

type showDeviceRequest struct {
	Name string
}

type showDeviceResponse struct {
	*munki.Device
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r showDeviceResponse) error() error { return r.Err }

func makeShowDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(showDeviceRequest)
		dev, err := svc.ShowDevice(ctx, req.Name)
		return showDeviceResponse{Device: dev, Err: err}, nil
	}
}

type createDeviceRequest struct {
	SerialNumber string `plist:"serial_number,omitempty" json:"serial_number,omitempty"`
	*munki.Device
}

type createDeviceResponse struct {
	*munki.Device
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r createDeviceResponse) error() error { return r.Err }

func makeCreateDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createDeviceRequest)
		dev, err := svc.CreateDevice(ctx, req.SerialNumber, req.Device)
		return createDeviceResponse{Device: dev, Err: err}, nil
	}
}

type deleteDeviceRequest struct {
	SerialNumber string `plist:"serial_number,omitempty" json:"serial_number,omitempty"`
}

type deleteDeviceResponse struct {
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r deleteDeviceResponse) error() error { return r.Err }

func makeDeleteDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteDeviceRequest)
		err := svc.DeleteDevice(ctx, req.SerialNumber)
		return deleteDeviceResponse{Err: err}, nil
	}
}

type replaceDeviceRequest struct {
	SerialNumber string `plist:"serial_number,omitempty" json:"serial_number,omitempty"`
	*munki.Device
}

type replaceDeviceResponse struct {
	*munki.Device
	Err error `json:"error,omitempty" plist:"error,omitempty"`
}

func (r replaceDeviceResponse) error() error { return r.Err }

func makeReplaceDeviceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(replaceDeviceRequest)
		dev, err := svc.ReplaceDevice(ctx, req.SerialNumber, req.Device)
		return replaceDeviceResponse{Device: dev, Err: err}, nil
	}
}
