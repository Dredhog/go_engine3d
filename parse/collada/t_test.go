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
	mesh, err := ParseToMesh("/../../data/model/t.dae")
	if err != nil {
		log.Fatalf("Collada test: %v\n", err)
	}
	fmt.Printf("%#v", mesh)
}

func BenchmarkParseToMesh(b *testing.B) {
	for i := 0; i < b.N; i++{
		_, err := ParseToMesh("t.dae")
		if err != nil {
			log.Fatalf("Collada test: %v\n", err)
		}
	}
}
