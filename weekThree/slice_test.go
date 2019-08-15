package main

import (
	"testing"
)

const iter = 1000

// appending to slice with pre allocated size
func allocSlice(i int) {
	bufSlice := make([]int, 0, iter)
	for j := 0; j < i; j++ {
		bufSlice = append(bufSlice, j)
	}
}


// appending to empty slice, each time len is equal to size memory will be allocated for slice with 2 x size
func emptySlice(i int){
	empSlice := make([]int, 0)
	for j := 0; j < i; j++ {
		empSlice = append(empSlice, j)
	}
}


func BenchmarkEmptySlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		emptySlice(iter)
	}
}


func BenchmarkAllocSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allocSlice(iter)
	}
}
