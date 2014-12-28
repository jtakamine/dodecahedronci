package main

import (
	"flag"
	"log"
	"strconv"
)

func main() {
	port := parseArgs()
	log.Printf("Listening on port %v\n", port)

	err := ListenAndServe(":" + strconv.Itoa(port))
	if err != nil {
		log.Println("An error occurred while instantiating the server:\n", err)
	} else {
		log.Println("Server exited.")
	}
}

func parseArgs() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this server will listen")
	flag.Parse()
	return *portPtr
}
