package anim

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func NewSkeleton(bones []Bone, bindShapeMatrix mgl32.Mat4, rootIndex int) *Skeleton {
	s := Skeleton{Bones: bones, BindShapeMatrix: bindShapeMatrix, RootIndex: rootIndex}
	s.FinalMatrices = make([]mgl32.Mat4, len(bones))
	return &s
}

func (s *Skeleton) ApplyKeyframe(keyframe *Keyframe) error {
	if len(keyframe.Transforms) != len(s.Bones) {
		return fmt.Errorf("anim: Wrong number of transforms for skeleton in keyframe. expected %v, but got %v", len(s.Bones), len(keyframe.Transforms))
	}
	for i := range s.Bones {
		s.Bones[i].SetIndependentBoneMatrix(keyframe.Transforms[i], s.BindShapeMatrix)
	}
	for i := range s.Bones {
		_ = s.CalcCumulativeBoneMatrix(i)
	}
	for i := range s.Bones {
		s.Bones[i].CumulativeSet = false
	}
	return nil
}

func (s *Skeleton) CalcCumulativeBoneMatrix(boneIndex int) mgl32.Mat4 {
	b := &s.Bones[boneIndex]
	if b.Index == s.RootIndex {
		s.FinalMatrices[b.Index] = b.IndependentMatrix
	} else if !b.CumulativeSet {
		s.FinalMatrices[b.Index] = s.CalcCumulativeBoneMatrix(b.ParentIndex).Mul4(b.IndependentMatrix)
		b.CumulativeSet = true
	}
	return s.FinalMatrices[b.Index]
}

func (b *Bone) SetIndependentBoneMatrix(transform Transform, bindShape mgl32.Mat4) {
	transformation := TransformToMat4(transform)
	b.IndependentMatrix = b.BindPose.Mul4(transformation.Mul4(b.InverseBindPose))
}
