package anim

import "testing"

var firstKeyframe Keyframe= Keyframe{Ticks: 0, Transforms: []Transform{
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var secondKeyframe Keyframe = Keyframe{Ticks: 30, Transforms: []Transform{
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}}}}
var thirdKeyframe Keyframe = Keyframe{Ticks: 120, Transforms: []Transform{
	Transform{[3]float32{1, 1, 1}, [3]float32{0, 1, 0}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var fourthKeyframe Keyframe = Keyframe{Ticks: 180, Transforms: []Transform{
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{0.8, 0.8, 0.8}, [3]float32{}, [3]float32{0, 0, -30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, -30}}}}
var fifthKeyframe Keyframe = Keyframe{Ticks: 200, Transforms: []Transform{
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
var walkAnimation Animation = Animation{Keyframes: []Keyframe{firstKeyframe, secondKeyframe, thirdKeyframe, fourthKeyframe, fifthKeyframe}}

var workingKeyframe Keyframe= Keyframe{Ticks: 0, Transforms: []Transform{
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}},
	Transform{[3]float32{1, 1, 1}, [3]float32{}, [3]float32{0, 0, 0}}}}
func BenchmarkAnimationSampling(b *testing.B){
	walkAnimation.SetTotalTicks()
	for i := 0; i < b.N; i++{
		walkAnimation.LoopedLinearSample(float32(b.N), &workingKeyframe)
	}
}
