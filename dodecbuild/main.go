package main

import (
	"flag"
	"fmt"
)

func main() {
	port := parseArgs()
	fmt.Printf("Listening on port %v\n", port)

	err := rpcRegisterService()
	if err != nil {
		panic("Error registering with controller: " + err.Error())
	}

	err = rpcListen(port)
	if err != nil {
		panic("Error listening on TCP port: " + err.Error())
	}
}

var parseArgs = func() (port int) {
	flag.IntVar(&port, "port", 80, "The port on which this service will listen")
	flag.Parse()
	return port
}
