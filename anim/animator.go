package anim

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func NewAnimator(skeleton *Skeleton, anims []Animation) (*Animator, error) {
	//Error testing
	boneCount := len(skeleton.Bones)
	for a := range anims {
		for b := range anims[a].Keyframes {
			if len(anims[a].Keyframes[b].Transforms) != boneCount {
				return nil, fmt.Errorf("anim: Incompatible keyframe %v of animation %v, expected %v, but go %v transforms", b, a, boneCount, len(anims[a].Keyframes[b].Transforms))
			}
		}
	}
	a := Animator{s: skeleton, animations: anims}
	a.animationStates = make([]animationState, len(anims))
	for i := range a.animationStates {
		a.animationStates[i].loop = true
		a.animationStates[i].playbackRate = 0.5
	}
	a.localPose = Keyframe{Transforms: make([]Transform, boneCount)}
	a.localPoseMatrices = make([]mgl32.Mat4, boneCount)
	a.GlobalPoseMatrices = make([]mgl32.Mat4, boneCount)
	a.globalPosesSet = make([]bool, boneCount)
	return &a, nil
}

func (a *Animator) CalcGlobalPoseMatrices() {
	for i := range a.s.Bones {
		a.globalPosesSet[i] = false
	}
	for i := range a.localPoseMatrices {
		boneSpacePose := TransformToMat4(a.localPose.Transforms[i])
		a.localPoseMatrices[i] = a.s.Bones[i].BindPose.Mul4(boneSpacePose.Mul4(a.s.Bones[i].InverseBindPose))
	}
	for i := range a.s.Bones {
		_ = a.calcGlobalPoseMatrix(i)
	}
}

func (a *Animator) calcGlobalPoseMatrix(boneIndex int) mgl32.Mat4 {
	if boneIndex == a.s.RootIndex {
		a.GlobalPoseMatrices[boneIndex] = a.localPoseMatrices[boneIndex]
	} else if !a.globalPosesSet[boneIndex] {
		a.GlobalPoseMatrices[boneIndex] = a.calcGlobalPoseMatrix(a.s.Bones[boneIndex].ParentIndex).Mul4(a.localPoseMatrices[boneIndex])

	}
	a.globalPosesSet[boneIndex] = true
	return a.GlobalPoseMatrices[boneIndex]
}

func (a *Animator) Update(deltaTime float32) {
	a.globalTime += deltaTime
	a.sample(1)
	a.CalcGlobalPoseMatrices()
}

func (a *Animator) sample(animIndex int) {
	state := &a.animationStates[animIndex]
	animation := &a.animations[animIndex]
	t := state.playbackRate * (a.globalTime - state.startTime)
	if state.loop {
		t = t - animation.Duration*float32(int(t/(animation.Duration)))
	} else if t > animation.Duration {
		t = animation.Duration
	}
	animation.linearSample(t, &a.localPose)
}
