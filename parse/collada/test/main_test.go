package main

import (
	"training/engine/parse/collada"
	"log"
	"testing"
)

func TestParse(t *testing.T){
	if err := collada.Parse("cube.dae"); err != nil{
		log.Fatalf("collada test: %v", err)
	}
}
