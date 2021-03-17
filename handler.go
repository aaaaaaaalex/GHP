package main

import (
  log "github.com/sirupsen/logrus"

  "fmt"
  "net/http"
)

type Handler struct {}
func (h Handler) ServeHTTP (res http.ResponseWriter, req *http.Request) {
  log.Infof("%+v", *req)
  res.Write([]byte(fmt.Sprintf("%+v\n", *req)))
}

