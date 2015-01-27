package main

import (
	"flag"
	"github.com/jtakamine/dodecahedronci/configutil"
	"log"
	"strconv"
)

func main() {
	port := parseArgs()
	log.Printf("Listening on port %v\n", port)

	if !validateConfig() {
		return
	}

	err := ListenAndServe(":" + strconv.Itoa(port))
	if err != nil {
		log.Println("An error occurred while instantiating the service:\n", err)
	} else {
		log.Println("Server exited.")
	}
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
		log.Println(err)
		return false
	}

	return true
}
