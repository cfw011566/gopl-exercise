package echo

import "strings"

func echo1(args []string) {
	var s, sep string
	for i := 1; i < len(args); i++ {
		s += sep + args[i]
		sep = " "
	}
}

func echo2(args []string) {
	var s, sep string
	for _, arg := range args {
		s += sep + arg
		sep = " "
	}
}

func echo3(args []string) {
	strings.Join(args, " ")
}
