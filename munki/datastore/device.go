package datastore

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/groob/plist"
	"github.com/micromdm/squirrel/munki/munki"
)

// type DeviceStore interface {
// 	AllDevices() ([]Device, error)
// 	Device(name string) (*Device, error)
// 	NewDevice(name string) (*Device, error)
// 	SaveDevice(name string, dev *Device) error
// 	DeleteDevice(name string) error
// }

// AllDevices returns a list of all devices in a Munki repo
func (r *SimpleRepo) AllDevices() ([]munki.Device, error) {
	devices, err := loadDevices(r.Path)
	if err != nil {
		return nil, err
	}
	r.updateDeviceIndex(devices)
	return devices, nil
}

// Device returns a single device from the repo
func (r *SimpleRepo) Device(name string) (*munki.Device, error) {
	devices, err := loadDevices(r.Path)
	if err != nil {
		return nil, err
	}
	r.updateDeviceIndex(devices)
	dev, ok := r.indexDevices[name]
	if !ok {
		return nil, ErrNotFound
	}
	return &dev, nil
}

// NewDevice creates and returns a new device in the repo
func (r *SimpleRepo) NewDevice(name string) (*munki.Device, error) {
	dev := &munki.Device{SerialNumber: name}
	devPath := fmt.Sprintf("%v/devices/%v", r.Path, name)
	err := createFile(devPath)
	return dev, err
}

// SaveDevice saves a device in the datastore
func (r *SimpleRepo) SaveDevice(name string, device *munki.Device) error {
	if name == "" {
		return errors.New("must specify a device name(serial number)")
	}
	devPath := fmt.Sprintf("%v/devices/%v", r.Path, name)
	file, err := os.OpenFile(devPath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := plist.NewEncoder(file).Encode(device); err != nil {
		return err
	}
	return nil
}

// DeleteDevice removes a device file from the repository
func (r *SimpleRepo) DeleteDevice(name string) error {
	devPath := fmt.Sprintf("%v/devices/%v", r.Path, name)
	return deleteFile(devPath)
}

func loadDevices(path string) ([]munki.Device, error) {
	devicePath := fmt.Sprintf("%v/devices", path)
	devices := []munki.Device{}
	err := filepath.Walk(devicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if !info.IsDir() {
			var dev munki.Device
			if err := plist.NewDecoder(file).Decode(dev); err != nil {
				log.Printf("git-repo: failed to decode %v, skipping \n", info.Name())
				return nil
			}
			devices = append(devices, dev)
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return devices, nil
}

//
func (r *SimpleRepo) updateDeviceIndex(devices []munki.Device) {
	r.indexDevices = make(map[string]munki.Device, len(devices))
	for _, dev := range devices {
		r.indexDevices[dev.SerialNumber] = dev
	}
}
