package anim

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var first = Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}
var second = Transform{[3]float32{1, 1, 1}, [3]float32{0, 0.025, 0}, [3]float32{0, 0, 0}}

var one = mgl32.Vec3{1, 2, 3}
var two = mgl32.Vec3{2, 3, 4}
var three = mgl32.Vec3{3, 4, 5}
var four = mgl32.Vec3{4, 5, 6}
var five = mgl32.Vec3{5, 6, 7}
var six = mgl32.Vec3{6, 7, 8}

func BenchmarkFuncLerp(b *testing.B) {
	var result Transform = Transform{}
	var t float32 = 0.5
	for i := 0; i < b.N; i++ {
		lerpTransform(&first, &second, t, &result)
	}
}

func BenchmarkInlineExpanded(b *testing.B) {
	var result Transform
	var t float32 = 0.5
	for i := 0; i < b.N; i++ {
		result.Translate[0] = first.Translate[0]*(1-t) + second.Translate[0]*t
		result.Translate[1] = first.Translate[1]*(1-t) + second.Translate[1]*t
		result.Translate[2] = first.Translate[2]*(1-t) + second.Translate[2]*t
		result.Rotation[0] = first.Rotation[0]*(1-t) + second.Rotation[0]*t
		result.Rotation[1] = first.Rotation[1]*(1-t) + second.Rotation[1]*t
		result.Rotation[2] = first.Rotation[2]*(1-t) + second.Rotation[2]*t
		result.Scale[0] = first.Scale[0]*(1-t) + second.Scale[0]*t
		result.Scale[1] = first.Scale[1]*(1-t) + second.Scale[1]*t
		result.Scale[2] = first.Scale[2]*(1-t) + second.Scale[2]*t
	}
}

func BenchmarkMathlib(b *testing.B) {
	var result mgl32.Vec3
	var t float32
	for i := 0; i < b.N; i++ {
		t = float32(i) / float32(b.N)
		result = four.Mul(1 - t).Add(one.Mul(t))
		result = five.Mul(1 - t).Add(two.Mul(t))
		result = six.Mul(1 - t).Add(three.Mul(t))
	}
	_ = result
}

func BenchmarkNormalizationMathLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		one = one.Normalize()
	}
}

func BenchmarkNormalization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if temp := one.Add(two.Mul(10)); temp.Len() < 100 {
			one = temp
		}
	}
}

func BenchmarkZeroing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		three = mgl32.Vec3{}
	}
}

func BenchmarkZeroingInline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		three[0] = 0
		three[1] = 0
		three[2] = 0
	}
}
