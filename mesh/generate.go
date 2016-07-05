package mesh

func GeneratePlane(w, h int) (floats []float32, indices []uint32, err error) {
	//Mesh variables
	floatsPerVertex := 6
	vertexCount := w * h * 6
	indiceCount := vertexCount

	floats = make([]float32, vertexCount*floatsPerVertex)
	indices = make([]uint32, indiceCount)

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			color := []float32{float32(i)/float32(w), 0.5, float32(j)/float32(w)}
			var vertices [6][]float32

			//Get the slices for the 6 vertices of a quad
			for k := range vertices{
				vertices[k] = floats[floatsPerVertex*(6*i*w+6*j+k) : floatsPerVertex*(6*i*w+6*j+k+1)]
			}
			//Top left vertex
			vertices[0][0] = float32(i)
			vertices[0][1] = float32(0)
			vertices[0][2] = float32(j)
			vertices[0][3] = color[0]
			vertices[0][4] = color[1]
			vertices[0][5] = color[2]

			vertices[1][0] = float32(i+1)
			vertices[1][1] = float32(0)
			vertices[1][2] = float32(j-1)
			vertices[1][3] = color[0]
			vertices[1][4] = color[1]
			vertices[1][5] = color[2]

			vertices[2][0] = float32(i)
			vertices[2][1] = float32(0)
			vertices[2][2] = float32(j-1)
			vertices[2][3] = color[0]
			vertices[2][4] = color[1]
			vertices[2][5] = color[2]

			vertices[3][0] = float32(i)
			vertices[3][1] = float32(0)
			vertices[3][2] = float32(j)
			vertices[3][3] = color[0]
			vertices[3][4] = color[1]
			vertices[3][5] = color[2]

			vertices[4][0] = float32(i+1)
			vertices[4][1] = float32(0)
			vertices[4][2] = float32(j)
			vertices[4][3] = color[0]
			vertices[4][4] = color[1]
			vertices[4][5] = color[2]

			vertices[5][0] = float32(i+1)
			vertices[5][1] = float32(0)
			vertices[5][2] = float32(j-1)
			vertices[5][3] = color[0]
			vertices[5][4] = color[1]
			vertices[5][5] = color[2]
		}
	}
	for i := 0; i < w*h*6; i++ {
		indices[i] = uint32(i)
	}
	err = nil
	return
}
