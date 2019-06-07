// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func inputText(c net.Conn, in chan<- bool) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		in <- false
		go echo(c, input.Text(), 1*time.Second)
	}
	in <- true
}

//!+
func handleConn(c net.Conn) {
	in := make(chan bool)
	go inputText(c, in)
	ticker := time.NewTicker(1 * time.Second)
	countdown := 10
	for countdown > 0 {
		select {
		case <-ticker.C:
			countdown--
			log.Println(countdown)
		case end := <-in:
			if end {
				countdown = 0
			} else {
				countdown = 10
			}
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	log.Println(c.RemoteAddr(), "out")
	ticker.Stop()
	c.Close()
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
