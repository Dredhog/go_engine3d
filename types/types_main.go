package types

type Vertex struct {
	Position [3]float32
	Normal   [3]float32
	Color    [3]float32
}

type Texture struct {
	Id      uint32
	TexType string
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
	Meshes []Mesh
}
