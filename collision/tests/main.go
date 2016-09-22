package main

import (
	"fmt"
	"training/engine/collision"

	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	inputs := [][]mgl32.Vec3{
		[]mgl32.Vec3{mgl32.Vec3{1, 0, -1}, mgl32.Vec3{1, 0, 1}, mgl32.Vec3{-1, 0, 1}, mgl32.Vec3{0, 1, 0}},
		[]mgl32.Vec3{mgl32.Vec3{1, 0, -1}, mgl32.Vec3{1, 0, 1}, mgl32.Vec3{-1, 0, 1}, mgl32.Vec3{0, 0, 0}},
		[]mgl32.Vec3{mgl32.Vec3{1, 0, -1}, mgl32.Vec3{1, 0, 1}, mgl32.Vec3{-1, 0, 1}, mgl32.Vec3{0, 0.1, 0}},
		[]mgl32.Vec3{mgl32.Vec3{1, -1, -1}, mgl32.Vec3{1, -1, 1}, mgl32.Vec3{-1, -1, 1}, mgl32.Vec3{0, 0, 0}},
		[]mgl32.Vec3{mgl32.Vec3{1, -1, -1}, mgl32.Vec3{1, -1, 1}, mgl32.Vec3{-1, -1, 1}, mgl32.Vec3{1, 0, 1}},
		[]mgl32.Vec3{mgl32.Vec3{1, -1, -1}, mgl32.Vec3{1, -1, 1}, mgl32.Vec3{-1, -1, 1}, mgl32.Vec3{-1, 0, -1}},
		[]mgl32.Vec3{mgl32.Vec3{1, -1, -1}, mgl32.Vec3{1, -1, 1}, mgl32.Vec3{-1, -1, 1}, mgl32.Vec3{1, 0, -1}},
	}
	order := 3
	dir := mgl32.Vec3{0, 0, 1}
	for i := range inputs {
		fmt.Printf("before doSimplex3: simplex %#v;\tdir %v\n", inputs[i], dir)
		result := collision.DoSimplex3(inputs[i], &dir, &order)
		fmt.Printf("test: %v returned: %v\n", i, result)
		fmt.Printf("after doSimplex3: simplex %#v;\tdir %v; order %v\n", inputs[i], dir, order)
		fmt.Println("----------------------------------------------------------")
	}
}
