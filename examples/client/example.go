package main

import (
	"github.com/yaccio/gnc"
	"time"
)

func main() {
	channel := make(chan string)
	gnc.EstablishChannelAsClient(channel, "", "localhost:8080")

	for {
		channel <- "This is an example"
		time.Sleep(time.Second * 3)
	}
}
