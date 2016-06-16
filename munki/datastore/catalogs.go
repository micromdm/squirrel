package datastore

import (
	"fmt"
	"log"
	"os"

	"github.com/groob/ape/models"
	"github.com/groob/plist"
)

var makecatalogs = make(chan bool, 1)

func (r *SimpleRepo) makeCatalogs(done chan bool) {
	catalogs := map[string]*models.Catalogs{}
	pkgsinfos, err := r.AllPkgsinfos()
	if err != nil {
		log.Println(err)
	}
	allCatalogs := pkgsinfos.Catalog("all")
	catalogs["all"] = allCatalogs
	for _, info := range *allCatalogs {
		for _, catalogName := range info.Catalogs {
			catalogs[catalogName] = pkgsinfos.Catalog(catalogName)
		}
	}

	for k, v := range catalogs {
		err = r.saveCatalog(k, v)
		if err != nil {
			log.Println(err)
		}
	}
	done <- true
}

func (r *SimpleRepo) saveCatalog(name string, catalogs *models.Catalogs) error {
	catalogsPath := fmt.Sprintf("%v/catalogs/%v", r.Path, name)
	var file *os.File
	var err error
	if _, err := os.Stat(catalogsPath); err != nil {
		file, err = os.Create(catalogsPath)
	} else {
		file, err = os.OpenFile(catalogsPath, os.O_TRUNC|os.O_WRONLY, 0755)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	enc := plist.NewEncoder(file)
	enc.Indent("  ")
	return enc.Encode(catalogs)

}

//WatchCatalogs creates catalogs from pkgsinfos
func (r *SimpleRepo) WatchCatalogs() {
	done := make(chan bool, 1)
	for {
		<-makecatalogs
		go r.makeCatalogs(done)
		<-done
	}
}
