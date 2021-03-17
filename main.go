package main

import (
  "flag"

  log "github.com/sirupsen/logrus"

  "net"
  "net/http/fcgi"
)

func main() {
  network := flag.String("n", "tcp", "network to communicate over")
  address := flag.String("a", "127.0.0.1:9000", "address to listen on")
  flag.Parse()

  log.Info("Starting GHP...")

  listener, err := net.Listen(*network, *address)
  if err != nil {
    log.Fatalf("Couldn't listen on %s: %s", *address, err.Error())
    return
  }

  log.Infof("Listening on %s", *address)
  server := Handler{}
  err = fcgi.Serve(listener, server)
  if err != nil {
    log.Fatal("Fatal error while listening on %s: %s", *address, err.Error())
    return
  }
}

