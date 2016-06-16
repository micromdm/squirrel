package models

// ManifestPayload represents a payload type
// The payload type is what the client would send over the wire
type ManifestPayload struct {
	Filename          *string      `plist:"filename,omitempty" json:"filename,omitempty"`
	Catalogs          *[]string    `plist:"catalogs,omitempty" json:"catalogs,omitempty"`
	DisplayName       *string      `plist:"display_name,omitempty" json:"display_name,omitempty"`
	IncludedManifests *[]string    `plist:"included_manifests,omitempty" json:"included_manifests,omitempty"`
	Notes             *string      `plist:"notes,omitempty" json:"notes,omitempty"`
	User              *string      `plist:"user,omitempty" json:"user,omitempty"`
	ConditionalItems  *[]condition `plist:"conditional_items,omitempty" json:"conditional_items,omitempty"`
	manifestItemsPayload
}

type manifestItemsPayload struct {
	OptionalInstalls  *[]string `plist:"optional_installs,omitempty" json:"optional_installs,omitempty"`
	ManagedInstalls   *[]string `plist:"managed_installs,omitempty" json:"managed_installs,omitempty"`
	ManagedUninstalls *[]string `plist:"managed_uninstalls,omitempty" json:"managed_uninstalls,omitempty"`
	ManagedUpdates    *[]string `plist:"managed_updates,omitempty" json:"managed_updates,omitempty"`
}
