package api

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
	"time"
)

func Subscribe(channel string, address string) (subChan <-chan Message, err error) {
	address = strings.TrimSuffix(address, "/") + "/"
	origin := "http://localhost/"
	url := address + channel

	fmt.Println("Url: " + url)
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}

	subChan_bi := make(chan Message)

	go func() {
		defer conn.Close()
		for {
			data := make([]byte, 512)
			n, err := conn.Read(data)
			if err != nil {
				return
			}

			fmt.Printf("Received: %s.\n", data[:n])

			msg := Message{}

			json.Unmarshal(data, &msg)

			subChan_bi <- msg
		}
	}()

	subChan = subChan_bi
	return subChan, nil
}

func Publish(msg string, channel string, url string) (err error) {
	url = strings.TrimSuffix(url, "/") + "/"
	url = url + channel

	reqObj := Message{
		Text: msg,
		Time: time.Now(),
	}

	reqData, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Post(url, "application/json", strings.NewReader(string(reqData)))
	if err != nil {
		return err
	}

	return nil
}
