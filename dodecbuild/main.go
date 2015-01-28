package main

import (
	"flag"
	"fmt"
	"github.com/jtakamine/dodecahedronci/configutil"
	"strconv"
)

func main() {
	port := parseArgs()
	fmt.Printf("Listening on port %v\n", port)

	if !validateConfig() {
		return
	}

	err := ListenAndServe(":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("An error occurred while instantiating the service:\n" + err.Error())
	} else {
		fmt.Println("Server exited.")
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
		fmt.Println("An error occurred:\n" + err.Error())
		return false
	}

	return true
}
