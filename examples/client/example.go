package main

import (
	"fmt"
	"github.com/yaccio/gonetchan"
)

func main() {
	channel := make(chan string, 2)
	gonetchan.EstablishChannelAsClient(channel, "", "localhost:8080")

	for {
		var input string
		fmt.Scanln(&input)
		channel <- input
	}
}
