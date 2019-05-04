package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Clock struct {
	City    string
	Address string
	Time    string
	Conn    net.Conn
}

func main() {
	var clocks []Clock
	for _, setting := range os.Args[1:] {
		s := strings.Split(setting, "=")
		if len(s) < 2 {
			log.Fatal("setting error")
			continue
		}
		city := s[0]
		address := s[1]
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer conn.Close()
		clocks = append(clocks, Clock{City: city, Address: address, Conn: conn})
	}

	for {
		var out string
		for _, clock := range clocks {
			var time string
			if _, err := fmt.Fscanf(clock.Conn, "%s", &time); err == nil {
				clock.Time = time
				s := fmt.Sprintf("%s %s\t", clock.City, clock.Time)
				out = out + s
			}
		}
		if out == "" {
			break
		}
		fmt.Println(out)
	}
}
