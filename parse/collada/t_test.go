package collada

import (
	"encoding/xml"
	"fmt"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	coll, err := Parse("../../data/model/t_baboon.dae")
	if err != nil {
		log.Fatalf("collada test: %v", err)
	}
	backInXml, err := xml.MarshalIndent(coll.LibraryControllers, "", "   ")
	if err != nil {
		log.Fatalf("collada test: marshalling back to xml error %v", err)
	}
	fmt.Println(string(backInXml))
}

/*
func TestParseToMesh(t *testing.T) {
	mesh, err := ParseToMesh("/../../data/model/t.dae")
	if err != nil {
		log.Fatalf("Collada test: %v\n", err)
	}
	fmt.Printf("%#v", mesh)
}
*/
func BenchmarkParseToMesh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseToMesh("t.dae")
		if err != nil {
			log.Fatalf("Collada test: %v\n", err)
		}
	}
}
