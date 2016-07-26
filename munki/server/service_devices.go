package munkiserver

import (
	"github.com/micromdm/squirrel/munki/munki"
	"golang.org/x/net/context"
)

// DeviceService describes the actions of a munki server
type DeviceService interface {
	ListDevices(ctx context.Context) ([]munki.Device, error)
	ShowDevice(ctx context.Context, name string) (*munki.Device, error)
	CreateDevice(ctx context.Context, name string, dev *munki.Device) (*munki.Device, error)
	ReplaceDevice(ctx context.Context, name string, dev *munki.Device) (*munki.Device, error)
	DeleteDevice(ctx context.Context, name string) error
}

func (svc service) ListDevices(ctx context.Context) ([]munki.Device, error) {
	return svc.repo.AllDevices()
}

func (svc service) ShowDevice(ctx context.Context, name string) (*munki.Device, error) {
	return svc.repo.Device(name)
}

func (svc service) CreateDevice(ctx context.Context, name string, dev *munki.Device) (*munki.Device, error) {
	_, err := svc.repo.NewDevice(name)
	if err != nil {
		return nil, err
	}
	if err := svc.repo.SaveDevice(name, dev); err != nil {
		return nil, err
	}
	return dev, nil
}

func (svc service) ReplaceDevice(ctx context.Context, name string, dev *munki.Device) (*munki.Device, error) {
	if err := svc.repo.DeleteDevice(name); err != nil {
		return nil, err
	}
	return svc.CreateDevice(ctx, name, dev)
}

func (svc service) DeleteDevice(ctx context.Context, name string) error {
	return svc.repo.DeleteDevice(name)
}
