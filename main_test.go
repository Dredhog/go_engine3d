package main

import (
	"runtime"
	"testing"
)

func init() {
	runtime.LockOSThread()
}
func TestMain(t *testing.T) {
	openglWork()
}
