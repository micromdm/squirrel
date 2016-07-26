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
	SerialNumber     string   `plist:"serial_number" json:"serial_number"`
	HostName         string   `plist:"hostname,omitempty" json:"hostname,omitempty"`
	DisplayName      string   `plist:"display_name,omitempty" json:"display_name,omitempty"`
	Notes            string   `plist:"notes,omitempty" json:"notes,omitempty"`
	TemplateManifest string   `plist:"template_manifest,omitempty" json:"template_manifest,omitempty"`
	User             string   `plist:"user,omitempty" json:"user,omitempty"`
	Catalogs         []string `plist:"catalogs,omitempty" json:"catalogs,omitempty"`
	DEPStatus        string   `plist:"dep_status,omitempty" json:"dep_status,omitempty"`
}
