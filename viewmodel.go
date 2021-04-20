package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// implmentation of indexServer's viewmodel interface
type indexViewmodel struct {
	dir http.File // the directory being indexed
	data interface{} // can be any structure accepted as data by text/template.Execute
	funcs template.FuncMap
}

/* Data is used exclusively to expose the contents of the directory being served.
	this conveniently allows templates to reference the "current dir" with a dot.
	See sample index file for elaboration.
*/
func (v indexViewmodel) Data() interface{} {
	return v.data
}

// Funcs provides readonly access to state (e.g. the request)
func (v indexViewmodel) Funcs() template.FuncMap {
	return v.funcs
}


// NewIndexViewmodel builds a new indexviewmodel for request r, implements ViewmodelBuilder
func NewIndexViewmodel(r *http.Request, dir http.File) (v Viewmodel, err error) {
	// expose the index's siblings
        siblings, err := dir.Readdir(0)
        if err != nil {
		return v, fmt.Errorf("Error reading index's directory: %s", err.Error())
        }

	funcMap := template.FuncMap{
		"Request": func() *http.Request {
			return r
		},
	}

	v = indexViewmodel {
		dir: dir,
		data: siblings,
		funcs: funcMap,
	}

	return v, nil
}
