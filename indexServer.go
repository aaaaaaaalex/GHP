package IndexServer

import (
  "net/http"
  "io"
  "os"
  "time"
)

/*
  IndexServer is a more configurable file server
*/
type IndexServer struct {
  /* the name of the index files to serve
  e.g. "index.php" */
  index string

  /* root dir to serve */
  root string

  /* when on, will redirect paths ending with the index filename
   to the same path without it, e.g /some/path/index.php -> /some/path*/
  autoRedirect bool

}


// ServeHTTP to implement the Handler interface
func (s *indexServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  filePath := root + req.URL.Path
  //TODO
  http.ServeContent(w, req, filePath)
}
