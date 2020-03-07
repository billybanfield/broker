package main

import (
	"flag"
	"log"

	"github.com/billybanfield/broker/pkg/server"
)

var (
	sockAddr = flag.String("sockAddr", "/tmp/test.sock", "")
)

func main() {
	flag.Parse()
	srv := &server.Server{
		Addr: "/tmp/test.sock",
	}
	log.Fatal(srv.ListenAndServe())
}
