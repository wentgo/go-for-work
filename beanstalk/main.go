package main

import (
	"fmt"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

func main() {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		panic(err)
	}

	// Produce jobs:
	for i := 1; i < 100; i++ {
		msg := fmt.Sprintf("Message~%d", i)
		id, err := c.Put([]byte(msg), 1, 0, 120*time.Second)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Put #%d\n", id)
	}

	// Consume jobs:
	for i := 1; i < 100; i++ {
		id, body, err := c.Reserve(5 * time.Second)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Get #%d %s\n", id, string(body))
		c.Delete(id)
	}
}
