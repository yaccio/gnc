package main

import (
	"fmt"
	"github.com/yaccio/gonetchan"
	"time"
)

func main() {
	channel := make(chan string, 2)
	gonetchan.EstablishChannelAsHost(channel, "", ":8080")

	for {
		select {
		case msg := <-channel:
			fmt.Println(msg)
		case <-time.After(time.Second):
			fmt.Println("Timed out, try again")
		}
	}
}
