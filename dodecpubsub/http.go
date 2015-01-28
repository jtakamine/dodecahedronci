package main

import (
	"bytes"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strings"
)

func ListenAndServe(addr string) (err error) {
	http.HandleFunc("/publish/", httpHandlePub)
	http.Handle("/subscribe/", websocket.Handler(wsHandleSub))

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}

	return nil
}

func httpHandlePub(w http.ResponseWriter, r *http.Request) {
	channel := getChannelName(r)

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Panicf("Error reading request body: %v\n", err)
	}

	msg := buf.String()

	err = publish(msg, channel)
	if err != nil {
		log.Panicf("Error publishing message: %v\n", err)
	}
}

func wsHandleSub(conn *websocket.Conn) {
	channel := getChannelName(conn.Request())

	ch, err := subscribe(channel)
	if err != nil {
		log.Panic(err)
	}

	for {
	}
	for msg := range ch {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Panic(err)
		}
	}
}

func getChannelName(r *http.Request) (channel string) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	channel = strings.Join(parts[1:], ".")

	return channel
}
