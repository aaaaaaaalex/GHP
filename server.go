package main

import (
  "bytes"
  "errors"
  "net/http"
  "io"
  "os"
  "strings"
  "text/template"

  log "github.com/sirupsen/logrus"
)

const defaultIndexFile = "index.gohtml"

type indexHandler struct {
  root http.FileSystem
  index string
}

/*
ServeHTTP serves directories that have index files, or file contents if the request doesn't target a directory.
  Borrows heavily from net/http.fileHandler source code.
*/
func (i *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var c io.ReadSeeker // the content that will ultimately be served by http.ServeContent for a successful request. Not necessarily a file.

  path := r.URL.Path // the sanitised path to use to target files
  if !strings.HasPrefix(path, "/") {
    path = "/" + path
    r.URL.Path = path
  }

  // direct requests to a dir's index file should be redirected to the parent dir's path (or should they?)
  // if strings.HasSuffix(path, i.index) {
    // TODO: redirct this
  //   return
  // {

  // can we open a file using the request path?
  f, err := i.root.Open(path)
  if err != nil {
    _, code := toHTTPError(err)
    log.Error(err)
    http.Error(w, http.StatusText(code), code)
    return
  }
  defer f.Close()
  c = f // for now, let's assume this is the source of content to serve

  // file stat modtime is used for Last-Modified headers, etc.
  s, err := f.Stat()
  if err != nil {
    _, code := toHTTPError(err)
    log.Error(err)
    http.Error(w, http.StatusText(code), code)
    return
  }

  // when the req targets a dir, we search for an index file inside
  if s.IsDir() {
    if path[len(path)-1] != '/' {
      // TODO redirect to canonical
      path = path + "/"
    }

    log.Debugf("Request '%s' targets a directory - searching dir for index '%s'", path, i.index)
    fi, err := i.root.Open(path + i.index)
    if err != nil {
      // if the dir is missing an index, we 404 (why would we assume an auto-indexed page is desired??)
      _, code := toHTTPError(err)
      log.Error(err)
      http.Error(w, http.StatusText(code), code)
      return
    }
    defer fi.Close()

    si, err := fi.Stat()
    if err != nil {
      _, code := toHTTPError(err)
      log.Error(err)
      http.Error(w, http.StatusText(code), code)
      return
    }

    // the index file is our template
    bi := make([]byte, si.Size())
    bytesRead, err := fi.Read(bi)
    if err != nil {
      _, code := toHTTPError(err)
      log.Errorf("Error reading template file '%s': %s", si.Name(),  err.Error())
      http.Error(w, http.StatusText(code), code)
      return
    }
    log.Debugf("Reading %d bytes from index template '%s%s'", bytesRead, path, i.index)

    t, err := template.New("bodyTemplate").Parse(string(bi))
    if err != nil {
      _, code := toHTTPError(err)
      log.Errorf("Error parsing template: %s", err.Error())
      http.Error(w, http.StatusText(code), code)
      return
    }

    // expose the index's siblings
    siblings, err := f.Readdir(0)
    if err != nil {
      log.Errorf("Error reading index's directory: %s", err.Error())
    }

    // prepare a buffer for the rendered content
    rb := new(bytes.Buffer)
    err = t.Execute(rb, map[string]interface{}{
      "siblings": siblings,
    })
    if err != nil {
      _, code := toHTTPError(err)
      log.Errorf("Error executing response template with model: %s", err.Error())
      http.Error(w, http.StatusText(code), code)
      return
    }
    log.Debugf("Rendered %d bytes from '%s%s'", len(rb.Bytes()), path, i.index)

    c = bytes.NewReader(rb.Bytes())
  }

  log.Debugf("Serving content for request to '%s'", path)
  http.ServeContent(w, r, s.Name(), s.ModTime(), c)
  _ = f.Close()
}

/* 
IndexServer acts much like the native http.FileServer implementation.
  However, IndexServer will *not* autoindex; if a directory has no index file, requests to it will 404.
  Requests to non-directory files will serve the file's contents.
  The name of directory index files is customisable via the 'index' argument - pass "" for the default designated by defaultIndexFile.
*/
func IndexServer(root http.FileSystem, index string) (http.Handler, error) {
  if index == "" {
    index = defaultIndexFile
  }

  if strings.Contains(index, "/") {
     return nil, errors.New("Invalid index filename passed: name should not include a leading slash, or parent directories.")
  }

  return &indexHandler{
    root: root,
    index: index,
  }, nil
}

// Copyright (c) 2009 The Go Authors. All rights reserved.
// -------------------------------------------------------
// toHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that toHTTPError does not
// actually return err.Error(), since msg and httpStatus are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
