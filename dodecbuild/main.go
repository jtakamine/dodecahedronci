package main

import (
	"flag"
	"fmt"
)

func main() {
	port := parseArgs()
	fmt.Printf("Listening on port %v\n", port)

	fmt.Println("Registering with controller")
	err := rpcRegisterService()
	if err != nil {
		panic("Error registering with controller: " + err.Error())
	}
	fmt.Println("Successfully registered with controller")

	err = rpcListen(port)
	if err != nil {
		panic("Error listening on TCP port: " + err.Error())
	}
}

var parseArgs = func() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this service will listen")
	flag.Parse()
	return *portPtr
}
