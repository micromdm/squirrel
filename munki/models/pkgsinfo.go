package models

import "time"

// PkgsinfoStore is an interface for interacting with Pkgsinfo types
type PkgsinfoStore interface {
	AllPkgsinfos() (*PkgsInfoCollection, error)
	Pkgsinfo(name string) (*PkgsInfo, error)
	NewPkgsinfo(name string) (*PkgsInfo, error)
	SavePkgsinfo(manifest *PkgsInfo) error
	DeletePkgsinfo(name string) error
}

// PkgsInfo represents the structure of a pkgsinfo file
type PkgsInfo struct {
	pkgsinfo
	Filename string   `plist:"-" json:"-"`
	Metadata metadata `plist:"_metadata,omitempty" json:"_metadata,omitempty"`
	adobeRelatedItems
}

type pkgsinfo struct {
	Autoremove               bool                  `plist:"autoremove,omitempty" json:"autoremove,omitempty"`
	Catalogs                 []string              `plist:"catalogs,omitempty" json:"catalogs,omitempty"`
	Category                 string                `plist:"category,omitempty" json:"category,omitempty"`
	CopyLocal                bool                  `plist:"copy_local,omitempty" json:"copy_local,omitempty"`
	Description              string                `plist:"description,omitempty" json:"description,omitempty"`
	Developer                string                `plist:"developer,omitempty" json:"developer,omitempty"`
	DisplayName              string                `plist:"display_name,omitempty" json:"display_name,omitempty"`
	ForceInstallAfterDate    time.Time             `plist:"force_install_after_date,omitempty" json:"force_install_after_date,omitempty"`
	IconName                 string                `plist:"icon_name,omitempty" json:"icon_name,omitempty"`
	InstallableCondition     string                `plist:"installable_condition,omitempty" json:"installable_condition,omitempty"`
	InstalledSize            int                   `plist:"installed_size,omitempty" json:"installed_size,omitempty"`
	InstallerItemHash        string                `plist:"installer_item_hash,omitempty" json:"installer_item_hash,omitempty"`
	InstallerItemLocation    string                `plist:"installer_item_location,omitempty" json:"installer_item_location,omitempty"`
	InstallerItemSize        int                   `plist:"installer_item_size,omitempty" json:"installer_item_size,omitempty"`
	InstallerType            string                `plist:"installer_type,omitempty" json:"installer_item_type,omitempty"`
	Installs                 []install             `plist:"installs,omitempty" json:"installs,omitempty"`
	Receipts                 []receipt             `plist:"receipts,omitempty" json:"receipts,omitempty"`
	ItemsToCopy              []itemsToCopy         `plist:"items_to_copy,omitempty" json:"items_to_copy,omitempty"`
	MinimumMunkiVersion      string                `plist:"minimum_munki_version,omitempty" json:"minimum_munki_version,omitempty"`
	MinimumOSVersion         string                `plist:"minimum_os_version,omitempty" json:"minimum_os_version,omitempty"`
	MaximumOSVersion         string                `plist:"maximum_os_version,omitempty" json:"maximum_os_version,omitempty"`
	Name                     string                `plist:"name,omitempty" json:"name,omitempty"`
	Notes                    string                `plist:"notes,omitempty" json:"notes,omitempty"`
	PackageCompleteURL       string                `plist:"PackageCompleteURL,omitempty" json:"PackageCompleteURL,omitempty"`
	PackageURL               string                `plist:"PackageURL,omitempty" json:"PackageURL,omitempty"`
	PackagePath              string                `plist:"package_path,omitempty" json:"package_path,omitempty"`
	InstallCheckScript       string                `plist:"installcheck_script,omitempty" json:"installcheck_script,omitempty"`
	UninstallCheckScript     string                `plist:"uninstallcheck_script,omitempty" json:"uninstallcheck_script,omitempty"`
	OnDemand                 bool                  `plist:"OnDemand,omitempty" json:"OnDemand,omitempty"`
	PostInstallScript        string                `plist:"postinstall_script,omitempty" json:"postinstall_script,omitempty"`
	PreInstallScript         string                `plist:"preinstall_script,omitempty" json:"preinstall_script,omitempty"`
	PostUninstallScript      string                `plist:"postuninstall_script,omitempty" json:"postuninstall_script,omitempty"`
	SuppressBundleRelocation bool                  `plist:"suppress_bundle_relocation,omitempty" json:"suppress_bundle_relocation,omitempty"`
	UnattendedInstall        bool                  `plist:"unattended_install,omitempty" json:"unattended_install,omitempty"`
	UnattendedUninstall      bool                  `plist:"unattended_uninstall,omitempty" json:"unattended_uninstall,omitempty"`
	Requires                 []string              `plist:"requires,omitempty" json:"requires,omitempty"`
	RestartAction            string                `plist:"RestartAction,omitempty" json:"RestartAction,omitempty"`
	Uninstallmethod          string                `plist:"uninstall_method,omitempty" json:"uninstall_method,omitempty"`
	UninstallScript          string                `plist:"uninstall_script,omitempty" json:"uninstall_script,omitempty"`
	UninstallerItemLocation  string                `plist:"uninstaller_item_location,omitempty" json:"uninstaller_item_location,omitempty"`
	AppleItem                bool                  `plist:"apple_item,omitempty" json:"apple_item,omitempty"`
	Uninstallable            bool                  `plist:"uninstallable,omitempty" json:"uninstallable,omitempty"`
	BlockingApplications     []string              `plist:"blocking_applications,omitempty" json:"blocking_applications,omitempty"`
	SupportedArchitectures   []string              `plist:"supported_architectures,omitempty" json:"supported_architectures,omitempty"`
	UpdateFor                []string              `plist:"update_for,omitempty" json:"update_for,omitempty"`
	Version                  string                `plist:"version,omitempty" json:"version,omitempty"`
	InstallerChoicesXML      []installerChoicesXML `plist:"installer_choices_xml,omitempty" json:"installer_choices_xml,omitempty"`
	InstallerEnvironment     map[string]string     `plist:"installer_environment,omitempty" json:"installer_environment,omitempty"`
}

type metadata struct {
	CreatedBy    string    `plist:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedDate  time.Time `plist:"creation_date,omitempty" json:"created_date,omitempty"`
	MunkiVersion string    `plist:"munki_version,omitempty" json:"munki_version,omitempty"`
	OSVersion    string    `plist:"os_version,omitempty" json:"os_version,omitempty"`
}

type install struct {
	CFBundleIdentifier         string `plist:"CFBundleIdentifier,omitempty" json:"CFBundleIdentifier,omitempty"`
	CFBundleName               string `plist:"CFBundleName,omitempty" json:"CFBundleName,omitempty"`
	CFBundleShortVersionString string `plist:"CFBundleShortVersionString,omitempty" json:"CFBundleShortVersionString,omitempty"`
	CFBundleVersion            string `plist:"CFBundleVersion,omitempty" json:"CFBundleVersion,omitempty"`
	MD5Checksum                string `plist:"md5checksum,omitempty" json:"md5checksum,omitempty"`
	MinOSVersion               string `plist:"minosversion,omitempty" json:"min_os_version,omitempty"`
	Path                       string `plist:"path,omitempty" json:"path,omitempty"`
	Type                       string `plist:"type,omitempty" json:"type,omitempty"`
	VersionComparisonKey       string `plist:"version_comparison_key,omitempty" json:"version_comparision_key,omitempty"`
}

type receipt struct {
	Filename      string `plist:"filename,omitempty" json:"filename,omitempty"`
	InstalledSize int    `plist:"installed_size,omitempty" json:"installed_size,omitempty"`
	Name          string `plist:"name,omitempty" json:"name,omitempty"`
	PackageID     string `plist:"packageid,omitempty" json:"packageid,omitempty"`
	Version       string `plist:"version,omitempty" json:"version,omitempty"`
	Optional      bool   `plist:"optional,omitempty" json:"optional,omitempty"`
}

type itemsToCopy struct {
	DestinationPath string `plist:"destination_path" json:"destination_path"`
	Group           string `plist:"group,omitempty" json:"group,omitempty"`
	Mode            string `plist:"mode,omitempty" json:"mode,omitempty"`
	SourceItem      string `plist:"source_item" json:"source_item"`
	User            string `plist:"user,omitempty" json:"user,omitempty"`
}

type installerChoicesXML struct {
	AttributeSetting int    `plist:"attributeSetting,omitempty" json:"attributeSetting,omitempty"`
	ChoiceAttribute  string `plist:"choiceAttribute,omitempty" json:"choiceAttribute,omitempty"`
	ChoiceIdentifier string `plist:"choiceIdentifier,omitempty" json:"choiceIdentifier,omitempty"`
}

type adobeRelatedItems struct {
	AdobeSetupType string                   `plist:"AdobeSetupType,omitempty" json:"AdobeSetupType,omitempty"`
	Payloads       []map[string]interface{} `plist:"payloads,omitempty" json:"payloads,omitempty"`
	// Only available in CS, not CC
	AdobeInstallInfo adobeInstallInfo `plist:"adobe_install_info,omitempty" json:"adobe_install_info,omitempty"`
}

type adobeInstallInfo struct {
	SerialNumber         string `plist:"serialnumber,omitempty" json:"serialnumber,omitempty"`
	InstallXML           string `plist:"installxml,omitempty" json:"installxml,omitempty"`
	UninstallXML         string `plist:"uninstallxml,omitempty" json:"uninstallxml,omitempty"`
	MediaSignature       string `plist:"media_signature,omitempty" json:"media_signature,omitempty"`
	MediaDigest          string `plist:"media_digest,omitempty" json:"media_signature,omitempty"`
	PayloadCount         int    `plist:"payload_count,omitempty" json:"payload_count,omitempty"`
	SuppressRegistration bool   `plist:"suppress_registration,omitempty" json:"suppress_registration,omitempty"`
	SuppressUpdates      bool   `plist:"suppress_updates,omitempty" json:"suppress_updates,omitempty"`
}

// View returns an API view
func (p *PkgsInfo) View(accept string) (*Response, error) {
	if p == nil {
		return nil, ErrNoData
	}
	return marshal(p, accept)
}

// Catalog converts a pkgsinfo to a catalog
func (p *PkgsInfo) catalog() *Catalog {
	catalog := &Catalog{
		p.pkgsinfo,
	}
	return catalog
}

type pkgsinfoView struct {
	Filename string `plist:"filename,omitempty" json:"filename,omitempty"`
	*PkgsInfo
}

// PkgsInfoCollection is a collection of pkgsinfos
type PkgsInfoCollection []*PkgsInfo

// Catalog return a specific catalog
func (p *PkgsInfoCollection) Catalog(name string) *Catalogs {
	catalogs := Catalogs{}
	var pkgsinfos *PkgsInfoCollection
	if name != "all" {
		filtered := p.ByCatalog(name)
		pkgsinfos = filtered
	} else {
		pkgsinfos = p
	}
	for _, info := range *pkgsinfos {
		catalog := info.catalog()
		catalogs = append(catalogs, catalog)
	}
	return &catalogs
}

// View returns an api response view
func (p *PkgsInfoCollection) View(accept string) (*Response, error) {
	var view []*pkgsinfoView
	for _, item := range *p {
		viewItem := &pkgsinfoView{
			item.Filename,
			item,
		}
		view = append(view, viewItem)
	}
	return marshal(view, accept)
}

// ByCatalog returns an array of items in catalog
func (p *PkgsInfoCollection) ByCatalog(catalogs ...string) *PkgsInfoCollection {
	byCatalogIndex := map[string]*PkgsInfo{}
	byCatalog := PkgsInfoCollection{}
	for _, item := range *p {
		for _, catalog := range catalogs {
			if containsString(item.Catalogs, catalog) {
				byCatalogIndex[item.Filename] = item
			}
		}
	}

	for _, v := range byCatalogIndex {
		byCatalog = append(byCatalog, v)
	}

	return &byCatalog
}

// ByName returns a list of pkgsinfos filtered by name
func (p *PkgsInfoCollection) ByName(name string) *PkgsInfoCollection {
	byName := PkgsInfoCollection{}
	for _, item := range *p {
		if item.Name == name {
			byName = append(byName, item)
		}
	}
	return &byName
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
