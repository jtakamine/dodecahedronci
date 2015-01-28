package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {
	port := parseArgs()
	fmt.Printf("Listening on port %v\n", port)

	err := ListenAndServe(":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("An error occurred while instantiating the service:\n" + err.Error())
	} else {
		fmt.Println("Service exited.")
	}
}

var parseArgs = func() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this service will listen")
	flag.Parse()
	return *portPtr
}
