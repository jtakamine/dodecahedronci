package main

import (
	"flag"
	"fmt"
	"github.com/jtakamine/dodecahedronci/configutil"
	"net"
	"net/rpc"
	"strconv"
)

func main() {
	port := parseArgs()
	fmt.Printf("Listening on port %v\n", port)

	if !validateConfig() {
		return
	}

	err := rpc.Register(&Build{})
	if err != nil {
		panic("Error creating RPC server: " + err.Error())
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic("Error listening on TCP port: " + err.Error())
	}

	rpc.Accept(l)
}

var parseArgs = func() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this service will listen")
	flag.Parse()
	return *portPtr
}

func validateConfig() bool {
	requiredConfig := []string{
		"DODEC_HOME",
		"DODEC_DOCKER_USER",
		"DODEC_DOCKER_PASSWORD",
		"DODEC_DOCKER_EMAIL",
	}

	err := configutil.Require(requiredConfig)
	if err != nil {
		fmt.Println("An error occurred:\n" + err.Error())
		return false
	}

	return true
}
