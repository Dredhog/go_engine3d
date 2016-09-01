package anim

import "github.com/go-gl/mathgl/mgl32"

type Skeleton struct {
	Transform
	RootIndex          int
	Bones              []Bone
	BindShapeMatrix    mgl32.Mat4
	GlobalPoseMatrices []mgl32.Mat4
}

type Bone struct {
	Name                string
	Index               int
	ParentIndex         int
	BindPose            mgl32.Mat4
	InverseBindPose     mgl32.Mat4
	localPoseMatrix     mgl32.Mat4
	globalPoseMatrixSet bool
}

type Animator struct {
	Clips             []Clip
	CurrentClip       *Clip
	UpcomingClip      *Clip
	CurrentKeyframe   Keyframe
	UpcomingKeyframe  Keyframe
	ResultKeyframe    Keyframe
}

type Clip struct {
	Keyframes []Keyframe
	Duration  float32
}

type ClipState struct {
	index        int
	startTime    float32
	playbackRate float32
	LocalPose    Keyframe
	GlobalPose   Keyframe
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
