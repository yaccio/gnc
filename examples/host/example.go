package main

import (
	"fmt"
	"github.com/yaccio/gnc"
)

func main() {
	channel := make(chan string)
	gnc.EstablishChannelAsHost(channel, "", ":8080")

	for {
		msg := <-channel
		fmt.Println(msg)
	}
}
