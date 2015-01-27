package api

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func Subscribe(channel string, address string) (subChan <-chan string, err error) {
	subConn, err := redis.Dial("tcp", address)
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

func Publish(msg string, channel string, address string) (err error) {
	pubC, err := redis.Dial("tcp", address)
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
