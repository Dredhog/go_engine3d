package anim

import (
	"testing"
	"fmt"
)

func TestTransformToMat4(t *testing.T){
	inputs := []Transform{
		Transform{[3]float32{3, 10, 0}, [3]float32{0, 10, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{100,17,0}, [3]float32{10, 0, 0}, [3]float32{4, 10, 0}},
		Transform{[3]float32{0, 10, 0}, [3]float32{0,  0, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{0, 20, 0}, [3]float32{0, 20, 0}, [3]float32{10, 0, 0}},
		Transform{[3]float32{0, 19, 0}, [3]float32{0, 19, 0}, [3]float32{100,17,0}},
		Transform{[3]float32{10, 0, 0}, [3]float32{4, 10, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{0,  0, 0}, [3]float32{0,  5, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{0, 20, 0}, [3]float32{3, 10, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{0, 19, 0}, [3]float32{100,17,0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{4, 10, 0}, [3]float32{0, 10, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{0,  5, 0}, [3]float32{0, 20, 0}, [3]float32{10, 0, 0}},
		Transform{[3]float32{3, 10, 0}, [3]float32{0, 19, 0}, [3]float32{100,17,0}},
		Transform{[3]float32{100,17,0}, [3]float32{10, 0, 0}, [3]float32{4, 10, 0}},
		Transform{[3]float32{0, 10, 0}, [3]float32{0,  0, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{10, 0, 0}, [3]float32{0, 20, 0}, [3]float32{10, 0, 0}},
		Transform{[3]float32{0,  0, 0}, [3]float32{0, 19, 0}, [3]float32{100,17,0}},
		Transform{[3]float32{0, 20, 0}, [3]float32{4, 10, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{0, 19, 0}, [3]float32{0,  5, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{4, 10, 0}, [3]float32{3, 10, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{0,  5, 0}, [3]float32{100,17,0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{3, 10, 0}, [3]float32{0, 10, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{100,17,0}, [3]float32{10, 0, 0}, [3]float32{4, 10, 0}},
		Transform{[3]float32{0, 10, 0}, [3]float32{0,  0, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{10, 0, 0}, [3]float32{0, 20, 0}, [3]float32{10, 0, 0}},
		Transform{[3]float32{0,  0, 0}, [3]float32{0, 19, 0}, [3]float32{100,17,0}},
		Transform{[3]float32{0, 20, 0}, [3]float32{4, 10, 0}, [3]float32{0, 0,  30}},
		Transform{[3]float32{0, 19, 0}, [3]float32{0,  5, 0}, [3]float32{0, 0,  15}},
		Transform{[3]float32{4, 10, 0}, [3]float32{3, 10, 0}, [3]float32{0, 0,  45}},
		Transform{[3]float32{0,  5, 0}, [3]float32{100,17,0}, [3]float32{0, 0, 135}},
		Transform{[3]float32{3, 10, 0}, [3]float32{0, 10, 0}, [3]float32{0, 0,  10}},
		Transform{[3]float32{100,17,0}, [3]float32{10, 0, 0}, [3]float32{0, 0, -60}},
		Transform{[3]float32{0, 10, 0}, [3]float32{0,  0, 0}, [3]float32{0, 0, 10}},
		Transform{[3]float32{1, 1, 1}, [3]float32{0, 20, 0}, [3]float32{10, 0, 0}},
		Transform{[3]float32{1, 1, 1}, [3]float32{0, 19, 0}, [3]float32{100,17,0}},
		Transform{[3]float32{1, 1, 1}, [3]float32{4, 10, 0}, [3]float32{0, 0,  30}},
		Transform{[3]float32{1, 1, 1}, [3]float32{0,  5, 0}, [3]float32{0, 0,  15}},
		Transform{[3]float32{1, 1, 1}, [3]float32{3, 10, 0}, [3]float32{0, 0,  45}},
		Transform{[3]float32{1, 1, 1}, [3]float32{100,17,0}, [3]float32{0, 0, 135}},
		Transform{[3]float32{1, 1, 1}, [3]float32{0, 10, 0}, [3]float32{0, 0,  10}},
		Transform{[3]float32{1, 1, 1}, [3]float32{0, 20, 0}, [3]float32{0, 0, -60}},
	}
	for _, input := range inputs{
		output := TransformToMat4(&input)
		_ = output
		fmt.Printf("%v\n%v", input, output)
	}
}
