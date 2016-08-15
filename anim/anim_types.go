package anim

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Skeleton struct {
	Transform
	RootIndex            int
	Bones                []Bone
	BindShapeMatrix      mgl32.Mat4
	FinalTransformations []mgl32.Mat4
}

type Bone struct {
	Name             string
	Index            int
	ParentIndex      int
	BindPose         mgl32.Mat4
	InverseBindPose  mgl32.Mat4
	IndividualMatrix mgl32.Mat4
	CumulativeSet    bool
}

type Animator struct {
	Animations      []Animation
	CurrenAnimIndex int
	NextAnimIndex   int
	LastKeyIndex    int
	NextKeyIndex    int
	timeElapsed     time.Duration
}

type Animation struct {
	Keyframes     []Keyframe
	TotalDuration time.Duration
}

type Keyframe struct {
	Transforms []Transform
	Duration   time.Duration
}

type Transform struct {
	Scale       [3]float32
	Translation [3]float32
	Rotation    [3]float32
}
