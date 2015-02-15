package main

import (
	"flag"
	"fmt"
)

var dockerUser = "jtakamine"

func main() {
	port, rpcPort := parseArgs()
	fmt.Printf("Listening on port %d for HTTP, port %d for RPC\n", port, rpcPort)

	go rpcListen(rpcPort)
	httpListen(port)
}

var parseArgs = func() (port int, rpcPort int) {
	flag.IntVar(&port, "port", 80, "The port on which this service will listen for HTTP requests")
	flag.IntVar(&rpcPort, "rpcport", 90, "The port on which this service will listen for RPC commands")

	flag.Parse()
	return port, rpcPort
}
