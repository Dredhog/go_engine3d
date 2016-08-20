package anim

import "github.com/go-gl/mathgl/mgl32"

type Skeleton struct {
	Transform
	RootIndex       int
	Bones           []Bone
	BindShapeMatrix mgl32.Mat4
	FinalMatrices   []mgl32.Mat4
}

type Bone struct {
	Name             string
	Index            int
	ParentIndex      int
	BindPose         mgl32.Mat4
	InverseBindPose  mgl32.Mat4
	IndependentMatrix mgl32.Mat4
	CumulativeSet    bool
}

type Animator struct {
	Animations        []Animation
	CurrentAnimation  *Animation
	UpcomingAnimation *Animation
	CurrentKeyframe   *Keyframe
	UpcomingKeyframe  *Keyframe
	TicksPerSecond    float32
	TicksIntoCurrent  float32
	TicksIntoUpcoming float32
}

type Animation struct {
	Keyframes          []Keyframe
	TotalTicks float32
}

type Keyframe struct {
	Transforms    []Transform
	Ticks float32
}

type Transform struct {
	Scale       [3]float32
	Translation [3]float32
	Rotation    [3]float32
}
