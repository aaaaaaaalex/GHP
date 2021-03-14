package main

import (
  log "github.com/sirupsen/logrus"

  "net"
  "net/http/fcgi"
)

func main() {
  log.Info("Starting GHP...")
  server := Handler{}

  listener, err := net.Listen("tcp", ":9000")
  if err != nil {
    log.Fatal("Couldn't listen on port :9000")
    return
  }

  log.Info("Listening on port :9000")
  fcgi.Serve(listener, server)
  return
}

