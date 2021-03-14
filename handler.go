package main

import (
  "net/http"
)

type Handler struct {}
func (h Handler) ServeHTTP (rw http.ResponseWriter, req *http.Request) {
  rw.Write([]byte("hello"))
}

