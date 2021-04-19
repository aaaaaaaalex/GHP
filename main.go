package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"net"
	"net/http"
	"net/http/fcgi"
)

func main() {
	network := flag.String("n", "tcp", "network to communicate over")
	address := flag.String("a", "127.0.0.1:9000", "address to listen on")
	rootPath := flag.String("d", "/var/www/", "directory to serve from")
	index := flag.String("i", "", "filename of directory indexes")
	logLevel := flag.String("loglevel", "info", "one of: panic, fatal, error, warn, info, debug, trace")
	flag.Parse()

	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid log level '%s'. Exiting...", *logLevel)
		return
	}
	log.SetLevel(ll)

	log.Info("Starting GHP...")
	root := http.Dir(*rootPath)
	server, err := IndexServer(root, *index, NewIndexViewmodel)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Infof("Serving directory \"%s\"", *rootPath)

	listener, err := net.Listen(*network, *address)
	if err != nil {
		log.Fatalf("Couldn't listen on %s: %s", *address, err.Error())
		return
	}
	log.Infof("Listening on %s", *address)

	if err := fcgi.Serve(listener, server); err != nil {
		log.Fatalf("Fatal error while listening on %s: %s", *address, err.Error())
		return
	}
}
