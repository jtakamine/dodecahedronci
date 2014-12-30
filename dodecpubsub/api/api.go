package api

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func Subscribe(topic string) (err error) {
	subC, err := redis.Dial("tcp", "localhost:8000")
	if err != nil {
		return err
	}

	//the below is just testing code
	//TODO: return a channel through which messages will be passed back
	subPsc := redis.PubSubConn{subC}
	subPsc.Subscribe("example")

	go func() {
		for {
			switch v := subPsc.Receive().(type) {
			case redis.Message:
				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)

	pubC, err := redis.Dial("tcp", "localhost:8000")
	if err != nil {
		return err
	}

	err = pubC.Send("PUBLISH", "example", "hi redis!")
	if err != nil {
		return err
	}

	err = pubC.Flush()
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	return nil
}
