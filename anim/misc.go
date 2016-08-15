package anim

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func TransformToMat4(t Transform) mgl32.Mat4 {
	// translate := mgl32.Translate3D(t.Translation[0], t.Translation[1], t.Translation[2])
	// scale := mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])
	rotateX := mgl32.HomogRotate3DX(mgl32.DegToRad(t.Rotation[0]))
	rotateY := mgl32.HomogRotate3DY(mgl32.DegToRad(t.Rotation[1]))
	rotateZ := mgl32.HomogRotate3DZ(mgl32.DegToRad(t.Rotation[2]))
	return rotateZ.Mul4(rotateY.Mul4(rotateX)) //translate.Mul4(scale.Mul4(rotate))
}

func InterpolateTransform(first, second *Transform, t float32) *Transform {
	result := Transform{Rotation: [3]float32{first.Rotation[0]*(1-t) + second.Rotation[0]*t, first.Rotation[1]*(1-t) + second.Rotation[1]*t, first.Rotation[2]*(1-t) + second.Rotation[2]*t}}
	return &result
}

func InterpolateKeyframe(first, second *Keyframe, t float32) (*Keyframe, error) {
	if len(first.Transforms) != len(second.Transforms) {
		return nil, fmt.Errorf("anim: mismatched keyframes when interpolating, first has %v, second has %v transforms", len(first.Transforms), len(second.Transforms))
	}
	result := Keyframe{Transforms: make([]Transform, len(first.Transforms))}
	for i := 0; i < len(first.Transforms); i++ {
		result.Transforms[i] = *InterpolateTransform(&first.Transforms[i], &second.Transforms[i], t)
	}
	return &result, nil
}
