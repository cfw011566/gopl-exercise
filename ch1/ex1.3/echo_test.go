package echo

import "testing"

func BenchmarkEcho1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		args := []string{"1 22 333 4444 55555"}
		echo1(args)
	}
}

func BenchmarkEcho2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		args := []string{"1 22 333 4444 55555"}
		echo2(args)
	}
}

func BenchmarkEcho3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		args := []string{"1 22 333 4444 55555"}
		echo3(args)
	}
}
