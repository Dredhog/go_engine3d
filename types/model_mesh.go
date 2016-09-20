package types

import (
	"fmt"
	"strconv"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func (m *Mesh) Init(floats []float32, indices []uint32, attrMask uint32, offsets [6]int, textures []Texture) {
	m.Floats = floats
	m.Indices = indices
	m.Textures = textures
	m.AttrMask = attrMask
	m.Offsets = offsets

	m.setUpMesh()
}

func (m *Mesh) setUpMesh() {
	gl.GenVertexArrays(1, &m.VAO)
	gl.GenBuffers(1, &m.VBO)
	gl.GenBuffers(1, &m.EBO)

	gl.BindVertexArray(m.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Floats)*4, gl.Ptr(m.Floats), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*4, gl.Ptr(m.Indices), gl.STATIC_DRAW)

	if m.AttrMask&USE_POSITIONS != 0 {
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(4*m.Offsets[0]))
	}
	if m.AttrMask&USE_NORMALS != 0 {
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(4*m.Offsets[1]))
	}
	if m.AttrMask&USE_TEXCOORDS != 0 {
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 0, gl.PtrOffset(4*m.Offsets[2]))
	}
	if m.AttrMask&USE_COLORS != 0 {
		gl.EnableVertexAttribArray(3)
		gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 0, gl.PtrOffset(4*m.Offsets[3]))
	}
	if m.AttrMask&USE_BONES != 0 {
		gl.EnableVertexAttribArray(4)
		gl.VertexAttribPointer(4, 3, gl.FLOAT, false, 0, gl.PtrOffset(4*m.Offsets[4]))
		gl.EnableVertexAttribArray(5)
		gl.VertexAttribPointer(5, 3, gl.FLOAT, false, 0, gl.PtrOffset(4*m.Offsets[5]))
	}
	fmt.Printf("Mesh Floats: %v;\tbytest: %v\n", len(m.Floats), 4*len(m.Floats))

	//Rebind default array object
	gl.BindVertexArray(0)
}

func (m *Mesh) Draw(shader uint32) {
	diffuseCount := 0
	specularCount := 0
	for i := 0; i < len(m.Textures); i++ {
		switch m.Textures[i].Type {
		case "specular":
			gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("texture_specular"+strconv.Itoa(specularCount)+"\x00")), int32(i))
			specularCount++
		case "diffuse":
			gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("texture_diffuse"+strconv.Itoa(diffuseCount)+"\x00")), int32(i))
			diffuseCount++
		default:
			panic(fmt.Errorf("mesh: unsupported texture type: %v", m.Textures[i].Type))
		}
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, m.Textures[i].Id)
	}
	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindVertexArray(m.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}
