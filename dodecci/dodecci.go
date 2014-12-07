package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"github.com/jtakamine/dodecahedronci/handlers"
)

func main() {
	port := parseArgs()

	http.HandleFunc("/", handlers.Handle)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	if(err != nil) {
		fmt.Println("An error occurred while instantiating the http server:\n", err)
	} else {
		fmt.Println("Http server exited.")
	}
}

func parseArgs() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this server will listen")
	flag.Parse()
	return *portPtr
}
