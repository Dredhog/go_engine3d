package anim

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func NewSkeleton(bones []Bone) *Skeleton {
	s := Skeleton{Bones: bones}
	s.localTransformations = make([]mgl32.Mat4, len(bones))
	s.FinalTransformations = make([]mgl32.Mat4, len(bones))
	return &s
}

func TransformToMat4(t Transform) mgl32.Mat4 {
	translate := mgl32.Translate3D(t.Translation[0], t.Translation[1], t.Translation[2])
	scale := mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])
	rotate := mgl32.HomogRotate3DZ(mgl32.DegToRad(t.Rotation[2]))
	return translate.Mul4(scale.Mul4(rotate))
}

func (s *Skeleton) CalculateFinalTransformations(transforms ...Transform) error {
	if len(transforms) != len(s.Bones) {
		return fmt.Errorf("anim: Wrong number of transforms for skeleton. expected %v, got %v", len(s.Bones), len(transforms))
	}
	for i := range s.Bones {
		s.Bones[i].Transform = transforms[i]
	}
	for _, b := range s.Bones {
		_ = CalcFinalTransform(&b, s)
	}
	for i := range s.Bones {
		s.Bones[i].FinalSet = false
		s.Bones[i].LocalSet = false
	}
	return nil
}

func CalcLocalTransform(b *Bone, s *Skeleton) mgl32.Mat4 {
	if !b.LocalSet {
		b.LocalSet = true
		bindMatrix := b.InverseBindMatrix.Inv()
		transformation := TransformToMat4(b.Transform)
		s.localTransformations[b.Index] = bindMatrix.Mul4(transformation.Mul4(b.InverseBindMatrix))
	}
	return s.localTransformations[b.Index]
}

func CalcFinalTransform(b *Bone, s *Skeleton) mgl32.Mat4 {
	if !b.FinalSet {
		if b.Index == s.RootIndex {
			s.FinalTransformations[b.Index] = TransformToMat4(b.Transform)
		} else {
			s.FinalTransformations[b.Index] = CalcFinalTransform(&s.Bones[b.ParentIndex], s).Mul4(CalcLocalTransform(b, s))
		}
		b.FinalSet = true
	}
	return s.FinalTransformations[b.Index]
}
