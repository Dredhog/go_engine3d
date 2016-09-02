package anim

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func NewAnimator(skeleton *Skeleton, animations []Animation) (*Animator, error) {
	//Error testing
	boneCount := len(skeleton.Bones)
	for a := range animations {
		for b := range animations[a].Keyframes {
			if len(animations[a].Keyframes[b].Transforms) != boneCount {
				return nil, fmt.Errorf("anim: Incompatible keyframe %v of animation %v, expected %v, but go %v transforms", b, a, boneCount, len(animations[a].Keyframes[b].Transforms))
			}
		}
	}
	a := Animator{skeleton: skeleton, animations: animations}
	a.animationStates = make([]animationState, len(animations))
	for i := range a.animationStates {
		a.animationStates[i].loop = true
		a.animationStates[i].playbackRate = 1
	}
	a.animationStates[2].loop = false
	a.workingPoses = make([]Keyframe, 3)
	for i := range a.workingPoses {
		a.workingPoses[i].Transforms = make([]Transform, boneCount)
	}
	a.localPose.Transforms = make([]Transform, boneCount)
	a.localPoseMatrices = make([]mgl32.Mat4, boneCount)
	a.GlobalPoseMatrices = make([]mgl32.Mat4, boneCount)
	a.globalPosesSet = make([]bool, boneCount)
	return &a, nil
}

func (a *Animator) CalcGlobalPoseMatrices() {
	for i := range a.skeleton.Bones {
		a.globalPosesSet[i] = false
	}
	for i := range a.localPoseMatrices {
		boneSpacePose := TransformToMat4(a.localPose.Transforms[i])
		a.localPoseMatrices[i] = a.skeleton.Bones[i].BindPose.Mul4(boneSpacePose.Mul4(a.skeleton.Bones[i].InverseBindPose))
	}
	for i := range a.skeleton.Bones {
		_ = a.calcGlobalPoseMatrix(i)
	}
}

func (a *Animator) calcGlobalPoseMatrix(boneIndex int) mgl32.Mat4 {
	if boneIndex == a.skeleton.RootIndex {
		a.GlobalPoseMatrices[boneIndex] = a.localPoseMatrices[boneIndex]
	} else if !a.globalPosesSet[boneIndex] {
		a.GlobalPoseMatrices[boneIndex] = a.calcGlobalPoseMatrix(a.skeleton.Bones[boneIndex].ParentIndex).Mul4(a.localPoseMatrices[boneIndex])

	}
	a.globalPosesSet[boneIndex] = true
	return a.GlobalPoseMatrices[boneIndex]
}

func (a *Animator) Update(deltaTime float32, s, h float32) {
	a.globalTime += deltaTime
	a.SampleAtGlobalTime(0, 0)
	a.SampleAtGlobalTime(1, 1)
	a.LinearBlend(0, 1, s, 0)
	a.SampleLinear(2, h, 1)
	a.AdditiveBlend(0, 1, 1.0, 0)
	a.localPose = a.workingPoses[0]
	a.CalcGlobalPoseMatrices()
}

func (a *Animator) SampleLinear(sampleIndex int, t float32, resultIndex int) {
	a.animations[sampleIndex].linearSample(t, &a.workingPoses[resultIndex])
}

func (a *Animator) LinearBlend(firstIndex, secondIndex int, t float32, resultIndex int){
	lerpKeyframe(&a.workingPoses[firstIndex], &a.workingPoses[secondIndex], t, &a.workingPoses[resultIndex])
}

func (a *Animator) AdditiveBlend(baseIndex, additiveIndex int, t float32, resultIndex int) {
	for i := 0; i < len(a.workingPoses[baseIndex].Transforms); i++ {
		addTransforms(&a.workingPoses[baseIndex].Transforms[i], &a.workingPoses[additiveIndex].Transforms[i], t, &a.workingPoses[resultIndex].Transforms[i])
	}
}


func (a *Animator) SampleAtGlobalTime(sampleIndex, resultIndex int) {
	state := &a.animationStates[sampleIndex]
	animation := &a.animations[sampleIndex]
	t := state.playbackRate * (a.globalTime - state.startTime)
	if state.loop {
		t = t - animation.Duration*float32(int(t/(animation.Duration)))
	} else if t > animation.Duration {
		t = animation.Duration
	}
	animation.linearSample(t, &a.workingPoses[resultIndex])
}
