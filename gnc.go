package gnc

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
			setupchannels(channel, value, conn)
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

	go setupchannels(channel, value, conn)

	return nil
}

/*
setupchannels setups the channels to communicate over the network.
*/
func setupchannels(channel, value interface{}, conn net.Conn) {
	cnl := make(chan interface{})   //Wrapper of channel
	write := make(chan interface{}) //Channel that writes over network
	read := make(chan interface{})  //Channel that reads from network

	go wrapchannel(channel, cnl)
	go handlereads(read, value, conn)
	go handlewrites(write, conn)

	syncchannels(cnl, write, read)
}

/*
Wrapchannel converts any channel into an interface channel.
*/
func wrapchannel(channel interface{}, cnl chan interface{}) {
	c := reflect.ValueOf(channel)
	go func() {
		for {
			v := <-cnl
			c.Send(reflect.ValueOf(v))
		}
	}()
	go func() {
		for {
			v, ok := c.Recv()
			if !ok {
				close(cnl)
			}
			cnl <- v.Interface()
		}
	}()
}

func syncchannels(cnl, write, read chan interface{}) {
	for {
		var justRead interface{}
		select {
		case r := <-read:
			justRead = r
			cnl <- r
		case c := <-cnl:
			if c != justRead {
				write <- c
			}
		}
	}
}

func handlereads(channel chan interface{}, value interface{}, conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
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

		channel <- value
	}
}

func handlewrites(channel chan interface{}, conn net.Conn) {
	for {
		v := <-channel
		data, err := json.Marshal(v)
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
