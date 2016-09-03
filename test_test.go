package main

import (
	"testing"
	"runtime"
)

func init(){
	runtime.LockOSThread()
}

func Benchmark(b *testing.B){
	for i := 0; i < b.N; i++ {
		program()
	}
}
