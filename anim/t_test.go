package anim

import (
	"testing"
	"fmt"
)

func BenchmarkCalculateFInalTransformations(b *testing.B){
	arm := NewSkeleton(6)
	arm.Bones[0] = Bone{Name: "root          ", Index: 0, ParentIndex: 0, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{},	      [3]float32{}}}
	arm.Bones[1] = Bone{Name: "neck          ", Index: 1, ParentIndex: 0, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{},	      [3]float32{}}}
	arm.Bones[2] = Bone{Name: "left_clavicle ", Index: 2, ParentIndex: 0, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{},          [3]float32{}}}
	arm.Bones[3] = Bone{Name: "right_clavicle", Index: 3, ParentIndex: 0, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{},          [3]float32{}}}
	arm.Bones[4] = Bone{Name: "left_shoulder ", Index: 4, ParentIndex: 2, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{0.2, 0, 0}, [3]float32{}}}
	arm.Bones[5] = Bone{Name: "right_shoulder", Index: 5, ParentIndex: 3, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{-0.2, 0, 0},[3]float32{}}}
	//arm.Bones[6] = Bone{Name: "left_elbow    ", Index: 6, ParentIndex: 4, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{0.4, 0, 0}, [3]float32{}}}
	//arm.Bones[7] = Bone{Name: "right_elbow   ", Index: 7, ParentIndex: 5, ToRoot: Transform{[3]float32{1, 1, 1}, [3]float32{-0.4, 0, 0},[3]float32{}}}

	angle0 := float32(0)
	angle1 := float32(0)
	angle2 := float32(0)
	angle3 := float32(0)
	angle4 := float32(0)
	angle5 := float32(0)
	//angle6 := float32(0)
	//angle7 := float32(0)

	for i := 0; i < b.N; i++{
		arm.CalculateFinalTransformations(
						Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle0}},
						Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle1}},
						Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle2}},
						Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle3}},
						Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle4}},
						Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle5}},
						//Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle6}},
						//Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, angle7}},
		)
	}
}

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
		output := TransformToMat4(input)
		_ = output
		fmt.Printf("%v\n%v", input, output)
	}
}
