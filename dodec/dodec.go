package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"github.com/jtakamine/dodecahedronci/handlers"
)

var q messageQueue

func main() {
	q = inMemoryMsgQ{}

	port := parseArgs()

	http.HandleFunc("/", handlers.Handle)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	if(err != nil) {
		fmt.Println("An error occurred while instantiating the http server:\n", err)
	} else {
		fmt.Println("Http server exited.")
	}
}

type repository struct {
	Ssh_url string
}

type webHookMsg struct {
	Repository repository
}

type messageQueue interface {
	Send(msg string) int
	Receive(num int) []string
}

type inMemoryMsgQ struct {
	msgs []string
}

func (q inMemoryMsgQ) New() {
	q.msgs = make([]string, 1)
}

func (q inMemoryMsgQ) Send(msg string) int {
	q.msgs = append(q.msgs, msg)
	return 0
}

func (q inMemoryMsgQ) Receive(num int) []string {
	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
	dumpRequest(r)
	q.Send("message")
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
