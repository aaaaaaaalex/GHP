package main

import "text/template"

// Viewmodel encapsulates a set of decorations to apply to a template, eg. data, functions to expose
type Viewmodel interface {
  Data() interface{}
  Funcs() template.FuncMap
}


