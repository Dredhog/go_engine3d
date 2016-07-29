package collada

import (
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	if err := Parse("cube.dae"); err != nil {
		log.Fatalf("collada test: %v", err)
	}
}
