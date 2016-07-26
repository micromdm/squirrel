package munki

// DeviceStore provides methods for managing device types in a database
type DeviceStore interface {
	AllDevices() ([]Device, error)
	Device(name string) (*Device, error)
	NewDevice(name string) (*Device, error)
	SaveDevice(name string, dev *Device) error
	DeleteDevice(name string) error
}

// Device represents a macOS device
type Device struct {
	SerialNumber     string
	DisplayName      string
	Notes            string
	TemplateManifest string
	User             string
	Catalogs         []string
	DEPStatus        string
}
