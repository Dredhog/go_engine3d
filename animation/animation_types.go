package animation

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	Position mgl32.Vec3
	Quat     mgl32.Quat
}

type Bone struct {
	Root   Transform
	Parent *Bone
	quat   mgl32.Quat
	Length float32
	Skin   []uint32
}

type Skeleton struct {
	Transform
	Joints []Bone
}
