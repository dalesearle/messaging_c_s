package main

import (
	"fmt"
	"net"
	"playground/messaging/container"
)

func main() {
	ln, err := net.Listen("tcp", ":6123")
	if err != nil {
		// TODO: handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		// The accept is an inbound connection, create a client connection here
		fmt.Println("client connected")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for true {
		pkg := container.NewPackage()
		d := container.NewDecoder(conn)
		if err := d.DecodeContainer(pkg); err != nil {
			fmt.Println(err)
			break
		}
		content, err := pkg.Content()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Sprintf("%v\n", content)
		conn.Write([]byte("rcd"))
	}
	conn.Close()
}
