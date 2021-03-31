package main

import (
  "errors"
  "net/http"
  "strings"
)

const DefaultIndexFile = "index.gohtml"

type indexHandler struct {
  root http.FileSystem
  index string
}

func (i *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  upath := r.URL.Path
  if !strings.HasPrefix(upath, "/") {
    upath = "/" + upath
    r.URL.Path = upath
  }

  w.Write([]byte("WAAAAAA"))
}

/* 
IndexServer acts much like the native http.FileServer implementation.
  The main difference between the two, is IndexServer will *not* autoindex.
  The default name of a directory's index file is also customisable.
*/
func IndexServer(root http.FileSystem, index string) (http.Handler, error) {
  if strings.Contains(index, "/") {
     return nil, errors.New("Invalid index filename passed: name should not include a leading slash, or parent directories.")
  }

  return &indexHandler{
    root: root,
    index: index,
  }, nil
}

