package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Vertex struct {
	position [3]float32
	normal   [3]float32
	color	 [3]float32
}

type Texture struct {
	id      uint32
	texType string
}

type Mesh struct {
	Vertices []float32
	Indices  []uint32
	Textures []Texture
	VAO      uint32
	VBO      uint32
	EBO      uint32
}

func (m *Mesh) Init(vertices []float32, indices []uint32, textures []Texture){
	m.Vertices = vertices
	m.Indices = indices
	m.Textures = textures

	m.setUpMesh()
}

func (m *Mesh) setUpMesh(){
	gl.GenVertexArrays(1, &m.VAO)
	gl.GenBuffers(1, &m.VBO)
	gl.GenBuffers(1, &m.EBO)

	gl.BindVertexArray(m.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*4, gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*4, gl.Ptr(m.Indices), gl.STATIC_DRAW)

	//Vertex positions
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	//Vertex texture coords
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	//Vertex normals
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(5*4))
/*
	//Vertex bone indices
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 10*4, gl.PtrOffset(6*4))
	//Vertex bone weights
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 10*4, gl.PtrOffset(8*4))
*/

	//Rebind default array object
	gl.BindVertexArray(0)
}

func( m *Mesh) Draw(){
	if len(m.Textures) > 0{
		gl.BindTexture(gl.TEXTURE_2D, m.Textures[0].id)
	}
	gl.BindVertexArray(m.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}
