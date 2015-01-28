package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func publish(msg string, channel string) (err error) {
	pubC, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}

	err = pubC.Send("PUBLISH", channel, msg)
	if err != nil {
		return err
	}

	err = pubC.Flush()
	if err != nil {
		return err
	}

	return nil
}

func subscribe(channel string) (subChan <-chan string, err error) {
	subConn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}

	subPsc := redis.PubSubConn{subConn}
	subPsc.Subscribe(channel)

	subChan_bi := make(chan string)
	go func() {
		for {
			switch v := subPsc.Receive().(type) {
			case redis.Message:
				subChan_bi <- string(v.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				return
			}
		}
	}()

	subChan = subChan_bi
	return subChan, nil
}
