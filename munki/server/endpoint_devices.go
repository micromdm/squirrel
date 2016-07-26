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
