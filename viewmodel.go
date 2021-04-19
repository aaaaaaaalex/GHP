package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"text/template"
)

// implmentation of indexServer's viewmodel interface
type indexViewmodel struct {}

// Data may return any data structure accepted as data by text/template.Execute
func (v indexViewmodel) Data(r *http.Request, cd http.File) interface{} {
	// expose the index's siblings
        siblings, err := cd.Readdir(0)
        if err != nil {
		log.Errorf("Error reading index's directory: %s", err.Error())
        }
	return siblings
}

func (v indexViewmodel) Funcs() template.FuncMap {
	return template.FuncMap{}
}
