package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", handle)

	port := parseArgs()
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	fmt.Println(err)
}

func handle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
	dumpRequest(r)
}

func parseArgs() (port int) {
	portPtr := flag.Int("port", 80, "The port on which this server will listen")
	flag.Parse()
	return *portPtr
}

func dumpRequest(r *http.Request) {
	fmt.Printf("%v\n", r)
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	fmt.Printf("%v\n", buf.String())
}
