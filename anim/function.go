package anim

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func NewSkeleton(bones []Bone, bindShapeMatrix mgl32.Mat4, rootIndex int) *Skeleton {
	s := Skeleton{Bones: bones, BindShapeMatrix: bindShapeMatrix, RootIndex: rootIndex}
	s.FinalTransformations = make([]mgl32.Mat4, len(bones))
	return &s
}

func TransformToMat4(t Transform) mgl32.Mat4 {
	// translate := mgl32.Translate3D(t.Translation[0], t.Translation[1], t.Translation[2])
	// scale := mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])
	rotate := mgl32.HomogRotate3DZ(mgl32.DegToRad(t.Rotation[2]))
	return rotate //translate.Mul4(scale.Mul4(rotate))
}

func (s *Skeleton) CalculateFinalTransformations(transforms ...Transform) error {
	if len(transforms) != len(s.Bones) {
		return fmt.Errorf("anim: Wrong number of transforms for skeleton. expected %v, but got %v", len(s.Bones), len(transforms))
	}
	for i := range s.Bones {
		s.Bones[i].SetIndependentTransform(transforms[i], s.BindShapeMatrix)
	}
	for i := range s.Bones {
		_ = CalcFinalTransform(&s.Bones[i], s)
	}
	for i := range s.Bones {
		s.Bones[i].FinalSet = false
	}
	return nil
}

func (b *Bone) SetIndependentTransform(transform Transform, bindShape mgl32.Mat4) {
	transformation := TransformToMat4(transform)
	b.IndividualTransform = b.BindPose.Mul4(transformation.Mul4(b.InverseBindPose))
}

func CalcFinalTransform(b *Bone, s *Skeleton) mgl32.Mat4 {
	if b.Index == s.RootIndex {
		s.FinalTransformations[b.Index] = b.IndividualTransform
	} else if !b.FinalSet {
		s.FinalTransformations[b.Index] = CalcFinalTransform(&s.Bones[b.ParentIndex], s).Mul4(b.IndividualTransform)
		b.FinalSet = true
	}
	return s.FinalTransformations[b.Index]
}
