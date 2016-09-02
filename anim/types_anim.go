package anim

import "github.com/go-gl/mathgl/mgl32"

type Skeleton struct {
	Bones           []Bone
	BindShapeMatrix mgl32.Mat4
	RootIndex       int
}

type Bone struct {
	Name            string
	BindPose        mgl32.Mat4
	InverseBindPose mgl32.Mat4
	ParentIndex     int
	Index           int
}

type Animator struct {
	animations         []Animation
	animationStates    []animationState
	GlobalPoseMatrices []mgl32.Mat4
	localPoseMatrices  []mgl32.Mat4
	workingPoses	   []Keyframe
	localPose          Keyframe
	globalPosesSet     []bool
	globalTime         float32
	skeleton           *Skeleton
}

type Animation struct {
	Keyframes []Keyframe
	Duration  float32
}

type animationState struct {
	animationIndex int
	startTime      float32
	playbackRate   float32
	loop           bool
}

type Keyframe struct {
	Transforms []Transform
	SampleTime float32
}

type Transform struct {
	Scale     [3]float32
	Translate [3]float32
	Rotation  [3]float32
}
