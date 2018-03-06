package main

import (
	"github.com/nats-io/go-nats"
	"fmt"
	"time"
	"sync"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer c.Close()

	// Simple Publisher

	var wait = sync.WaitGroup{}
	// Simple Async Subscriber
	var count int
	c.Subscribe("foo", func(s string) {
		count ++
		wait.Done()
		fmt.Printf("Received a message: %s %d\n", s, count)
	})

	go func() {
		for i := 0; i < 10000; i++ {
			time.Sleep(time.Nanosecond)
			c.Publish("foo", "Hello World")
		}
	}()

	wait.Wait()
	// EncodedConn can Publish any raw Go type using the registered Encoder
	type person struct {
		Name    string
		Address string
		Age     int
	}

	// Go type Subscriber
	c.Subscribe("hello", func(p *person) {
		fmt.Printf("Received a person:::: %+v\n", p)
	})

	me := &person{Name: "derek", Age: 25, Address: "140 New Montgomery Street, San Francisco, CA"}

	// Go type Publisher
	c.Publish("hello", me)

	// Unsubscribe
	sub, err := c.Subscribe("foo", nil)
	//...
	sub.Unsubscribe()

	// Requests
	var response string
	err = c.Request("help", "help me", &response, 1000*time.Millisecond)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	}

	// Replying
	c.Subscribe("help", func(subj, reply string, msg string) {
		c.Publish(reply, "I can help!")
	})

	// Close connection
	select {}
	c.Close();
}
