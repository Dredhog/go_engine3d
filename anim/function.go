package anim

import "github.com/go-gl/mathgl/mgl32"
import "fmt"

func NewSkeleton(boneCount int) *Skeleton {
	s := Skeleton{}
	s.Bones = make([]Bone, boneCount)
	s.ToRootTransforms = make([]mgl32.Mat4, boneCount)
	s.LocalTransforms = make([]mgl32.Mat4, boneCount)
	s.FinalTransforms = make([]mgl32.Mat4, boneCount)
	return &s
}

func TransformToMat4(t Transform) mgl32.Mat4 {
	translate := mgl32.Translate3D(t.Translation[0], t.Translation[1], t.Translation[2])
	scale := mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])
	rotate := mgl32.HomogRotate3DZ(mgl32.DegToRad(t.Rotation[2]))
	return translate.Mul4(scale.Mul4(rotate))
}

func (s *Skeleton) CalculateFinalTransformations(transforms ...Transform) error {
	if len(s.Bones) != len(transforms) {
		return fmt.Errorf("anim: Wrong number of transforms for skeleton. expected %v, got %v", len(s.Bones), len(transforms))
	}
	for i := range s.Bones {
		s.Bones[i].Transform = transforms[i]
		s.ToRootTransforms[i] = TransformToMat4(s.Bones[i].ToRoot)
	}
	for _, b := range s.Bones {
		_ = CalcFinalTransform(&b, s)
	}
	for i := range s.Bones{
		s.Bones[i].FinalSet = false
		s.Bones[i].LocalSet = false
	}
	return nil
}

func CalcLocalTransform(b *Bone, s *Skeleton) mgl32.Mat4 {
	if !b.LocalSet {
		b.LocalSet = true
		toRoot := TransformToMat4(b.ToRoot)
		toRootInv := toRoot.Inv()
		transformation := TransformToMat4(b.Transform)
		s.LocalTransforms[b.Index] = toRootInv.Mul4(transformation.Mul4(toRoot))
		//s.LocalTransforms[b.Index] = TransformToMat4(b.ToRoot).Inv().Mul4(TransformToMat4(b.Transform).Mul4(TransformToMat4(b.ToRoot)))
	}
	return s.LocalTransforms[b.Index]
}

func CalcFinalTransform(b *Bone, s *Skeleton) mgl32.Mat4 {
	if !b.FinalSet {
		if b.Index == s.RootIndex {
			s.FinalTransforms[b.Index] = TransformToMat4(b.Transform)
		} else {
			s.FinalTransforms[b.Index] = CalcFinalTransform(&s.Bones[b.ParentIndex], s).Mul4(CalcLocalTransform(b, s))
		}
		b.FinalSet = true
	}
	return s.FinalTransforms[b.Index]
}
