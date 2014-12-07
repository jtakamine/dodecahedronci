package main

import (
	"flag"
	"net/http"
	"log"
	"strconv"
	"github.com/jtakamine/dodecahedronci/dodecci/internal/handlers"
)

func main() {
	port := parseArgs()

	http.HandleFunc("/", handlers.Handle)

	log.Printf("Listening on port %v\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	
	if(err != nil) {
		log.Println("An error occurred while instantiating the http server:\n", err)
	} else {
		log.Println("Http server exited.")
	}
}

func parseArgs() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this server will listen")
	flag.Parse()
	return *portPtr
}
