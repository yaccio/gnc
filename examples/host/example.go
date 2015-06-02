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
		time.Sleep(2 * time.Second)
		select {
		case msg := <-channel:
			fmt.Println(msg)
		default:
			fmt.Println("No new msg")
		}
	}
}
