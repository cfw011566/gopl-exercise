package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var portNumber = flag.Int("port", 8000, "port number")

func main() {
	flag.Parse()
	address := fmt.Sprintf("localhost:%d", *portNumber)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr())
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
