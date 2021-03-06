package main

import (
	"flag"
	"log"
)

func main() {
	port := parseArgs()
	log.Printf("Listening on port %v\n", port)

	err := rpcListen(port)
	if err != nil {
		panic("Error listening on TCP port: " + err.Error())
	}
}

var parseArgs = func() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this service will listen")
	flag.Parse()
	return *portPtr
}
