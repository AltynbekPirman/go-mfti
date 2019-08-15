package main

import (
	"./pp"
	"testing"
)

var data = []byte{
	128, 36, 17, 0,
	9, 0, 0, 0,
	118, 46, 114, 111, 109, 97, 110, 111, 118,
	16, 0, 0, 0,
}


func BenchmarkCodegen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := &pp.User{}
		u.Unpack(data)
	}
}

func BenchmarkReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := &pp.User{}
		unpackRef(u, data)
	}
}
