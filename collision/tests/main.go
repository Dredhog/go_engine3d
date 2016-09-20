package main

import (
	"fmt"
	"training/engine/collision"

	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	inputs := [][]mgl32.Vec3{
		[]mgl32.Vec3{mgl32.Vec3{1, 0, 0}, mgl32.Vec3{0, -1, -1}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 0, -0.001}},
	}
	a := 0
	v := mgl32.Vec3{0, 0, 1}
	for i := range inputs {
		result := collision.DoSimplex3(inputs[i], &v, &a)
		fmt.Printf("test: %v returned: %v\n", i, result)
	}
}
