package collision

import (
	"training/engine/types"
	"github.com/go-gl/mathgl/mgl32"
)

func GenerateCollisionPointsFromConvexMesh(mesh *types.Mesh) []mgl32.Vec3{
	posIndCount := len(mesh.Floats)
	for i := 1; i < len(mesh.Offsets); i++{
		if offset := mesh.Offsets[i]; offset > 1 && offset < posIndCount {
			posIndCount = offset
		}
	}
	vertCount := posIndCount/9
	uniqueVerts := make(map[mgl32.Vec3]mgl32.Vec3, vertCount)
	for i := 0; i < vertCount; i++{
		for j := 0; j < 3; j++{
			vert := mgl32.Vec3{mesh.Floats[9*i+3*j], mesh.Floats[9*i+3*j+1], mesh.Floats[9*i+3*j+2]}
			_, present := uniqueVerts[vert]
			if !present {
				uniqueVerts[vert] = vert
			}
		}
	}
	collisionPoints := make([]mgl32.Vec3, len(uniqueVerts))
	i := 0
	for _, v := range uniqueVerts{
		collisionPoints[i] = v
		i++
	}
	return collisionPoints
}
