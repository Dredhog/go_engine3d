package anim

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

func TransformToMat4(t Transform) mgl32.Mat4 {
	translate := mgl32.Translate3D(t.Translate[0], t.Translate[1], t.Translate[2])
	scale := mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])
	rotateX := mgl32.HomogRotate3DX(mgl32.DegToRad(t.Rotation[0]))
	rotateY := mgl32.HomogRotate3DY(mgl32.DegToRad(t.Rotation[1]))
	rotateZ := mgl32.HomogRotate3DZ(mgl32.DegToRad(t.Rotation[2]))
	return translate.Mul4(scale.Mul4(rotateZ.Mul4(rotateY.Mul4(rotateX))))
}

func LerpTransform(first, second *Transform, t float32, result *Transform) {
	result.Translate[0] = first.Translate[0]*(1-t) + second.Translate[0]*t
	result.Translate[1] = first.Translate[1]*(1-t) + second.Translate[1]*t
	result.Translate[2] = first.Translate[2]*(1-t) + second.Translate[2]*t
	result.Rotation[0] = first.Rotation[0]*(1-t) + second.Rotation[0]*t
	result.Rotation[1] = first.Rotation[1]*(1-t) + second.Rotation[1]*t
	result.Rotation[2] = first.Rotation[2]*(1-t) + second.Rotation[2]*t
	result.Scale[0] = first.Scale[0]*(1-t) + second.Scale[0]*t
	result.Scale[1] = first.Scale[1]*(1-t) + second.Scale[1]*t
	result.Scale[2] = first.Scale[2]*(1-t) + second.Scale[2]*t
}

func LerpKeyframe(first, second *Keyframe, t float32, result *Keyframe) error {
	if len(first.Transforms) != len(second.Transforms) {
		return fmt.Errorf("anim: Lerp: mismatched keyframe types, first has %v, second has %v transforms, expected result is %v length", len(first.Transforms), len(second.Transforms), len(result.Transforms))
	}
	for i := 0; i < len(first.Transforms); i++ {
		LerpTransform(&first.Transforms[i], &second.Transforms[i], t, &result.Transforms[i])
	}
	return nil
}
