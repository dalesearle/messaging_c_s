package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"playground/messaging"
	"playground/messaging/container"
	"playground/messaging/content/tabledata"
	"time"
)

var count int

func main() {
	conn, err := net.Dial("tcp", ":6123")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	go handleReads(conn)
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		content := createContent()
		if err := send(conn, content); err != nil {
			fmt.Println(err)
			break
		}
	}

	defer func() {
		time.Sleep(2 * time.Second)
		conn.Close()
	}()
}

func handleReads(conn net.Conn) {
	for true {
		buff := make([]byte, 3)
		_, e := conn.Read(buff)
		if e != nil {
			break
		}
		fmt.Println(string(buff))
	}
}

func createContent() messaging.Content {
	count++
	return tabledata.NewBuilder().
		SetDatabase("MSSQL").
		SetData("SOMEID", 13).
		SetData("SOMENAME", "Name").
		SetData("SOMETIME", time.Now()).
		SetData("SOMEFLOAT", 3.14).
		SetPMS("DENTRIXENTERPRISE").
		SetVersion("8.0.23.456")
}

func send(conn net.Conn, content messaging.Content) error {
	gob.Register(time.Time{})
	buf := new(bytes.Buffer)
	pkg := container.NewPackage().
		SetContent(content).
		SetReturnAddress(43).
		SetVertical(messaging.ClassicVertical)
	if err := container.NewEncoder(pkg).EncodeContainer(buf); err != nil {
		return err
	}
	fmt.Println("writing: ", buf.Len())
	conn.Write(buf.Bytes())
	return nil
}
