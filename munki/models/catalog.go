package models

// Catalog is a munki catalog
type Catalog struct {
	pkgsinfo
}

// Catalogs is an array of catalogs
type Catalogs []*Catalog
