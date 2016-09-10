package main

import (
	"os"
	"runtime"
	"testing"
)

var (
	closures = make(chan func())
	acks     = make(chan struct{})
)

func TestMain(m *testing.M) {
	runtime.LockOSThread()
	//flag.Parse()
	go func() {
		os.Exit(m.Run())
	}()
	for f := range closures {
		f()
		acks <- struct{}{}
	}
}

func doOnMain(f func()) {
	closures <- f
	<-acks
}

func TestEngine(t *testing.T) {
	doOnMain(func() {
		main()
	})
}
