package anim

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Scale       [3]float32
	Translation [3]float32
	Rotation    [3]float32
}

type Bone struct {
	Name                string
	Index               int
	ParentIndex         int
	BindPose            mgl32.Mat4
	InverseBindPose     mgl32.Mat4
	IndividualTransform mgl32.Mat4
	FinalSet            bool
}

type Skeleton struct {
	Transform
	RootIndex            int
	Bones                []Bone
	BindShapeMatrix      mgl32.Mat4
	FinalTransformations []mgl32.Mat4
}

type BoneState struct {
	BoneIndex int
	Transform
}

type KeyFrame struct {
	BoneStates []BoneState
	Duration   time.Duration
}

type Animation struct {
	KeyFrames     []KeyFrame
	TotalDuration time.Duration
}
