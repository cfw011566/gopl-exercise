package main

import (
	"bytes"
	"fmt"
)

func main() {
	s := "123"
	fmt.Println(s)
	fmt.Println(comma(s))
}

func comma(s string) string {
	var buf bytes.Buffer
	t := len(s) % 3
	for i := 0; i < len(s); i++ {
		if i%3 == t && i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(s[i])
	}
	return buf.String()
}
