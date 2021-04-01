package main

import (
  "errors"
  "net/http"
  "os"
  "strings"
)

const DefaultIndexFile = "index.gohtml"

type indexHandler struct {
  root http.FileSystem
  index string
}

/*
ServeHTTP serves directories that have index files, or file contents if the request doesn't target a directory.
  Borrows heavily from net/http.fileHandler source code, but with a few modifications / deletions.
*/
func (i *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  if !strings.HasPrefix(path, "/") {
    path = "/" + path
    r.URL.Path = path
  }

  // direct requests to a dir's index file should be redirected to the parent dir
  // if strings.HasSuffix(path, i.index) {
    // TODO: redirct this
  //   return
  // {

  // find our file to serve
  f, err := i.root.Open(path)
  if err != nil {
    // there is no file or dir found at path. Oh, well.
    _, code := toHTTPError(err)
    http.Error(w, http.StatusText(code), code)
    return
  }

  // get file metadata
  s, err := f.Stat()
  if err != nil {
    _, code := toHTTPError(err)
    http.Error(w, http.StatusText(code), code)
    return
  }

  // search for index file within requested dirs
  if s.IsDir() {
    if path[len(path)-1] != '/' {
      // TODO redirect to canonical
      path = path + "/"
    }

    ff, err := i.root.Open(path + i.index)
    if err != nil {
      _, code := toHTTPError(err)
      http.Error(w, http.StatusText(code), code)
      return
    }

    ss, err := ff.Stat()
    if err != nil {
      _, code := toHTTPError(err)
      http.Error(w, http.StatusText(code), code)
      return
    }

    // swap f and s with the index we found
    _ = f.Close()
    f = ff
    s = ss
    // TODO set content-type before serveContent does
  }

  // We've finally found the right file! Serve its contents
  http.ServeContent(w, r, s.Name(), s.ModTime(), f)

  // didn't defer this because it felt cleaner in the case where the file handle reference can be replaced
  _ = f.Close()
}

/* 
IndexServer acts much like the native http.FileServer implementation.
  However, IndexServer will *not* autoindex; if a directory has no index file, requests to it will 404.
  Requests to non-directory files will serve the file's contents.
  The default name of a directory's index file is customisable via the 'index' argumentz.
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
