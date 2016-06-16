package models

import "io"

// PkgStore is an interface for interacting with binary files
type PkgStore interface {
	AddPkg(name string, r io.Reader) error
	DeletePkg(name string) error
}
