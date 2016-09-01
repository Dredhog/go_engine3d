package anim

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func NewSkeleton(bones []Bone, bindShapeMatrix mgl32.Mat4, rootIndex int) *Skeleton {
	s := Skeleton{Bones: bones, BindShapeMatrix: bindShapeMatrix, RootIndex: rootIndex}
	s.GlobalPoseMatrices = make([]mgl32.Mat4, len(bones))
	return &s
}

func (s *Skeleton) ApplyKeyframe(keyframe *Keyframe) error {
	if len(keyframe.Transforms) != len(s.Bones) {
		return fmt.Errorf("anim: Wrong number of transforms for skeleton in keyframe. expected %v, but got %v", len(s.Bones), len(keyframe.Transforms))
	}
	for i := range s.Bones {
		s.Bones[i].globalPoseMatrixSet = false
	}
	for i := range s.Bones {
		s.Bones[i].setLocalPoseMatrix(keyframe.Transforms[i])
	}
	for i := range s.Bones {
		_ = s.setGlobalPoseMatrix(i)
	}
	return nil
}

func (s *Skeleton) setGlobalPoseMatrix(boneIndex int) mgl32.Mat4 {
	b := &s.Bones[boneIndex]
	if b.Index == s.RootIndex {
		s.GlobalPoseMatrices[b.Index] = b.localPoseMatrix
	} else if !b.globalPoseMatrixSet {
		s.GlobalPoseMatrices[b.Index] = s.setGlobalPoseMatrix(b.ParentIndex).Mul4(b.localPoseMatrix)
		b.globalPoseMatrixSet = true
	}
	return s.GlobalPoseMatrices[b.Index]
}

func (b *Bone) setLocalPoseMatrix(transform Transform) {
	boneSpacePose := TransformToMat4(transform)
	b.localPoseMatrix = b.BindPose.Mul4(boneSpacePose.Mul4(b.InverseBindPose))
}
