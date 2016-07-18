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
	Transform   //The transformation relative to the parent
	Name        string
	ToRoot      Transform
	Index       int
	ParentIndex int
	LocalSet    bool
	FinalSet    bool
}

type Skeleton struct {
	Transform
	RootIndex        int
	Bones            []Bone
	LocalTransforms  []mgl32.Mat4
	ToRootTransforms []mgl32.Mat4
	FinalTransforms  []mgl32.Mat4
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
