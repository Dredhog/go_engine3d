package types

import "training/engine/anim"

type Texture struct {
	Id      uint32
	Type string
}

const (
	USE_POSITIONS = 1 << iota
	USE_NORMALS   = 1 << iota
	USE_TEXCOORDS = 1 << iota
	USE_COLORS    = 1 << iota
	USE_BONES     = 1 << iota
)

type Mesh struct {
	Floats   []float32
	Indices  []uint32
	Textures []Texture
	VAO      uint32
	VBO      uint32
	EBO      uint32

	AttrMask uint32
	Offsets  [6]int
}

type Model struct {
	Mesh     *Mesh
	Animator *anim.Animator
}
