package main

import (
	"flag"
	"log"
)

func main() {
	port := parseArgs()
	log.Printf("Listening on port %v\n", port)

	err := rpcRegisterService()
	if err != nil {
		panic("Error registering with controller: " + err.Error())
	}

	err = rpcListen(port)
	if err != nil {
		panic("Error listening on TCP port: " + err.Error())
	}
}

func parseArgs() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this service will listen")
	flag.Parse()
	return *portPtr
}
