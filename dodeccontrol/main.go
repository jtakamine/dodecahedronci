package main

import (
	"flag"
	"fmt"
)

func main() {
	port, rpcPort := parseArgs()
	fmt.Printf("Listening on port %d for HTTP, port %d for RPC\n", port, rpcPort)

	rpcListen(rpcPort)

	for {
	}
}

var parseArgs = func() (port int, rpcPort int) {
	portPtr := flag.Int("port", 80, "The port on which this service will listen (HTTP requests)")
	rpcPortPtr := flag.Int("rpcport", 90, "The port on which this service will list (RPC commands)")

	flag.Parse()
	return *portPtr, *rpcPortPtr
}
