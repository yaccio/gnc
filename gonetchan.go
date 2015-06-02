package gonetchan

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

/*
EstablishChannelAsHost establishes a channel as host. Takes the channel that data
will be sent to/from, the data send through the channel (must match channel
type) and a network address to host on.
*/
func EstablishChannelAsHost(channel interface{}, value interface{}, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Printf("Connection failed: %s\n", err.Error())
			}

			go handlereads(channel, value, conn)
			go handlewrites(channel, conn)
		}
	}()

	return nil
}

/*
EstablishChannelAsClient establishes a channel as a client. Takes the channel that data
will be sent to/from, the data send through the channel (must match channel
type) and the address of the host.
*/
func EstablishChannelAsClient(channel interface{}, value interface{}, addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go handlereads(channel, value, conn)
	go handlewrites(channel, conn)

	return nil
}

func handlereads(channel interface{}, value interface{}, conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		data, err := reader.ReadBytes(0x1F)

		if err != nil {
			fmt.Printf("Connection terminated: %s\n", err.Error())
			return
		}
		data = data[:len(data)-1]

		err = json.Unmarshal(data, &value)

		if err != nil {
			fmt.Printf("Protocol broken: %s\n", err.Error())
			return
		}

		c := reflect.ValueOf(channel)
		v := reflect.ValueOf(value)
		go c.Send(v)
	}
}

func handlewrites(channel interface{}, conn net.Conn) {
	for {
		c := reflect.ValueOf(channel)
		v, ok := c.Recv()
		if !ok {
			fmt.Println("Channel closed")
			return
		}

		data, err := json.Marshal(v.Interface())
		data = append(data, 0x1F)
		if err != nil {
			fmt.Printf("JSON Marshal failed: %s\n", err.Error())
		} else {
			_, err := conn.Write(data)
			if err != nil {
				fmt.Printf("Write failed: %s\n", err.Error())
			}
		}
	}
}
