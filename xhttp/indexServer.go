package xhttp

import (
  "net/http"
  "io"
  "os"
  "time"
)

/*
  indexServer is a more configurable file server
*/
type indexServer struct {
  /* default index filename
  e.g. "index.php" */
  index string

  /* root dir to serve */
  root string

  /* when on, will redirect uri paths ending with the index filename
   to the same path without it, e.g /some/path/index.php -> /some/path */
  autoRedirect bool
}

/*
  IndexServer returns a Handler similar to that of http.FileServer
*/
func IndexServer(index string, root string, autoRedirect bool) http.Handler{
  return indexServer{
    index: index,
    root: root,
    autoRedirect: autoRedirect,
  }
}

// ServeHTTP to implement the Handler interface
func (s *indexServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  filePath := root + req.URL.Path
  //TODO
  http.ServeContent(w, req, filePath)
}
