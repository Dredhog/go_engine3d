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

type Mesh struct {
	Vertices []float32
	Indices  []uint32
	Textures []Texture
	VAO      uint32
	VBO      uint32
	EBO      uint32
}

type Model struct {
	Meshes []Mesh
}
