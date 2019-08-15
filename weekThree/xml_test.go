package main

import "testing"

func BenchmarkCountDecoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountDecoder()
	}
}

func BenchmarkCountStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountStruct()
	}
}