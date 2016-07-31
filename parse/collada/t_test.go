package collada

import (
	"fmt"
	"log"
	"testing"
)
/*
func TestParse(t *testing.T) {
	coll, err := Parse("cube.dae")
	if err != nil {
		log.Fatalf("collada test: %v", err)
	}
	fmt.Printf("%#v\n", coll)
}
*/
func TestParseToMesh(t *testing.T) {
	mesh, err := ParseToMesh("/../../data/model/cube.dae")
	if err != nil {
		log.Fatalf("Collada test: %v\n", err)
	}
	fmt.Printf("%#v", mesh)
}
